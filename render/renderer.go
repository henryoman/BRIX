package render

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"

	"BRIX/assets"
	"BRIX/entities"
)

// Renderer handles all drawing operations
type Renderer struct {
	images  *assets.Images
	font    font.Face
	bigFont font.Face

	startTime time.Time // reference time for start-screen flash
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

	bigFontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    80, // Large font for score
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create big font face: %v", err)
	}

	return &Renderer{
		images:    images,
		font:      fontFace,
		bigFont:   bigFontFace,
		startTime: time.Now(),
	}, nil
}

// drawText draws text with the custom font
func (r *Renderer) drawText(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	text.Draw(screen, str, r.font, x, y, clr)
}

// DrawStartScreen draws the start screen
func (r *Renderer) DrawStartScreen(screen *ebiten.Image, levelName string) {
	// Decide which start image to show based on elapsed time in the current second
	elapsed := time.Since(r.startTime)
	ms := elapsed.Milliseconds() % 1000 // cycle every second

	var img *ebiten.Image
	if ms < 700 {
		img = r.images.StartScreen1
	} else {
		img = r.images.StartScreen2
	}

	// Scale to fit the full window (1440x1080 logical size)
	op := &ebiten.DrawImageOptions{}
	bounds := img.Bounds()
	scaleX := 1440.0 / float64(bounds.Dx())
	scaleY := 1080.0 / float64(bounds.Dy())
	op.GeoM.Scale(scaleX, scaleY)
	screen.DrawImage(img, op)
}

// DrawGame draws the main game screen
func (r *Renderer) DrawGame(screen *ebiten.Image, paddle *entities.Paddle, ball *entities.Ball, bricks []*entities.Brick, levelName string, levelNum, score int, lives int) {
	// Clear entire screen so borders remain black
	screen.Fill(color.Black)

	// HUD background (1440x60)
	hud := ebiten.NewImage(1440, 60)
	hud.Fill(color.Black)
	screen.DrawImage(hud, nil)

	// HUD text - single line with all info at y=55
	levelText := levelName
	if len(levelText) > 20 {
		levelText = levelText[:20] + "..."
	}
	r.drawText(screen, levelText, 20, 45, color.White)

	scoreText := fmt.Sprintf("Score: %d", score)
	r.drawText(screen, scoreText, 400, 45, color.White)

	// Lives display
	livesText := fmt.Sprintf("Lives: %d", lives)
	r.drawText(screen, livesText, 800, 45, color.White)

	// Bricks remaining
	activeBricks := 0
	for _, brick := range bricks {
		if brick.IsActive() {
			activeBricks++
		}
	}
	bricksText := fmt.Sprintf("Bricks: %d", activeBricks)
	r.drawText(screen, bricksText, 1200, 45, color.White)

	// Playfield background using level-specific image (1400x1000)
	backgroundImg := r.images.GetLevelBackground(levelNum)
	op := &ebiten.DrawImageOptions{}

	// Scale the background image to fit the 1400x1000 gameplay area
	imgBounds := backgroundImg.Bounds()
	scaleX := 1400.0 / float64(imgBounds.Dx())
	scaleY := 1000.0 / float64(imgBounds.Dy())
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(entities.GameAreaLeft, entities.GameAreaTop) // below HUD and with left border
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
	// Draw the game over screen image scaled to the window
	img := r.images.GameOverScreen
	if img == nil {
		// Fallback if image failed to load
		screen.Fill(color.Black)
		return
	}

	op := &ebiten.DrawImageOptions{}
	bounds := img.Bounds()
	scaleX := 1440.0 / float64(bounds.Dx())
	scaleY := 1080.0 / float64(bounds.Dy())
	op.GeoM.Scale(scaleX, scaleY)
	screen.DrawImage(img, op)

	// Only display the score number (no "Final Score:" text)
	// Position centered within the box: x: 215, y: 680, width: 300px, height: 120px
	scoreText := fmt.Sprintf("%d", score)

	// Calculate center position of the box
	boxCenterX := 215 + 300/2 // 365
	boxCenterY := 680 + 120/2 // 740

	// Estimate text width for centering (big font is much wider)
	textWidth := len(scoreText) * 48 // Roughly 48px per character for 80pt font
	textX := boxCenterX - textWidth/2

	// Position text at box center using big font, lowered by 30 pixels
	text.Draw(screen, scoreText, r.bigFont, textX, boxCenterY+40, color.White)
}

// DrawWaitingToContinue draws the waiting to continue screen
func (r *Renderer) DrawWaitingToContinue(screen *ebiten.Image, lives int) {
	// Draw the ball lost screen image scaled to the window
	img := r.images.BallLostScreen
	if img == nil {
		// Fallback if image failed to load
		screen.Fill(color.Black)
		return
	}

	op := &ebiten.DrawImageOptions{}
	bounds := img.Bounds()
	scaleX := 1440.0 / float64(bounds.Dx())
	scaleY := 1080.0 / float64(bounds.Dy())
	op.GeoM.Scale(scaleX, scaleY)
	screen.DrawImage(img, op)
}

// DrawPauseScreen draws the pause screen
func (r *Renderer) DrawPauseScreen(screen *ebiten.Image) {
	// Draw the supplied pause screen image scaled to the window.
	img := r.images.PauseScreen
	if img == nil {
		// If the image failed to load, just clear to black without overlay text.
		screen.Fill(color.Black)
		return
	}

	op := &ebiten.DrawImageOptions{}
	bounds := img.Bounds()
	scaleX := 1440.0 / float64(bounds.Dx())
	scaleY := 1080.0 / float64(bounds.Dy())
	op.GeoM.Scale(scaleX, scaleY)
	screen.DrawImage(img, op)
}

// DrawLevelComplete draws the level complete screen
func (r *Renderer) DrawLevelComplete(screen *ebiten.Image) {
	// Draw the supplied level complete screen image scaled to the window.
	img := r.images.LevelCompleteScreen
	if img == nil {
		screen.Fill(color.Black)
		return
	}

	op := &ebiten.DrawImageOptions{}
	bounds := img.Bounds()
	scaleX := 1440.0 / float64(bounds.Dx())
	scaleY := 1080.0 / float64(bounds.Dy())
	op.GeoM.Scale(scaleX, scaleY)
	screen.DrawImage(img, op)
}

// drawBricks draws all active bricks using sprite images
func (r *Renderer) drawBricks(screen *ebiten.Image, bricks []*entities.Brick) {

	for _, brick := range bricks {
		if !brick.IsActive() {
			continue
		}

		brickX, brickY := brick.GetScreenPosition()
		brickImg := r.images.GetBrickImage(brick.Type())
		brickWidth := float32(brick.Width())
		brickHeight := float32(brick.Height())

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
			brickWidth, brickHeight, 1.0, color.RGBA{255, 255, 255, 64}, false)

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

	imgBounds := r.images.Paddle.Bounds()
	scaleX := paddle.Width() / float64(imgBounds.Dx())
	scaleY := paddle.Height() / float64(imgBounds.Dy())
	op.GeoM.Scale(scaleX, scaleY)

	op.GeoM.Translate(paddle.X()-paddle.Width()/2, paddle.Y())
	screen.DrawImage(r.images.Paddle, op)
}

// drawBall draws the ball as a circle
func (r *Renderer) drawBall(screen *ebiten.Image, ball *entities.Ball) {
	vector.DrawFilledCircle(screen, float32(ball.X()), float32(ball.Y()),
		float32(ball.Radius()), color.White, false)
}
