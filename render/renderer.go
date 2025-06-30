package render

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"brick-breaker/assets"
	"brick-breaker/entities"
)

// Renderer handles all drawing operations
type Renderer struct {
	images *assets.Images
}

// NewRenderer creates a new renderer with loaded images
func NewRenderer() (*Renderer, error) {
	images, err := assets.LoadImages()
	if err != nil {
		return nil, fmt.Errorf("failed to load images: %v", err)
	}

	return &Renderer{
		images: images,
	}, nil
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
func (r *Renderer) DrawGame(screen *ebiten.Image, paddle *entities.Paddle, ball *entities.Ball, bricks []*entities.Brick, levelName string, levelNum, score int, lives int) {
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

	// Lives display
	livesText := fmt.Sprintf("Lives: %d", lives)
	ebitenutil.DebugPrintAt(screen, livesText, 300, 20)

	// Bricks remaining
	activeBricks := 0
	for _, brick := range bricks {
		if brick.IsActive() {
			activeBricks++
		}
	}
	bricksText := fmt.Sprintf("Bricks: %d", activeBricks)
	ebitenutil.DebugPrintAt(screen, bricksText, 500, 20)

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

// DrawWaitingToContinue draws the waiting to continue screen
func (r *Renderer) DrawWaitingToContinue(screen *ebiten.Image, lives int) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{0x20, 0x20, 0x30, 0xff})

	// Ball lost message
	ebitenutil.DebugPrintAt(screen, "BALL LOST!", 360-50, 300)

	// Lives remaining
	livesText := fmt.Sprintf("Lives Remaining: %d", lives)
	ebitenutil.DebugPrintAt(screen, livesText, 360-80, 350)

	// Continue instruction
	ebitenutil.DebugPrintAt(screen, "Press ANY KEY to Continue", 360-100, 450)
}

// drawBricks draws all active bricks using sprite images
func (r *Renderer) drawBricks(screen *ebiten.Image, bricks []*entities.Brick) {
	for _, brick := range bricks {
		if !brick.IsActive() {
			continue
		}

		brickX, brickY := brick.GetScreenPosition()
		brickImg := r.images.GetBrickImage(brick.Color())

		// Draw brick sprite scaled to the brick's configured size
		op := &ebiten.DrawImageOptions{}

		// Scale the sprite to fit exactly into the brick's size
		imgBounds := brickImg.Bounds()
		scaleX := float64(brick.Width()) / float64(imgBounds.Dx())
		scaleY := float64(brick.Height()) / float64(imgBounds.Dy())
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(brickX, brickY)

		screen.DrawImage(brickImg, op)

		// Draw white outline for better visibility
		vector.StrokeRect(screen, float32(brickX), float32(brickY),
			float32(brick.Width()), float32(brick.Height()), 1, color.White, false)

		// Show hit count if more than 1
		if brick.Hits() > 1 {
			hitText := fmt.Sprintf("%d", brick.Hits())
			ebitenutil.DebugPrintAt(screen, hitText,
				int(brickX)+brick.Width()/2-3, int(brickY)+brick.Height()/2-4)
		}
	}
}

// drawPaddle draws the paddle using sprite image
func (r *Renderer) drawPaddle(screen *ebiten.Image, paddle *entities.Paddle) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(paddle.X()-paddle.Width()/2, paddle.Y())
	screen.DrawImage(r.images.Paddle, op)
}

// drawBall draws the ball as a circle
func (r *Renderer) drawBall(screen *ebiten.Image, ball *entities.Ball) {
	vector.DrawFilledCircle(screen, float32(ball.X()), float32(ball.Y()),
		float32(ball.Radius()), color.White, false)
}
