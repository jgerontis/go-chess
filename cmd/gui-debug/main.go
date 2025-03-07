package main

// the user just clicks for both black and white, no AI

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jgerontis/go-chess/gui"
	"github.com/jgerontis/go-chess/internal/chess"
)

func main() {
	// start the game
	log.Println("Starting Go Chess")

	var fenString string
	if len(os.Args) > 1 {
		log.Println("Using provided FEN string: ", fenString)
		fenString = os.Args[1]
	} else {
		log.Println("No FEN string provided, starting with default position")
		fenString = chess.START_FEN
	}

	game := gui.NewGame(fenString)

	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Go Chess")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(160)

	// ebiten.SetScreenClearedEveryFrame(false)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Failed to run game, error: ", err)
	}
}
