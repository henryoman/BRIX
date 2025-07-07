package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"BRIX/config"
	"BRIX/game"
)

func main() {
	// Calculate an initial window size that fits on the current monitor.
	const baseW, baseH = 1440, 1080

	if mon := ebiten.Monitor(); mon != nil {
		mw, mh := mon.Size()
		const uiPadding = 80 // Rough allowance for OS menu + title bars
		usableH := float64(mh - uiPadding)
		if usableH <= 0 {
			usableH = float64(mh)
		}
		scale := math.Min(float64(mw)/baseW, usableH/float64(baseH))
		if scale < 1 {
			ebiten.SetWindowSize(int(float64(baseW)*scale), int(float64(baseH)*scale))
		} else {
			ebiten.SetWindowSize(baseW, baseH)
		}
	} else {
		// Fallback if monitor info unavailable
		ebiten.SetWindowSize(baseW, baseH)
	}

	// Allow the user to resize freely.
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	ebiten.SetWindowTitle("Brick Breaker")

	// Load brick & scoring configs
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	g := game.NewGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
