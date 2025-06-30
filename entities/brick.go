package entities

import (
	"image/color"
)

const (
	BrickCols = 12
	BrickRows = 10
)

// BrickType represents the type of brick with direct sprite mapping
type BrickType string

const (
	BrickTypeDefault  BrickType = "default"  // brick.png - red/orange/yellow
	BrickTypeGreen    BrickType = "green"    // brick-green.png
	BrickTypeBlue     BrickType = "blue"     // brick-blue.png
	BrickTypeColumbia BrickType = "columbia" // brick-columbia.png - white
	BrickTypeSupreme  BrickType = "supreme"  // brick-supreme.png - pink/purple
)

// Brick represents a single brick in the level
type Brick struct {
	x, y      int       // grid position
	brickType BrickType // type of brick (maps directly to sprite)
	hits      int       // hits required to destroy
	active    bool      // whether brick is still active

	// Level-specific sizing (set when brick is created)
	width, height      int
	spacingX, spacingY int

	// Field bounds for smart centering (set when brick is created)
	fieldMinX, fieldMaxX int
}

// LevelBrick represents a brick definition from level data
type LevelBrick struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	BrickType string `json:"bricktype"` // brick type name
	Hits      int    `json:"hits"`
}

// NewBrick creates a new brick at the specified grid position with custom sizing
func NewBrick(x, y int, brickType BrickType, hits int, width, height, spacingX, spacingY int) *Brick {
	return &Brick{
		x:         x,
		y:         y,
		brickType: brickType,
		hits:      hits,
		active:    true,
		width:     width,
		height:    height,
		spacingX:  spacingX,
		spacingY:  spacingY,
		fieldMinX: 0,
		fieldMaxX: 7, // default field bounds
	}
}

// NewBrickFromLevel creates a brick from level data with level's sizing
func NewBrickFromLevel(levelBrick LevelBrick, width, height, spacingX, spacingY int) *Brick {
	return &Brick{
		x:         levelBrick.X,
		y:         levelBrick.Y,
		brickType: ParseBrickType(levelBrick.BrickType),
		hits:      levelBrick.Hits,
		active:    true,
		width:     width,
		height:    height,
		spacingX:  spacingX,
		spacingY:  spacingY,
		fieldMinX: 0,
		fieldMaxX: 7, // default field bounds
	}
}

// NewBrickFromLevelWithBounds creates a brick from level data with calculated field bounds for smart centering
func NewBrickFromLevelWithBounds(levelBrick LevelBrick, width, height, spacingX, spacingY, fieldMinX, fieldMaxX int) *Brick {
	return &Brick{
		x:         levelBrick.X,
		y:         levelBrick.Y,
		brickType: ParseBrickType(levelBrick.BrickType),
		hits:      levelBrick.Hits,
		active:    true,
		width:     width,
		height:    height,
		spacingX:  spacingX,
		spacingY:  spacingY,
		fieldMinX: fieldMinX,
		fieldMaxX: fieldMaxX,
	}
}

// ParseBrickType converts a string to a BrickType (for backward compatibility)
func ParseBrickType(typeStr string) BrickType {
	switch typeStr {
	case "green":
		return BrickTypeGreen
	case "blue":
		return BrickTypeBlue
	case "columbia", "white":
		return BrickTypeColumbia
	case "supreme", "pink", "purple":
		return BrickTypeSupreme
	case "default", "red", "orange", "yellow":
		return BrickTypeDefault
	default:
		return BrickTypeDefault
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

// Type returns the brick's type
func (b *Brick) Type() BrickType {
	return b.brickType
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

// GetScreenPosition returns the pixel position of the brick on screen with smart centering
func (b *Brick) GetScreenPosition() (float64, float64) {
	const screenWidth = 720

	// Calculate the total field width based on actual field bounds
	fieldWidthInBricks := b.fieldMaxX - b.fieldMinX + 1
	totalFieldWidth := float64(fieldWidthInBricks*b.width + (fieldWidthInBricks-1)*b.spacingX)

	// Center the field on screen
	fieldStartX := (screenWidth - totalFieldWidth) / 2

	// Calculate this brick's position relative to the field start
	brickOffsetFromFieldStart := float64((b.x - b.fieldMinX) * (b.width + b.spacingX))

	screenX := fieldStartX + brickOffsetFromFieldStart
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

// GetDisplayColor returns the appropriate color for rendering based on the brick type
func (b *Brick) GetDisplayColor() color.Color {
	switch b.brickType {
	case BrickTypeDefault:
		return color.RGBA{255, 100, 100, 255} // red
	case BrickTypeGreen:
		return color.RGBA{100, 255, 100, 255} // green
	case BrickTypeBlue:
		return color.RGBA{100, 150, 255, 255} // blue
	case BrickTypeColumbia:
		return color.RGBA{255, 255, 255, 255} // white
	case BrickTypeSupreme:
		return color.RGBA{255, 100, 255, 255} // pink/purple
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
