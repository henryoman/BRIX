package entities

const (
	BallRadius = 10
	HUDHeight  = 30
)

// Ball represents the game ball
type Ball struct {
	x, y   float64 // center position
	vx, vy float64 // velocity
	speed  float64 // configured speed for this ball
}

// NewBall creates a new ball positioned above the paddle with default speed
func NewBall() *Ball {
	return NewBallWithSpeed(240) // default speed
}

// NewBallWithSpeed creates a new ball with configurable speed positioned at gameplay area center
func NewBallWithSpeed(speed float64) *Ball {
	return &Ball{
		x:     GameAreaLeft + GameAreaWidth/2, // center of gameplay area
		y:     PaddleY - 40,
		vx:    speed,
		vy:    -speed,
		speed: speed,
	}
}

// NewBallAbovePaddle creates a new ball positioned above the paddle with configurable speed
func NewBallAbovePaddle(paddleX float64, speed float64) *Ball {
	return &Ball{
		x:     paddleX, // position above the paddle's current location
		y:     PaddleY - 40,
		vx:    speed,
		vy:    -speed,
		speed: speed,
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

// IsLost returns true if the ball has fallen off the bottom of the gameplay area
func (b *Ball) IsLost() bool {
	return b.y > GameAreaBottom+50 // a bit below gameplay area bottom
}

// GetBounds returns the ball's bounding box for collision detection
func (b *Ball) GetBounds() (left, top, right, bottom float64) {
	left = b.x - BallRadius
	right = b.x + BallRadius
	top = b.y - BallRadius
	bottom = b.y + BallRadius
	return
}
