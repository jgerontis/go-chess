package main

// the user just clicks for both black and white, no AI

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jgerontis/go-chess/gui"
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
		fenString = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	}

	game := gui.NewGame(fenString)

	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Go Chess")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game.PrintBoard()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Failed to run game, error: ", err)
	}
}
