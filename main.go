package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"BRIX/config"
	"BRIX/game"
)

func main() {
	ebiten.SetWindowSize(1440, 1080)
	ebiten.SetWindowTitle("BRIX - Brick Breaker Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Set window size limits to common 4:3 aspect ratios
	// This constrains the OS-level resize to reduce snap-back
	minW, minH := 320, 240   // 4:3 minimum
	maxW, maxH := 2560, 1920 // 4:3 maximum
	ebiten.SetWindowSizeLimits(minW, minH, maxW, maxH)

	// Load brick & scoring configs
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	g := game.NewGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
