package assets

import (
	"bytes"
	_ "embed"
	"image"

	"brick-breaker/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

// Embed all sprite files
//
//go:embed paddle.png
var paddlePNG []byte

//go:embed brick.png
var brickPNG []byte

//go:embed brick-green.png
var brickGreenPNG []byte

//go:embed brick-blue.png
var brickBluePNG []byte

//go:embed brick-columbia.png
var brickColumbiaPNG []byte

//go:embed brick-supreme.png
var brickSupremePNG []byte

// Images holds all loaded game sprites
type Images struct {
	Paddle        *ebiten.Image
	Brick         *ebiten.Image
	BrickGreen    *ebiten.Image
	BrickBlue     *ebiten.Image
	BrickColumbia *ebiten.Image
	BrickSupreme  *ebiten.Image
}

// LoadImages loads all embedded sprites into memory
func LoadImages() (*Images, error) {
	paddle, err := loadImageFromBytes(paddlePNG)
	if err != nil {
		return nil, err
	}

	brick, err := loadImageFromBytes(brickPNG)
	if err != nil {
		return nil, err
	}

	brickGreen, err := loadImageFromBytes(brickGreenPNG)
	if err != nil {
		return nil, err
	}

	brickBlue, err := loadImageFromBytes(brickBluePNG)
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

	return &Images{
		Paddle:        paddle,
		Brick:         brick,
		BrickGreen:    brickGreen,
		BrickBlue:     brickBlue,
		BrickColumbia: brickColumbia,
		BrickSupreme:  brickSupreme,
	}, nil
}

// loadImageFromBytes converts embedded PNG bytes to ebiten.Image
func loadImageFromBytes(data []byte) (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

// GetBrickImage returns the appropriate brick sprite based on brick type
func (imgs *Images) GetBrickImage(brickType entities.BrickType) *ebiten.Image {
	switch brickType {
	case entities.BrickTypeGreen:
		return imgs.BrickGreen
	case entities.BrickTypeBlue:
		return imgs.BrickBlue
	case entities.BrickTypeColumbia:
		return imgs.BrickColumbia
	case entities.BrickTypeSupreme:
		return imgs.BrickSupreme
	case entities.BrickTypeDefault:
		return imgs.Brick
	default:
		return imgs.Brick // fallback to default brick sprite
	}
}
