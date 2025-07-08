package levels

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"BRIX/entities"
)

// Level represents a complete level configuration
type Level struct {
	Name string `json:"name"`

	// Grid-based format (legacy)
	BrickWidth    int `json:"brick_width,omitempty"`     // configurable brick width
	BrickHeight   int `json:"brick_height,omitempty"`    // configurable brick height
	BrickSpacingX int `json:"brick_spacing_x,omitempty"` // horizontal gap
	BrickSpacingY int `json:"brick_spacing_y,omitempty"` // vertical gap

	// Pixel-perfect format (new)
	UsePixelPositioning bool `json:"use_pixel_positioning,omitempty"` // flag to indicate pixel format
	DefaultBrickWidth   int  `json:"default_brick_width,omitempty"`   // default width if not specified per brick
	DefaultBrickHeight  int  `json:"default_brick_height,omitempty"`  // default height if not specified per brick

	BallSpeed float64               `json:"ball_speed"` // ball speed in pixels per second
	Bricks    []entities.LevelBrick `json:"bricks"`
}

// LoadLevel loads a level from a JSON file
func LoadLevel(levelNum int) (*Level, error) {
	filename := filepath.Join("levels", fmt.Sprintf("level%d.json", levelNum))
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read level file %s: %v", filename, err)
	}

	var level Level
	err = json.Unmarshal(data, &level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse level file %s: %v", filename, err)
	}

	// Auto-detect pixel positioning format
	if level.UsePixelPositioning || isPixelFormat(&level) {
		level.UsePixelPositioning = true
		// Set reasonable defaults for pixel format
		if level.DefaultBrickWidth == 0 {
			level.DefaultBrickWidth = 150 // default brick width
		}
		if level.DefaultBrickHeight == 0 {
			level.DefaultBrickHeight = 60 // default brick height
		}
	} else {
		// Legacy grid-based format - apply auto-fit if needed
		AutoFitLevel(&level)
	}

	// Validate the level
	if err := ValidateLevel(&level); err != nil {
		return nil, fmt.Errorf("level validation failed for %s: %v", filename, err)
	}

	return &level, nil
}

// isPixelFormat auto-detects if this is a pixel-perfect format based on the data
func isPixelFormat(level *Level) bool {
	// Check if any brick has "type" field (new format) or "pixel_x"/"pixel_y" fields
	for _, brick := range level.Bricks {
		if brick.Type != "" || brick.PixelX != 0 || brick.PixelY != 0 {
			return true
		}
	}
	// Check if grid sizing fields are missing (indicating pixel format)
	return level.BrickWidth == 0 && level.BrickHeight == 0
}

// CreateDefaultLevel creates a simple default level for fallback
func CreateDefaultLevel() *Level {
	return &Level{
		Name:          "Default Level",
		BrickWidth:    150,
		BrickHeight:   60,
		BrickSpacingX: 25,
		BrickSpacingY: 30,
		BallSpeed:     400,
		Bricks: []entities.LevelBrick{
			{X: 4, Y: 2, BrickType: "standard", Hits: 1},
			{X: 5, Y: 2, BrickType: "standard", Hits: 1},
			{X: 6, Y: 2, BrickType: "standard", Hits: 1},
			{X: 7, Y: 2, BrickType: "standard", Hits: 1},
		},
	}
}

