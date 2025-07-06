package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"brick-breaker/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

// Embed all sprite files

//go:embed paddles/paddle.png
var paddlePNG []byte

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

type Images struct {
	Paddle          *ebiten.Image
	BrickStandard   *ebiten.Image
	BrickColumbia   *ebiten.Image
	BrickSupreme    *ebiten.Image
	BrickTusi       *ebiten.Image
	BrickWeed       *ebiten.Image
	LevelBackground *ebiten.Image
}

func LoadImages() (*Images, error) {
	paddle, err := loadImageFromBytes(paddlePNG)
	if err != nil {
		return nil, err
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

	return &Images{
		Paddle:          paddle,
		BrickStandard:   brickStandard,
		BrickColumbia:   brickColumbia,
		BrickSupreme:    brickSupreme,
		BrickTusi:       brickTusi,
		BrickWeed:       brickWeed,
		LevelBackground: levelBackground,
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
