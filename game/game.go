package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"BRIX/entities"
	"BRIX/levels"
	"BRIX/physics"
	"BRIX/render"
)

// GameState represents the current state of the game
type GameState int

const (
	StateStart GameState = iota
	StatePlaying
	StatePaused
	StateLevelComplete
	StateWaitingToContinue
	StateGameOver
)

// Game encapsulates the whole game world
type Game struct {
	paddle *entities.Paddle
	ball   *entities.Ball
	bricks []*entities.Brick
	level  *levels.Level

	currentLevel int
	score        int
	lives        int // player lives
	state        GameState

	physics  *physics.CollisionSystem
	renderer *render.Renderer

	// Track the last enforced window size so we don't loop
	lastWindowW int
	lastWindowH int
}

// NewGame creates a new game instance
func NewGame() *Game {
	// Initialize renderer first since it can fail
	renderer, err := render.NewRenderer()
	if err != nil {
		log.Fatalf("Failed to create renderer: %v", err)
	}

	game := &Game{
		currentLevel: 1,
		score:        0,
		lives:        3,
		state:        StateStart,
		physics:      physics.NewCollisionSystem(),
		renderer:     renderer,
		lastWindowW:  1440,
		lastWindowH:  1080,
	}

	// Initialize game entities
	game.paddle = entities.NewPaddle()

	// Load the first level
	if err := game.loadLevel(1); err != nil {
		log.Printf("Failed to load level 1: %v", err)
		game.createFallbackLevel()
	}

	// Create ball with level's speed positioned above paddle
	game.ball = entities.NewBallAbovePaddle(game.paddle.X(), game.level.BallSpeed)

	return game
}

// loadLevel loads a level from the levels package
func (g *Game) loadLevel(levelNum int) error {
	level, err := levels.LoadLevel(levelNum)
	if err != nil {
		return err
	}

	// Guarantee score baseline: at least 1000 points per level number.
	baseline := levelNum * 1000
	if g.score < baseline {
		g.score = baseline
	}

	g.level = level
	g.bricks = make([]*entities.Brick, len(level.Bricks))

	if level.UsePixelPositioning {
		// New pixel-perfect format
		for i, levelBrick := range level.Bricks {
			g.bricks[i] = entities.NewBrickFromLevelPixel(levelBrick, level.DefaultBrickWidth, level.DefaultBrickHeight)
		}
	} else {
		// Legacy grid-based format with row-specific centering
		// Determine min/max X for every Y row so each row can be centred independently.
		rowMin := make(map[int]int)
		rowMax := make(map[int]int)
		for _, lb := range level.Bricks {
			if v, ok := rowMin[lb.Y]; !ok || lb.X < v {
				rowMin[lb.Y] = lb.X
			}
			if v, ok := rowMax[lb.Y]; !ok || lb.X > v {
				rowMax[lb.Y] = lb.X
			}
		}

		// Convert level bricks to game entities with row-specific bounds for centring.
		for i, levelBrick := range level.Bricks {
			minX := rowMin[levelBrick.Y]
			maxX := rowMax[levelBrick.Y]
			g.bricks[i] = entities.NewBrickFromLevelWithBounds(levelBrick,
				level.BrickWidth, level.BrickHeight, level.BrickSpacingX, level.BrickSpacingY,
				minX, maxX)
		}
	}

	log.Printf("Level loaded: %s with %d bricks (format: %s)", level.Name, len(g.bricks),
		map[bool]string{true: "pixel-perfect", false: "grid-based"}[level.UsePixelPositioning])
	return nil
}

// calculateBrickFieldBounds calculates the minimum and maximum X coordinates used in the level
func (g *Game) calculateBrickFieldBounds(level *levels.Level) (int, int) {
	if len(level.Bricks) == 0 {
		return 0, 0
	}

	minX := level.Bricks[0].X
	maxX := level.Bricks[0].X

	for _, brick := range level.Bricks {
		if brick.X < minX {
			minX = brick.X
		}
		if brick.X > maxX {
			maxX = brick.X
		}
	}

	return minX, maxX
}

// createFallbackLevel creates a simple level if loading fails
func (g *Game) createFallbackLevel() {
	g.level = &levels.Level{
		Name:          "Default Level",
		BrickWidth:    150,
		BrickHeight:   60,
		BrickSpacingX: 40,
		BrickSpacingY: 30,
		BallSpeed:     400,
	}

	// Create a simple pattern of bricks with fallback sizing
	g.bricks = []*entities.Brick{
		entities.NewBrickFromLevelWithBounds(entities.LevelBrick{X: 2, Y: 2, BrickType: "standard", Hits: 1}, 150, 60, 40, 30, 2, 5),
		entities.NewBrickFromLevelWithBounds(entities.LevelBrick{X: 3, Y: 2, BrickType: "standard", Hits: 1}, 150, 60, 40, 30, 2, 5),
		entities.NewBrickFromLevelWithBounds(entities.LevelBrick{X: 4, Y: 2, BrickType: "standard", Hits: 1}, 150, 60, 40, 30, 2, 5),
		entities.NewBrickFromLevelWithBounds(entities.LevelBrick{X: 5, Y: 2, BrickType: "standard", Hits: 1}, 150, 60, 40, 30, 2, 5),
	}
}

