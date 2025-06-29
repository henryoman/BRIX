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
	state        GameState

	physics  *physics.CollisionSystem
	renderer *render.Renderer
}

// NewGame creates a new game instance
func NewGame() *Game {
	game := &Game{
		currentLevel: 1,
		score:        0,
		state:        StateStart,
		physics:      physics.NewCollisionSystem(),
		renderer:     render.NewRenderer(),
	}

	// Initialize game entities
	game.paddle = entities.NewPaddle()
	game.ball = entities.NewBall()

	// Load the first level
	if err := game.loadLevel(1); err != nil {
		log.Printf("Failed to load level 1: %v", err)
		game.createFallbackLevel()
	}

	return game
}

// loadLevel loads a level from the levels package
func (g *Game) loadLevel(levelNum int) error {
	level, err := levels.LoadLevel(levelNum)
	if err != nil {
		return err
	}

	g.level = level
	g.bricks = make([]*entities.Brick, len(level.Bricks))

	// Convert level bricks to game entities
	for i, levelBrick := range level.Bricks {
		g.bricks[i] = entities.NewBrickFromLevel(levelBrick)
	}

	log.Printf("Level loaded: %s with %d bricks", level.Name, len(g.bricks))
	return nil
}

// createFallbackLevel creates a simple level if loading fails
func (g *Game) createFallbackLevel() {
	g.level = &levels.Level{
		Name: "Default Level",
	}

	// Create a simple pattern of bricks
	g.bricks = []*entities.Brick{
		entities.NewBrick(4, 2, "red", 1),
		entities.NewBrick(5, 2, "red", 1),
		entities.NewBrick(6, 2, "red", 1),
		entities.NewBrick(7, 2, "red", 1),
	}
}

// Update implements ebiten.Game interface
func (g *Game) Update() error {
	switch g.state {
	case StateStart:
		return g.updateStart()
	case StatePlaying:
		return g.updatePlaying()
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
	// Update paddle
	g.paddle.Update()

	// Update ball
	g.ball.Update()

	// Check collisions
	g.physics.CheckPaddleCollision(g.ball, g.paddle, &g.score)
	g.physics.CheckBrickCollisions(g.ball, g.bricks, &g.score)
	g.physics.CheckWallCollisions(g.ball)

	// Check if ball is lost
	if g.ball.IsLost() {
		g.state = StateGameOver
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
		g.state = StateGameOver
	}

	return nil
}

// updateGameOver handles game over state
func (g *Game) updateGameOver() error {
	// Could handle restart logic here
	return nil
}

// Draw implements ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateStart:
		g.renderer.DrawStartScreen(screen, g.level.Name)
	case StatePlaying:
		g.renderer.DrawGame(screen, g.paddle, g.ball, g.bricks, g.level.Name, g.currentLevel, g.score)
	case StateGameOver:
		g.renderer.DrawGameOver(screen, g.score)
	}
}

// Layout implements ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 720, 800
}
