package entities

import (
	"image/color"
)

const (
	BrickCols = 12
	BrickRows = 10
)

// Brick represents a single brick in the level
type Brick struct {
	x, y   int    // grid position
	color  string // color name
	hits   int    // hits required to destroy
	active bool   // whether brick is still active

	// Level-specific sizing (set when brick is created)
	width, height      int
	spacingX, spacingY int
}

// LevelBrick represents a brick definition from level data
type LevelBrick struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"`
	Hits  int    `json:"hits"`
}

// NewBrick creates a new brick at the specified grid position with custom sizing
func NewBrick(x, y int, color string, hits int, width, height, spacingX, spacingY int) *Brick {
	return &Brick{
		x:        x,
		y:        y,
		color:    color,
		hits:     hits,
		active:   true,
		width:    width,
		height:   height,
		spacingX: spacingX,
		spacingY: spacingY,
	}
}

// NewBrickFromLevel creates a brick from level data with level's sizing
func NewBrickFromLevel(levelBrick LevelBrick, width, height, spacingX, spacingY int) *Brick {
	return &Brick{
		x:        levelBrick.X,
		y:        levelBrick.Y,
		color:    levelBrick.Color,
		hits:     levelBrick.Hits,
		active:   true,
		width:    width,
		height:   height,
		spacingX: spacingX,
		spacingY: spacingY,
	}
}

// X returns the grid X position
func (b *Brick) X() int {
	return b.x
}

// Y returns the grid Y position
func (b *Brick) Y() int {
	return b.y
}

// Color returns the brick's color name
func (b *Brick) Color() string {
	return b.color
}

// Hits returns the remaining hits needed to destroy the brick
func (b *Brick) Hits() int {
	return b.hits
}

// IsActive returns whether the brick is still active
func (b *Brick) IsActive() bool {
	return b.active
}

// Hit reduces the brick's hit count and deactivates it if necessary
func (b *Brick) Hit() bool {
	if !b.active {
		return false
	}

	b.hits--
	if b.hits <= 0 {
		b.active = false
		return true // brick destroyed
	}
	return false // brick damaged but not destroyed
}

// GetScreenPosition returns the pixel position of the brick on screen with gaps
func (b *Brick) GetScreenPosition() (float64, float64) {
	screenX := float64(b.x * (b.width + b.spacingX))
	screenY := float64(HUDHeight + b.y*(b.height+b.spacingY))
	return screenX, screenY
}

// GetBounds returns the brick's bounding box for collision detection
func (b *Brick) GetBounds() (left, top, right, bottom float64) {
	screenX, screenY := b.GetScreenPosition()
	left = screenX
	right = screenX + float64(b.width)
	top = screenY
	bottom = screenY + float64(b.height)
	return
}

// GetColor returns the appropriate color for rendering based on the color name
func (b *Brick) GetColor() color.Color {
	switch b.color {
	case "red":
		return color.RGBA{255, 100, 100, 255}
	case "orange":
		return color.RGBA{255, 165, 0, 255}
	case "yellow":
		return color.RGBA{255, 255, 100, 255}
	case "green":
		return color.RGBA{100, 255, 100, 255}
	case "blue":
		return color.RGBA{100, 150, 255, 255}
	case "purple":
		return color.RGBA{200, 100, 255, 255}
	default:
		return color.RGBA{200, 200, 200, 255} // gray default
	}
}

// Width returns the brick's width
func (b *Brick) Width() int {
	return b.width
}

// Height returns the brick's height
func (b *Brick) Height() int {
	return b.height
}
