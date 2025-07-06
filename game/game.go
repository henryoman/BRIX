package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"brick-breaker/entities"
	"brick-breaker/levels"
	"brick-breaker/physics"
	"brick-breaker/render"
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

	// Calculate the brick field bounds for proper centering
	minX, maxX := g.calculateBrickFieldBounds(level)

	// Convert level bricks to game entities with level's sizing and field bounds
	for i, levelBrick := range level.Bricks {
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
		BrickWidth:    90,
		BrickHeight:   45,
		BrickSpacingX: 40,
		BrickSpacingY: 30,
		BallSpeed:     200,
	}

	// Create a simple pattern of bricks with fallback sizing
	// Field bounds: min=2, max=5 (4 bricks wide)
	g.bricks = []*entities.Brick{
		entities.NewBrickFromLevelWithBounds(entities.LevelBrick{X: 2, Y: 2, BrickType: "default", Hits: 1}, 90, 45, 40, 30, 2, 5),
		entities.NewBrickFromLevelWithBounds(entities.LevelBrick{X: 3, Y: 2, BrickType: "default", Hits: 1}, 90, 45, 40, 30, 2, 5),
		entities.NewBrickFromLevelWithBounds(entities.LevelBrick{X: 4, Y: 2, BrickType: "default", Hits: 1}, 90, 45, 40, 30, 2, 5),
		entities.NewBrickFromLevelWithBounds(entities.LevelBrick{X: 5, Y: 2, BrickType: "default", Hits: 1}, 90, 45, 40, 30, 2, 5),
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
	return 720, 540
}
