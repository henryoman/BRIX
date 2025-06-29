package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"brick-breaker/game"
)

func main() {
	ebiten.SetWindowSize(720, 800)
	ebiten.SetWindowTitle("Brick Breaker")

	g := game.NewGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
