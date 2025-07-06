package main

import (
	"log"
	"math/rand"
	"time"

	"run-boy-run/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("ROAD ADVENTURE")

	rand.Seed(time.Now().UnixNano())

	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}