// Update implements ebiten.Game interface
func (g *Game) Update() error {
	switch g.state {
	case StateStart:
		return g.updateStart()
	case StatePlaying:
		return g.updatePlaying()
	case StatePaused:
		return g.updatePaused()
	case StateLevelComplete:
		return g.updateLevelComplete()
	case StateWaitingToContinue:
		return g.updateWaitingToContinue()
	case StateGameOver:
		return g.updateGameOver()
	}
	return nil
}

// updateStart handles start screen input
func (g *Game) updateStart() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyD) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.state = StatePlaying
	}
	return nil
}

// updatePlaying handles main game logic
func (g *Game) updatePlaying() error {
	// Check for pause input using IsKeyJustPressed to prevent flickering
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.state = StatePaused
		return nil
	}

	// Update paddle
	g.paddle.Update()

	// Update ball
	g.ball.Update()

	// Check collisions
	g.physics.CheckPaddleCollision(g.ball, g.paddle, &g.score, g.lives)
	g.physics.CheckBrickCollisions(g.ball, g.bricks, &g.score, g.lives)
	g.physics.CheckWallCollisions(g.ball)

	// Check if ball is lost
	if g.ball.IsLost() {
		g.lives-- // Subtract life immediately when ball is lost
		if g.lives <= 0 {
			g.state = StateGameOver
		} else {
			g.state = StateWaitingToContinue
		}
	}

	// Check if level is complete
	activeBricks := 0
	for _, brick := range g.bricks {
		if brick.IsActive() {
			activeBricks++
		}
	}

	if activeBricks == 0 {
		// Level complete - could advance to next level here
		g.state = StateLevelComplete
	}

	return nil
}

// updateWaitingToContinue handles waiting to continue after losing a life
func (g *Game) updateWaitingToContinue() error {
	// Check for any input to continue
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyD) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		// Reset ball position and continue playing (life already decremented)
		g.ball = entities.NewBallAbovePaddle(g.paddle.X(), g.level.BallSpeed)
		g.state = StatePlaying
	}
	return nil
}

// updateGameOver handles game over state
func (g *Game) updateGameOver() error {
	// Could handle restart logic here
	return nil
}

// updatePaused handles pause screen input
func (g *Game) updatePaused() error {
	// Check for any input to resume using IsKeyJustPressed to prevent flickering
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyD) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.state = StatePlaying
	}
	return nil
}

// updateLevelComplete handles level complete state
func (g *Game) updateLevelComplete() error {
	// Check for any input to advance to next level
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyD) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		// Try to advance to the next level
		nextLevel := g.currentLevel + 1
		if err := g.loadLevel(nextLevel); err != nil {
			// No more levels - game complete!
			log.Printf("No level %d found, game complete!", nextLevel)
			g.state = StateGameOver
		} else {
			// Successfully loaded next level
			g.currentLevel = nextLevel
			g.ball = entities.NewBallAbovePaddle(g.paddle.X(), g.level.BallSpeed)
			g.state = StatePlaying
			log.Printf("Advanced to level %d", nextLevel)
		}
	}
	return nil
}

// Draw implements ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateStart:
		g.renderer.DrawStartScreen(screen, g.level.Name)
	case StatePlaying:
		g.renderer.DrawGame(screen, g.paddle, g.ball, g.bricks, g.level.Name, g.currentLevel, g.score, g.lives)
	case StatePaused:
		g.renderer.DrawPauseScreen(screen)
	case StateLevelComplete:
		g.renderer.DrawLevelComplete(screen)
	case StateWaitingToContinue:
		g.renderer.DrawWaitingToContinue(screen, g.lives)
	case StateGameOver:
		g.renderer.DrawGameOver(screen, g.score)
	}
}

// Layout implements ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Always render the game at the fixed logical resolution.
	const logicalW, logicalH = 1440, 1080

	// Ignore the initial call where the outside size can be zero.
	if outsideWidth == 0 || outsideHeight == 0 {
		return logicalW, logicalH
	}

	// Record the very first *valid* outside size we get so we have a baseline for later comparison.
	if g.lastWindowW == 0 && g.lastWindowH == 0 {
		g.lastWindowW = outsideWidth
		g.lastWindowH = outsideHeight
		return logicalW, logicalH
	}

	// If the window is already a perfect 4:3 fit, just remember it and return.
	if outsideWidth*3 == outsideHeight*4 {
		g.lastWindowW = outsideWidth
		g.lastWindowH = outsideHeight
		return logicalW, logicalH
	}

	// Decide which dimension (width or height) the user is primarily dragging by checking
	// which changed the most since the last call.
	dw := abs(outsideWidth - g.lastWindowW)
	dh := abs(outsideHeight - g.lastWindowH)

	var desiredW, desiredH int
	if dw >= dh {
		// Width is dominating the resize – adjust height to keep 4:3.
		desiredW = outsideWidth
		desiredH = desiredW * 3 / 4
	} else {
		// Height is dominating – adjust width to keep 4:3.
		desiredH = outsideHeight
		desiredW = desiredH * 4 / 3
	}

	// Ensure we don’t enter an infinite loop by only updating when necessary.
	if desiredW != g.lastWindowW || desiredH != g.lastWindowH {
		ebiten.SetWindowSize(desiredW, desiredH)
		g.lastWindowW = desiredW
		g.lastWindowH = desiredH
	}

	return logicalW, logicalH
}

// abs is a tiny helper since Go’s standard library lacks maths on ints until Go 1.21.
func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
