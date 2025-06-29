package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	PaddleWidth  = 120
	PaddleHeight = 20
	PaddleSpeed  = 650 // pixels per second
	PaddleY      = 700 // Y position (near bottom)

	ScreenWidth = 720
	Tick        = 1.0 / 60.0 // fixed timestep
)

// Paddle represents the player's paddle
type Paddle struct {
	x float64 // center X position
}

// NewPaddle creates a new paddle at the center of the screen
func NewPaddle() *Paddle {
	return &Paddle{
		x: ScreenWidth / 2,
	}
}

// Update handles paddle movement based on input
func (p *Paddle) Update() {
	// Left movement
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.x -= PaddleSpeed * Tick
		if p.x < PaddleWidth/2 {
			p.x = PaddleWidth / 2
		}
	}

	// Right movement
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.x += PaddleSpeed * Tick
		if p.x > ScreenWidth-PaddleWidth/2 {
			p.x = ScreenWidth - PaddleWidth/2
		}
	}
}

// X returns the center X position of the paddle
func (p *Paddle) X() float64 {
	return p.x
}

// Y returns the Y position of the paddle
func (p *Paddle) Y() float64 {
	return PaddleY
}

// Width returns the width of the paddle
func (p *Paddle) Width() float64 {
	return PaddleWidth
}

// Height returns the height of the paddle
func (p *Paddle) Height() float64 {
	return PaddleHeight
}

// GetBounds returns the paddle's bounding box for collision detection
func (p *Paddle) GetBounds() (left, top, right, bottom float64) {
	left = p.x - PaddleWidth/2
	right = p.x + PaddleWidth/2
	top = PaddleY
	bottom = PaddleY + PaddleHeight
	return
}
