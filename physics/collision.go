package physics

import (
	"brick-breaker/entities"
	"math"
)

// CollisionSystem handles all collision detection in the game
type CollisionSystem struct{}

// NewCollisionSystem creates a new collision system
func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{}
}

// CheckPaddleCollision checks if the ball collides with the paddle
func (cs *CollisionSystem) CheckPaddleCollision(ball *entities.Ball, paddle *entities.Paddle, score *int) {
	if ball.VY() <= 0 {
		return // ball moving upward, no collision possible
	}

	ballLeft, ballTop, ballRight, ballBottom := ball.GetBounds()
	paddleLeft, paddleTop, paddleRight, paddleBottom := paddle.GetBounds()

	// Check if ball overlaps with paddle
	if ballBottom >= paddleTop && ballTop <= paddleBottom &&
		ballRight >= paddleLeft && ballLeft <= paddleRight {
		// Compute offset from paddle center (-1 .. 1)
		offset := (ball.X() - paddle.X()) / (paddle.Width() / 2)
		if offset < -1 {
			offset = -1
		}
		if offset > 1 {
			offset = 1
		}

		// Maintain current speed magnitude but adjust direction
		speed := math.Hypot(ball.VX(), ball.VY())
		if speed == 0 {
			speed = 240 // fallback speed
		}

		// Limit the horizontal component to prevent shallow bounces
		// Max horizontal is 75% of speed, ensuring minimum 25% vertical
		maxHorizontal := speed * 0.75
		newVX := offset * maxHorizontal

		// Ensure strong upward movement after bounce - minimum 50% of speed
		minVertical := speed * 0.5
		verticalFromHorizontal := math.Sqrt(speed*speed - newVX*newVX)
		var newVY float64
		if verticalFromHorizontal < minVertical {
			newVY = -minVertical
			// Recalculate horizontal to maintain speed
			newVX = math.Copysign(math.Sqrt(speed*speed-newVY*newVY), newVX)
		} else {
			newVY = -verticalFromHorizontal
		}

		ball.SetVelocity(newVX, newVY)

		*score += 10 // Add points for hitting paddle
	}
}

// CheckBrickCollisions checks if the ball collides with any bricks
func (cs *CollisionSystem) CheckBrickCollisions(ball *entities.Ball, bricks []*entities.Brick, score *int, lives int) {
	ballLeft, ballTop, ballRight, ballBottom := ball.GetBounds()

	for _, brick := range bricks {
		if !brick.IsActive() {
			continue
		}

		brickLeft, brickTop, brickRight, brickBottom := brick.GetBounds()

		// Check if ball overlaps with brick
		if ballRight >= brickLeft && ballLeft <= brickRight &&
			ballBottom >= brickTop && ballTop <= brickBottom {

			// Hit the brick
			destroyed := brick.Hit()

			// Calculate points based on lives remaining
			var points int
			switch lives {
			case 3:
				points = 20
			case 2:
				points = 10
			case 1:
				points = 5
			default:
				points = 5 // fallback for any edge case
			}

			if destroyed {
				*score += points // Points for destroying a brick based on lives
			} else {
				*score += points / 2 // Half points for just hitting a brick
			}

			// Determine collision direction and bounce ball
			cs.resolveBrickCollision(ball, brickLeft, brickTop, brickRight, brickBottom)

			// Only handle one collision per frame
			break
		}
	}
}

// CheckWallCollisions checks if the ball collides with gameplay area boundaries
func (cs *CollisionSystem) CheckWallCollisions(ball *entities.Ball) {
	ballLeft, ballTop, ballRight, _ := ball.GetBounds()

	// Left and right walls of gameplay area
	if ballLeft <= entities.GameAreaLeft && ball.VX() < 0 {
		ball.ReverseX()
	}
	if ballRight >= entities.GameAreaRight && ball.VX() > 0 {
		ball.ReverseX()
	}

	// Top wall of gameplay area
	if ballTop <= entities.GameAreaTop && ball.VY() < 0 {
		ball.ReverseY()
	}

	// Note: We don't handle bottom wall here as that's handled as "ball lost" in game logic
}

// resolveBrickCollision determines the appropriate bounce direction for brick collisions
func (cs *CollisionSystem) resolveBrickCollision(ball *entities.Ball, brickLeft, brickTop, brickRight, brickBottom float64) {
	ballX, ballY := ball.X(), ball.Y()

	// Calculate distances to each edge
	distLeft := ballX - brickLeft
	distRight := brickRight - ballX
	distTop := ballY - brickTop
	distBottom := brickBottom - ballY

	// Find the minimum distance to determine collision side
	minDist := distLeft
	if distRight < minDist {
		minDist = distRight
	}
	if distTop < minDist {
		minDist = distTop
	}
	if distBottom < minDist {
		minDist = distBottom
	}

	// Bounce based on which side was hit
	if minDist == distLeft || minDist == distRight {
		ball.ReverseX()
	} else {
		ball.ReverseY()
	}
}
