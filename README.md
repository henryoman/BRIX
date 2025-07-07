<div align="center">
  <img src="logo.png" alt="Logo" width="600"/>
</div>

<h3 align="center">BRIX</h3>

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/your-repo)](https://goreportcard.com/report/github.com/henryoman/brick-breaker)
[![Go Reference](https://pkg.go.dev/badge/github.com/your-username/your-repo.svg)](https://pkg.go.dev/github.com/your-username/your-repo)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Ebitengine](https://img.shields.io/badge/Made%20with-Ebitengine-blue.svg)](https://ebitengine.org/)

</div>

---

A modern brick breaker game built with Go and Ebitengine, featuring a flexible JSON-based level system.

## Features

- **Circular Ball**: Smooth, realistic ball physics with proper circular collision
- **Configurable Levels**: Easy-to-create JSON level files
- **Multiple Brick Types**: Different colors and hit requirements
- **Score System**: Points for hitting paddles and destroying bricks
- **Modern Graphics**: Vector-based rendering with Ebitengine

## Controls

- **Arrow Keys** or **A/D Keys**: Move paddle left/right
- The game automatically starts when you run it

## Level System

### Creating New Levels

Levels are stored as JSON files in the `levels/` directory. Each level file should be named `levelX.json` where X is the level number.

### Level Format

```json
{
  "name": "Level Name",
  "bricks": [
    {"x": 0, "y": 1, "color": "red", "hits": 1},
    {"x": 1, "y": 1, "color": "blue", "hits": 2}
  ]
}
```

### Grid System

- **Screen Width**: 720 pixels
- **Brick Width**: 60 pixels (12 bricks fit horizontally)
- **Brick Height**: 20 pixels
- **Grid Coordinates**: 
  - X: 0-11 (12 columns)
  - Y: 0-9 (10 rows maximum)

### Brick Properties

- **x, y**: Grid position (0-based)
- **color**: Brick color (affects appearance)
- **hits**: Number of hits required to destroy the brick

### Available Colors

- `red`: Red bricks (1 hit typically)
- `orange`: Orange bricks 
- `yellow`: Yellow bricks
- `green`: Green bricks
- `blue`: Blue bricks (often 2+ hits)
- `purple`: Purple bricks (often 3+ hits)
- `pink`: Pink bricks

### Scoring

- **Paddle Hit**: 10 points
- **Brick Hit**: 25 points
- **Brick Destroyed**: 50 points

## Building and Running

```bash
# Build the game
go build -o brick-breaker main.go

# Run the game
./brick-breaker
```

## Example Levels

### Level 1 - Easy Start
Simple horizontal rows with single-hit bricks.

### Level 2 - Getting Harder  
Mixed brick types with some requiring multiple hits.

## Adding New Levels

1. Create a new JSON file in the `levels/` directory
2. Name it `levelX.json` where X is the next level number
3. Use the grid system to place bricks logically
4. Test different color combinations and hit requirements
5. The game will automatically load the new level

## Technical Details

- Built with **Go 1.24+**
- Uses **Ebitengine v2.8+** for graphics
- Vector-based rendering for smooth graphics
- JSON-based configuration for maximum flexibility
- Modular design for easy expansion

## Tips for Level Design

1. **Start Simple**: Begin with single-hit bricks in basic patterns
2. **Add Complexity**: Introduce multi-hit bricks gradually
3. **Use Colors Meaningfully**: Different colors can indicate difficulty
4. **Leave Gaps**: Strategic gaps make levels more interesting
5. **Test Thoroughly**: Play-test your levels to ensure they're fun and fair 