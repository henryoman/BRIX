package entities

const (
	BallRadius = 10
	BallSpeed  = 240 // pixels per second
	HUDHeight  = 80
)

// Ball represents the game ball
type Ball struct {
	x, y   float64 // center position
	vx, vy float64 // velocity
}

// NewBall creates a new ball positioned above the paddle
func NewBall() *Ball {
	return &Ball{
		x:  ScreenWidth / 2,
		y:  PaddleY - 40,
		vx: BallSpeed,
		vy: -BallSpeed,
	}
}

// Update handles ball movement
func (b *Ball) Update() {
	b.x += b.vx * Tick
	b.y += b.vy * Tick
}

// X returns the center X position of the ball
func (b *Ball) X() float64 {
	return b.x
}

// Y returns the center Y position of the ball
func (b *Ball) Y() float64 {
	return b.y
}

// VX returns the X velocity of the ball
func (b *Ball) VX() float64 {
	return b.vx
}

// VY returns the Y velocity of the ball
func (b *Ball) VY() float64 {
	return b.vy
}

// SetVelocity sets the ball's velocity
func (b *Ball) SetVelocity(vx, vy float64) {
	b.vx = vx
	b.vy = vy
}

// ReverseX reverses the ball's X velocity
func (b *Ball) ReverseX() {
	b.vx = -b.vx
}

// ReverseY reverses the ball's Y velocity
func (b *Ball) ReverseY() {
	b.vy = -b.vy
}

// Radius returns the ball's radius
func (b *Ball) Radius() float64 {
	return BallRadius
}

// IsLost returns true if the ball has fallen off the bottom of the screen
func (b *Ball) IsLost() bool {
	return b.y > ScreenWidth+100 // a bit below screen
}

// GetBounds returns the ball's bounding box for collision detection
func (b *Ball) GetBounds() (left, top, right, bottom float64) {
	left = b.x - BallRadius
	right = b.x + BallRadius
	top = b.y - BallRadius
	bottom = b.y + BallRadius
	return
}
