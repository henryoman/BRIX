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

	return &level, nil
}

// CreateDefaultLevel creates a simple default level for fallback
func CreateDefaultLevel() *Level {
	return &Level{
		Name: "Default Level",
		Bricks: []entities.LevelBrick{
			{X: 4, Y: 2, BrickType: "default", Hits: 1},
			{X: 5, Y: 2, BrickType: "default", Hits: 1},
			{X: 6, Y: 2, BrickType: "default", Hits: 1},
			{X: 7, Y: 2, BrickType: "default", Hits: 1},
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
	}

	return nil
}
