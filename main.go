package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"brick-breaker/config"
	"brick-breaker/game"
)

func main() {
	ebiten.SetWindowSize(1440, 1080)
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
