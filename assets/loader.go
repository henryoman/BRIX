package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"BRIX/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

// Embed all sprite files

//go:embed paddles/paddle.png
var paddlePNG []byte

//go:embed paddles/paddle-silver.png
var paddleSilverPNG []byte

// NOTE: Currently, the generic brick.png sprite is reused for the new brick variants.
// Replace these files with real sprites of identical names when available.

//go:embed bricks/brick-standard.png
var brickStandardPNG []byte

//go:embed bricks/brick-tusi.png
var brickTusiPNG []byte

//go:embed bricks/brick-weed.png
var brickWeedPNG []byte

//go:embed bricks/brick-columbia.png
var brickColumbiaPNG []byte

//go:embed bricks/brick-supreme.png
var brickSupremePNG []byte

// Embed level background
//
//go:embed levels/level.png
var levelBackgroundPNG []byte

// Embed start screen images
//
//go:embed startscreens/start-screen-1.png
var startScreen1PNG []byte

//go:embed startscreens/start-screen-2.png
var startScreen2PNG []byte

// Embed additional UI screens
//
//go:embed startscreens/pause-screen.png
var pauseScreenPNG []byte

//go:embed startscreens/level-complete-screen.png
var levelCompleteScreenPNG []byte

//go:embed startscreens/game-over-screen.png
var gameOverScreenPNG []byte

//go:embed startscreens/ball-lost-screen.png
var ballLostScreenPNG []byte

type Images struct {
	Paddle          *ebiten.Image
	BrickStandard   *ebiten.Image
	BrickColumbia   *ebiten.Image
	BrickSupreme    *ebiten.Image
	BrickTusi       *ebiten.Image
	BrickWeed       *ebiten.Image
	LevelBackground *ebiten.Image
	StartScreen1    *ebiten.Image
	StartScreen2    *ebiten.Image

	PauseScreen         *ebiten.Image
	LevelCompleteScreen *ebiten.Image
	GameOverScreen      *ebiten.Image
	BallLostScreen      *ebiten.Image
}

func LoadImages() (*Images, error) {
	paddle, err := loadImageFromBytes(paddleSilverPNG)
	if err != nil {
		// fallback to old paddle.png
		paddle, err = loadImageFromBytes(paddlePNG)
		if err != nil {
			return nil, err
		}
	}

	brickStandard, err := loadImageFromBytes(brickStandardPNG)
	if err != nil {
		return nil, err
	}

	brickTusi, err := loadImageFromBytes(brickTusiPNG)
	if err != nil {
		return nil, err
	}

	brickWeed, err := loadImageFromBytes(brickWeedPNG)
	if err != nil {
		return nil, err
	}

	brickColumbia, err := loadImageFromBytes(brickColumbiaPNG)
	if err != nil {
		return nil, err
	}

	brickSupreme, err := loadImageFromBytes(brickSupremePNG)
	if err != nil {
		return nil, err
	}

	levelBackground, err := loadImageFromBytes(levelBackgroundPNG)
	if err != nil {
		return nil, err
	}

	// Load start screens
	start1, err := loadImageFromBytes(startScreen1PNG)
	if err != nil {
		return nil, err
	}
	start2, err := loadImageFromBytes(startScreen2PNG)
	if err != nil {
		return nil, err
	}

	// Load new UI screens
	pauseScreen, err := loadImageFromBytes(pauseScreenPNG)
	if err != nil {
		return nil, err
	}

	levelCompleteScreen, err := loadImageFromBytes(levelCompleteScreenPNG)
	if err != nil {
		return nil, err
	}

	gameOverScreen, err := loadImageFromBytes(gameOverScreenPNG)
	if err != nil {
		return nil, err
	}

	ballLostScreen, err := loadImageFromBytes(ballLostScreenPNG)
	if err != nil {
		return nil, err
	}

	return &Images{
		Paddle:          paddle,
		BrickStandard:   brickStandard,
		BrickColumbia:   brickColumbia,
		BrickSupreme:    brickSupreme,
		BrickTusi:       brickTusi,
		BrickWeed:       brickWeed,
		LevelBackground: levelBackground,
		StartScreen1:    start1,
		StartScreen2:    start2,

		PauseScreen:         pauseScreen,
		LevelCompleteScreen: levelCompleteScreen,
		GameOverScreen:      gameOverScreen,
		BallLostScreen:      ballLostScreen,
	}, nil
}

func loadImageFromBytes(data []byte) (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func (imgs *Images) GetBrickImage(brickType entities.BrickType) *ebiten.Image {
	switch brickType {
	case entities.BrickTypeColumbia:
		return imgs.BrickColumbia
	case entities.BrickTypeSupreme:
		return imgs.BrickSupreme
	case entities.BrickTypeTusi:
		return imgs.BrickTusi
	case entities.BrickTypeWeed:
		return imgs.BrickWeed
	case entities.BrickTypeStandard:
		fallthrough
	default:
		return imgs.BrickStandard
	}
}

func (imgs *Images) GetLevelBackground(levelNum int) *ebiten.Image {
	return imgs.LevelBackground
}