// ValidateLevel checks if a level has valid brick configurations
func ValidateLevel(level *Level) error {
	if level.Name == "" {
		return fmt.Errorf("level must have a name")
	}

	if len(level.Bricks) == 0 {
		return fmt.Errorf("level must have at least one brick")
	}

	for i, brick := range level.Bricks {
		if brick.X < 0 || brick.X >= entities.BrickCols {
			return fmt.Errorf("brick %d has invalid X position: %d", i, brick.X)
		}
		if brick.Y < 0 || brick.Y >= entities.BrickRows {
			return fmt.Errorf("brick %d has invalid Y position: %d", i, brick.Y)
		}
		if brick.Hits <= 0 {
			return fmt.Errorf("brick %d must have positive hits: %d", i, brick.Hits)
		}

		// --- Gameplay-area bounds checks ---
		// Determine horizontal field span (min/max X) to calculate pixel positions.
	}

	// Calculate brick field bounds (min & max X) for width validation.
	minX := level.Bricks[0].X
	maxX := level.Bricks[0].X
	minY := level.Bricks[0].Y
	maxY := level.Bricks[0].Y
	for _, b := range level.Bricks {
		if b.X < minX {
			minX = b.X
		}
		if b.X > maxX {
			maxX = b.X
		}
		if b.Y < minY {
			minY = b.Y
		}
		if b.Y > maxY {
			maxY = b.Y
		}
	}

	// --- Pixel-exact bounds validation ---
	fieldColumns := maxX - minX + 1
	fieldWidthPx := fieldColumns*level.BrickWidth + (fieldColumns-1)*level.BrickSpacingX
	if float64(fieldWidthPx) > entities.GameAreaWidth {
		return fmt.Errorf("brick field width (%d px) exceeds gameplay width (%.0f px)", fieldWidthPx, entities.GameAreaWidth)
	}

	// Horizontal centering offset (same calculation used by entities.Brick).
	fieldStartX := entities.GameAreaLeft + (entities.GameAreaWidth-float64(fieldWidthPx))/2
	fieldEndX := fieldStartX + float64(fieldWidthPx)
	if fieldStartX < entities.GameAreaLeft || fieldEndX > entities.GameAreaRight {
		return fmt.Errorf("brick field would render outside horizontal gameplay bounds (start=%.0f, end=%.0f)", fieldStartX, fieldEndX)
	}

	// Vertical bounds: emulate entities.Brick.GetScreenPosition logic.
	verticalStride := level.BrickHeight + level.BrickSpacingY
	topMostY := entities.GameAreaTop + float64(minY*verticalStride)
	bottomMostY := entities.GameAreaTop + float64(maxY*verticalStride+level.BrickHeight)
	if topMostY < entities.GameAreaTop {
		return fmt.Errorf("top bricks would render above gameplay area (y=%.0f)", topMostY)
	}
	if bottomMostY > entities.GameAreaBottom {
		return fmt.Errorf("bottom bricks would render below gameplay area (y=%.0f)", bottomMostY)
	}

	return nil
}

// AutoFitLevel modifies brick sizes so the field fits horizontally in GameAreaWidth
// without altering brick layout. It leaves spacing unchanged. If the bricks already
// fit, the level is returned untouched. Only brick_width (and proportionally
// brick_height, preserving aspect ratio) are reduced.
func AutoFitLevel(level *Level) {
	if len(level.Bricks) == 0 || level.BrickWidth <= 0 {
		return
	}

	// Determine min/max X to know number of columns
	minX := level.Bricks[0].X
	maxX := level.Bricks[0].X
	for _, b := range level.Bricks {
		if b.X < minX {
			minX = b.X
		}
		if b.X > maxX {
			maxX = b.X
		}
	}
	cols := maxX - minX + 1
	if cols <= 0 {
		return
	}

	// Compute current required width
	currentWidth := cols*level.BrickWidth + (cols-1)*level.BrickSpacingX
	if float64(currentWidth) <= entities.GameAreaWidth {
		return // already fits
	}

	// Compute new brick width that will fit exactly
	available := entities.GameAreaWidth - float64((cols-1)*level.BrickSpacingX)
	if available <= 0 {
		return // cannot fit, leave as is; validation will catch
	}
	newWidth := available / float64(cols)
	if newWidth <= 0 {
		return
	}

	// Scale height proportionally
	ratio := float64(level.BrickHeight) / float64(level.BrickWidth)
	level.BrickWidth = int(newWidth)
	level.BrickHeight = int(newWidth * ratio)
}
