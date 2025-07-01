package render

import (
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"

	"brick-breaker/assets"
	"brick-breaker/entities"
)

// Renderer handles all drawing operations
type Renderer struct {
	images *assets.Images
	font   font.Face
}

// NewRenderer creates a new renderer with loaded images
func NewRenderer() (*Renderer, error) {
	images, err := assets.LoadImages()
	if err != nil {
		return nil, fmt.Errorf("failed to load images: %v", err)
	}

	// --- Attempt to load Times New Roman ---
	var ttfBytes []byte
	fontCandidates := []string{
		"/System/Library/Fonts/Supplemental/Times New Roman.ttf",      // macOS Ventura +
		"/Library/Fonts/Times New Roman.ttf",                          // older macOS
		"/usr/share/fonts/truetype/msttcorefonts/Times_New_Roman.ttf", // Linux w/ msttcorefonts
		"assets/fonts/Times New Roman.ttf",                            // bundled fallback
	}

	for _, path := range fontCandidates {
		if data, err := os.ReadFile(path); err == nil {
			ttfBytes = data
			break
		}
	}

	if ttfBytes == nil {
		// Fallback to Go Regular if TNR not found
		ttfBytes = goregular.TTF
	}

	tt, err := opentype.Parse(ttfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %v", err)
	}

	const dpi = 72
	fontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20, // HUD font size
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %v", err)
	}

	return &Renderer{
		images: images,
		font:   fontFace,
	}, nil
}

// drawText draws text with the custom font
func (r *Renderer) drawText(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	text.Draw(screen, str, r.font, x, y, clr)
}

// DrawStartScreen draws the start screen
func (r *Renderer) DrawStartScreen(screen *ebiten.Image, levelName string) {
	// Clear screen with black background
	screen.Fill(color.Black)

	// Simple, readable title
	r.drawText(screen, "BRICK BREAKER", 360-80, 170, color.White)

	// Instructions
	r.drawText(screen, "Press ARROW KEYS or CLICK to Start", 360-140, 270, color.White)

	// Level info
	levelText := fmt.Sprintf("Level: %s", levelName)
	r.drawText(screen, levelText, 360-60, 370, color.White)
}

// DrawGame draws the main game screen
func (r *Renderer) DrawGame(screen *ebiten.Image, paddle *entities.Paddle, ball *entities.Ball, bricks []*entities.Brick, levelName string, levelNum, score int, lives int) {
	// HUD background (720x30)
	hud := ebiten.NewImage(720, 30)
	hud.Fill(color.Black)
	screen.DrawImage(hud, nil)

	// HUD text - single line with all info
	levelText := levelName
	if len(levelText) > 20 {
		levelText = levelText[:20] + "..."
	}
	r.drawText(screen, levelText, 10, 22, color.White)

	scoreText := fmt.Sprintf("Score: %d", score)
	r.drawText(screen, scoreText, 200, 22, color.White)

	// Lives display
	livesText := fmt.Sprintf("Lives: %d", lives)
	r.drawText(screen, livesText, 400, 22, color.White)

	// Bricks remaining
	activeBricks := 0
	for _, brick := range bricks {
		if brick.IsActive() {
			activeBricks++
		}
	}
	bricksText := fmt.Sprintf("Bricks: %d", activeBricks)
	r.drawText(screen, bricksText, 550, 22, color.White)

	// Playfield background using level-specific image (700x500 with 10px margins)
	backgroundImg := r.images.GetLevelBackground(levelNum)
	op := &ebiten.DrawImageOptions{}

	// Scale the background image to fit the 700x500 gameplay area
	imgBounds := backgroundImg.Bounds()
	scaleX := 700.0 / float64(imgBounds.Dx())
	scaleY := 500.0 / float64(imgBounds.Dy())
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(10, 30) // 10px left margin, 30px top (after HUD)
	screen.DrawImage(backgroundImg, op)

	// Draw bricks
	r.drawBricks(screen, bricks)

	// Draw paddle
	r.drawPaddle(screen, paddle)

	// Draw ball
	r.drawBall(screen, ball)
}

// DrawGameOver draws the game over screen
func (r *Renderer) DrawGameOver(screen *ebiten.Image, score int) {
	// Clear screen with black background
	screen.Fill(color.Black)

	// Game Over text
	r.drawText(screen, "GAME OVER", 360-50, 220, color.White)

	// Final score
	scoreText := fmt.Sprintf("Final Score: %d", score)
	r.drawText(screen, scoreText, 360-60, 320, color.White)
}

// DrawWaitingToContinue draws the waiting to continue screen
func (r *Renderer) DrawWaitingToContinue(screen *ebiten.Image, lives int) {
	// Clear screen with black background
	screen.Fill(color.Black)

	// Ball lost message
	r.drawText(screen, "BALL LOST!", 360-50, 200, color.White)

	// Lives remaining
	livesText := fmt.Sprintf("Lives Remaining: %d", lives)
	r.drawText(screen, livesText, 360-80, 270, color.White)

	// Continue instruction
	r.drawText(screen, "Press ANY KEY to Continue", 360-100, 370, color.White)
}

// DrawPauseScreen draws the pause screen
func (r *Renderer) DrawPauseScreen(screen *ebiten.Image) {
	// Clear screen with black background
	screen.Fill(color.Black)

	// Pause message
	r.drawText(screen, "GAME PAUSED", 360-60, 220, color.White)

	// Resume instruction
	r.drawText(screen, "Press ANY KEY to Resume", 360-100, 320, color.White)
}

// DrawLevelComplete draws the level complete screen
func (r *Renderer) DrawLevelComplete(screen *ebiten.Image) {
	// Clear screen with black background
	screen.Fill(color.Black)

	// Level complete message
	r.drawText(screen, "LEVEL COMPLETE!", 360-70, 220, color.White)

	// Continue instruction
	r.drawText(screen, "Press ANY KEY to Continue", 360-100, 320, color.White)
}

// drawBricks draws all active bricks using sprite images
func (r *Renderer) drawBricks(screen *ebiten.Image, bricks []*entities.Brick) {
	for _, brick := range bricks {
		if !brick.IsActive() {
			continue
		}

		brickX, brickY := brick.GetScreenPosition()
		brickImg := r.images.GetBrickImage(brick.Type())

		// Draw brick sprite scaled to the brick's configured size
		op := &ebiten.DrawImageOptions{}

		// Scale the sprite to fit exactly into the brick's size
		imgBounds := brickImg.Bounds()
		scaleX := float64(brick.Width()) / float64(imgBounds.Dx())
		scaleY := float64(brick.Height()) / float64(imgBounds.Dy())
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(brickX, brickY)

		screen.DrawImage(brickImg, op)

		// Draw white outline for better visibility (25% opacity)
		vector.StrokeRect(screen, float32(brickX), float32(brickY),
			float32(brick.Width()), float32(brick.Height()), 1, color.RGBA{255, 255, 255, 64}, false)

		// Show hit count if more than 1
		if brick.Hits() > 1 {
			hitText := fmt.Sprintf("%d", brick.Hits())
			r.drawText(screen, hitText,
				int(brickX)+brick.Width()/2-3, int(brickY)+brick.Height()/2-4, color.White)
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
