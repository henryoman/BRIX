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
	BrickTypeStandard BrickType = "standard" // brick-standard.png
	BrickTypeTusi     BrickType = "tusi"     // brick-tusi.png
	BrickTypeWeed     BrickType = "weed"     // brick-weed.png
	BrickTypeColumbia BrickType = "columbia" // brick-columbia.png
	BrickTypeSupreme  BrickType = "supreme"  // brick-supreme.png
)

// Brick represents a single brick in the level
type Brick struct {
	// Grid positioning (legacy)
	x, y int // grid position

	// Pixel positioning (new)
	pixelX, pixelY   float64 // absolute pixel position
	usePixelPosition bool    // whether to use pixel or grid positioning

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
	// Grid-based positioning (legacy)
	X         int    `json:"x"`
	Y         int    `json:"y"`
	BrickType string `json:"bricktype,omitempty"` // legacy field name

	// Pixel-perfect positioning (new format)
	PixelX int    `json:"pixel_x,omitempty"` // absolute pixel X position
	PixelY int    `json:"pixel_y,omitempty"` // absolute pixel Y position
	Type   string `json:"type,omitempty"`    // unified type field

	// Common fields
	Hits   int `json:"hits"`
	Width  int `json:"width,omitempty"`  // per-brick width override
	Height int `json:"height,omitempty"` // per-brick height override
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

// NewBrickPixelPosition creates a new brick at the specified pixel position
func NewBrickPixelPosition(pixelX, pixelY float64, brickType BrickType, hits, width, height int) *Brick {
	return &Brick{
		pixelX:           pixelX,
		pixelY:           pixelY,
		usePixelPosition: true,
		brickType:        brickType,
		hits:             hits,
		active:           true,
		width:            width,
		height:           height,
	}
}

// NewBrickFromLevelPixel creates a brick from pixel-perfect level data
func NewBrickFromLevelPixel(levelBrick LevelBrick, defaultWidth, defaultHeight int) *Brick {
	// Determine type from either "type" or "bricktype" field
	brickTypeStr := levelBrick.Type
	if brickTypeStr == "" {
		brickTypeStr = levelBrick.BrickType
	}

	// Use per-brick dimensions or defaults
	width := levelBrick.Width
	if width == 0 {
		width = defaultWidth
	}
	height := levelBrick.Height
	if height == 0 {
		height = defaultHeight
	}

	// Use pixel position if available, otherwise convert from grid
	var pixelX, pixelY float64
	if levelBrick.PixelX != 0 || levelBrick.PixelY != 0 {
		pixelX = float64(levelBrick.PixelX)
		pixelY = float64(levelBrick.PixelY)
	} else {
		// Convert from grid coordinates using default spacing
		pixelX = float64(levelBrick.X)
		pixelY = float64(levelBrick.Y)
	}

	return &Brick{
		pixelX:           pixelX,
		pixelY:           pixelY,
		usePixelPosition: true,
		brickType:        ParseBrickType(brickTypeStr),
		hits:             levelBrick.Hits,
		active:           true,
		width:            width,
		height:           height,
	}
}

// ParseBrickType converts a string to a BrickType (for backward compatibility)
func ParseBrickType(typeStr string) BrickType {
	switch typeStr {
	case "columbia", "white":
		return BrickTypeColumbia
	case "supreme", "pink", "purple":
		return BrickTypeSupreme
	case "tusi", "green":
		return BrickTypeTusi
	case "weed", "blue", "cyan":
		return BrickTypeWeed
	case "standard", "default", "red", "orange", "yellow":
		return BrickTypeStandard
	default:
		return BrickTypeStandard
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
	// If using pixel positioning, return absolute position within game area
	if b.usePixelPosition {
		return GameAreaLeft + b.pixelX, GameAreaTop + b.pixelY
	}

	// Legacy grid-based positioning with smart centering
	// Calculate horizontal centering offset based on field bounds.
	cols := b.fieldMaxX - b.fieldMinX + 1
	if cols <= 0 {
		cols = 1
	}

	fieldWidthPx := cols*b.width + (cols-1)*b.spacingX
	// Ensure we work with float64 for sub-pixel precision when scaling.
	offsetX := (GameAreaWidth - float64(fieldWidthPx)) / 2.0

	// Position within the grid relative to the leftmost brick (fieldMinX).
	dx := b.x - b.fieldMinX

	screenX := GameAreaLeft + offsetX + float64(dx*(b.width+b.spacingX))
	screenY := GameAreaTop + float64(b.y*(b.height+b.spacingY))

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
	case BrickTypeStandard:
		return color.RGBA{255, 100, 100, 255} // red-ish standard
	case BrickTypeTusi:
		return color.RGBA{255, 140, 255, 255} // tusi pinkish
	case BrickTypeWeed:
		return color.RGBA{100, 255, 100, 255} // green weed
	case BrickTypeColumbia:
		return color.RGBA{255, 255, 255, 255} // white
	case BrickTypeSupreme:
		return color.RGBA{255, 100, 255, 255} // pink/purple supreme
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
