package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// BrickTypeCfg describes a single brick type loaded from brick_types.json.
type BrickTypeCfg struct {
	Name        string  `json:"name"`
	Sprite      string  `json:"sprite"`
	Hits        int     `json:"hits"`
	SpeedFactor float64 `json:"speedFactor"`
	PowerUp     string  `json:"powerUp"`
}

// BrickTypes maps a brick shorthand / name to its config.
// Example keys: "default", "green", "blue".
type BrickTypes map[string]BrickTypeCfg

// PointsByLives stores points keyed by remaining lives ("3","2","1").
type PointsByLives map[string]int

// ScoringConfig controls point economy and now supports per-lives values.
type ScoringConfig struct {
	PaddleHit    PointsByLives            `json:"paddleHit"`
	LifeBonus    int                      `json:"lifeBonus"`
	BrickHit     map[string]PointsByLives `json:"brickHit"`
	BrickDestroy map[string]PointsByLives `json:"brickDestroy"`
	PowerUp      map[string]PointsByLives `json:"powerUp"`
}

var (
	// Brick holds the runtime-available brick palette.
	Brick BrickTypes
	// Score holds the scoring rules. Safe to read concurrently after Load returns.
	Score ScoringConfig
)

// Load reads brick_types.json and scoring.json into memory. Call this once at program start.
func Load() error {
	if err := loadBrickTypes("config/brick_types.json"); err != nil {
		return fmt.Errorf("load brick types: %w", err)
	}
	if err := loadScoring("config/scoring.json"); err != nil {
		return fmt.Errorf("load scoring: %w", err)
	}
	return nil
}

func loadBrickTypes(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var m BrickTypes
	if err := json.Unmarshal(raw, &m); err != nil {
		return err
	}
	if len(m) == 0 {
		return fmt.Errorf("brick_types.json contains no entries")
	}
	Brick = m
	return nil
}

func loadScoring(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var s ScoringConfig
	if err := json.Unmarshal(raw, &s); err != nil {
		return err
	}
	Score = s
	return nil
}
