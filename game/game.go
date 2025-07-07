package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

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

	// --- Row-specific centering ---
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

	log.Printf("Level loaded: %s with %d bricks", level.Name, len(g.bricks))
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
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyRight) ||
		ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.state = StatePlaying
	}
	return nil
}

// updatePlaying handles main game logic
func (g *Game) updatePlaying() error {
	// Check for pause input
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) {
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
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyRight) ||
		ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) ||
		ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

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
	// Check for any input to resume
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyRight) ||
		ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) ||
		ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.state = StatePlaying
	}
	return nil
}

// updateLevelComplete handles level complete state
func (g *Game) updateLevelComplete() error {
	// Check for any input to advance to next level
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyRight) ||
		ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) ||
		ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

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
	// Enforce a 4:3 aspect-ratio on the OS window. We clamp the window size to the
	// closest multiple of the base resolution (1440×1080) that fits in the space
	// the user just requested. This lets the user resize freely from a corner while
	// preventing side-only or top/bottom-only stretching.

	// Ignore the first call where outsideWidth/outsideHeight can be zero.
	if outsideWidth > 0 && outsideHeight > 0 {
		baseW, baseH := 1440, 1080

		scaleW := float64(outsideWidth) / float64(baseW)
		scaleH := float64(outsideHeight) / float64(baseH)

		// Pick the smaller scale to ensure the game always fits.
		scale := scaleW
		if scaleH < scale {
			scale = scaleH
		}

		// Round to nearest integer pixel size.
		desiredW := int(scale * float64(baseW))
		desiredH := int(scale * float64(baseH))

		// Only update if the desired size differs from the current size.
		if (desiredW != outsideWidth || desiredH != outsideHeight) &&
			(desiredW != g.lastWindowW || desiredH != g.lastWindowH) {
			ebiten.SetWindowSize(desiredW, desiredH)
			g.lastWindowW = desiredW
			g.lastWindowH = desiredH
		}
	}

	return 1440, 1080
}
