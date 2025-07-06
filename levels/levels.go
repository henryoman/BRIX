package levels

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"brick-breaker/entities"
)

// Level represents a complete level configuration
type Level struct {
	Name          string                `json:"name"`
	BrickWidth    int                   `json:"brick_width"`     // configurable brick width
	BrickHeight   int                   `json:"brick_height"`    // configurable brick height
	BrickSpacingX int                   `json:"brick_spacing_x"` // horizontal gap
	BrickSpacingY int                   `json:"brick_spacing_y"` // vertical gap
	BallSpeed     float64               `json:"ball_speed"`      // ball speed in pixels per second
	Bricks        []entities.LevelBrick `json:"bricks"`
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

	// Attempt to auto-fit the level into the 700×500 gameplay area if necessary.
	AutoFitLevel(&level)

	// Validate again after potential adjustments.
	if err := ValidateLevel(&level); err != nil {
		return nil, fmt.Errorf("level validation failed for %s: %v", filename, err)
	}

	return &level, nil
}

// CreateDefaultLevel creates a simple default level for fallback
func CreateDefaultLevel() *Level {
	return &Level{
		Name:          "Default Level",
		BrickWidth:    75,
		BrickHeight:   30,
		BrickSpacingX: 25,
		BrickSpacingY: 30,
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
	if fieldWidthPx > entities.GameAreaWidth {
		return fmt.Errorf("brick field width (%d px) exceeds gameplay width (%d px)", fieldWidthPx, entities.GameAreaWidth)
	}

	// Horizontal centering offset (same calculation used by entities.Brick).
	fieldStartX := entities.GameAreaLeft + (entities.GameAreaWidth-fieldWidthPx)/2
	fieldEndX := fieldStartX + fieldWidthPx
	if fieldStartX < entities.GameAreaLeft || fieldEndX > entities.GameAreaRight {
		return fmt.Errorf("brick field would render outside horizontal gameplay bounds (start=%d, end=%d)", fieldStartX, fieldEndX)
	}

	// Vertical bounds: emulate entities.Brick.GetScreenPosition logic.
	verticalStride := level.BrickHeight + level.BrickSpacingY
	topMostY := entities.GameAreaTop - 50 + minY*verticalStride
	bottomMostY := entities.GameAreaTop - 50 + maxY*verticalStride + level.BrickHeight
	if topMostY < entities.GameAreaTop {
		return fmt.Errorf("top bricks would render above gameplay area (y=%d)", topMostY)
	}
	if bottomMostY > entities.GameAreaBottom {
		return fmt.Errorf("bottom bricks would render below gameplay area (y=%d)", bottomMostY)
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
	if currentWidth <= entities.GameAreaWidth {
		return // already fits
	}

	// Compute new brick width that will fit exactly
	available := entities.GameAreaWidth - (cols-1)*level.BrickSpacingX
	if available <= 0 {
		return // cannot fit, leave as is; validation will catch
	}
	newWidth := available / cols
	if newWidth <= 0 {
		return
	}

	// Scale height proportionally
	ratio := float64(level.BrickHeight) / float64(level.BrickWidth)
	level.BrickWidth = newWidth
	level.BrickHeight = int(float64(newWidth) * ratio)
}
