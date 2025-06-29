package render

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"brick-breaker/entities"
)

// Renderer handles all drawing operations
type Renderer struct{}

// NewRenderer creates a new renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// DrawStartScreen draws the start screen
func (r *Renderer) DrawStartScreen(screen *ebiten.Image, levelName string) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{0x20, 0x20, 0x30, 0xff})

	// Simple, readable title
	ebitenutil.DebugPrintAt(screen, "BRICK BREAKER", 360-80, 200)

	// Instructions
	ebitenutil.DebugPrintAt(screen, "Press ARROW KEYS or CLICK to Start", 360-140, 350)

	// Level info
	levelText := fmt.Sprintf("Level: %s", levelName)
	ebitenutil.DebugPrintAt(screen, levelText, 360-60, 450)
}

// DrawGame draws the main game screen
func (r *Renderer) DrawGame(screen *ebiten.Image, paddle *entities.Paddle, ball *entities.Ball, bricks []*entities.Brick, levelName string, levelNum, score int) {
	// HUD background
	hud := ebiten.NewImage(720, 80)
	hud.Fill(color.Black)
	screen.DrawImage(hud, nil)

	// HUD text
	levelText := levelName
	if len(levelText) > 30 {
		levelText = levelText[:30] + "..."
	}
	ebitenutil.DebugPrintAt(screen, levelText, 10, 10)

	scoreText := fmt.Sprintf("Score: %d", score)
	ebitenutil.DebugPrintAt(screen, scoreText, 10, 30)

	levelNumText := fmt.Sprintf("Level: %d", levelNum)
	ebitenutil.DebugPrintAt(screen, levelNumText, 10, 50)

	// Bricks remaining
	activeBricks := 0
	for _, brick := range bricks {
		if brick.IsActive() {
			activeBricks++
		}
	}
	bricksText := fmt.Sprintf("Bricks: %d", activeBricks)
	ebitenutil.DebugPrintAt(screen, bricksText, 400, 20)

	// Playfield background
	playfield := ebiten.NewImage(720, 720)
	playfield.Fill(color.RGBA{0x80, 0x80, 0x80, 0xff})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 80)
	screen.DrawImage(playfield, op)

	// Draw bricks
	r.drawBricks(screen, bricks)

	// Draw paddle
	r.drawPaddle(screen, paddle)

	// Draw ball
	r.drawBall(screen, ball)
}

// DrawGameOver draws the game over screen
func (r *Renderer) DrawGameOver(screen *ebiten.Image, score int) {
	// Clear screen
	screen.Fill(color.RGBA{0x20, 0x20, 0x30, 0xff})

	// Game Over text
	ebitenutil.DebugPrintAt(screen, "GAME OVER", 360-50, 300)

	// Final score
	scoreText := fmt.Sprintf("Final Score: %d", score)
	ebitenutil.DebugPrintAt(screen, scoreText, 360-60, 400)
}

// drawBricks draws all active bricks
func (r *Renderer) drawBricks(screen *ebiten.Image, bricks []*entities.Brick) {
	for _, brick := range bricks {
		if !brick.IsActive() {
			continue
		}

		brickX, brickY := brick.GetScreenPosition()
		brickColor := brick.GetColor()

		// Draw brick as filled rectangle
		vector.DrawFilledRect(screen, float32(brickX), float32(brickY),
			float32(entities.BrickWidth), float32(entities.BrickHeight), brickColor, false)

		// Draw brick outline for better visibility
		vector.StrokeRect(screen, float32(brickX), float32(brickY),
			float32(entities.BrickWidth), float32(entities.BrickHeight), 1, color.White, false)

		// Show hit count if more than 1
		if brick.Hits() > 1 {
			hitText := fmt.Sprintf("%d", brick.Hits())
			ebitenutil.DebugPrintAt(screen, hitText,
				int(brickX)+entities.BrickWidth/2-3, int(brickY)+entities.BrickHeight/2-4)
		}
	}
}

// drawPaddle draws the paddle
func (r *Renderer) drawPaddle(screen *ebiten.Image, paddle *entities.Paddle) {
	paddleImg := ebiten.NewImage(int(paddle.Width()), int(paddle.Height()))
	paddleImg.Fill(color.RGBA{0x7f, 0xff, 0x7f, 0xff})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(paddle.X()-paddle.Width()/2, paddle.Y())
	screen.DrawImage(paddleImg, op)
}

// drawBall draws the ball as a circle
func (r *Renderer) drawBall(screen *ebiten.Image, ball *entities.Ball) {
	vector.DrawFilledCircle(screen, float32(ball.X()), float32(ball.Y()),
		float32(ball.Radius()), color.White, false)
}
