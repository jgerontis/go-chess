package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jgerontis/go-chess/internal/chess"
	"golang.org/x/exp/slog"
)

func main() {
	// set up slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// start the game
	slog.Log(nil, slog.LevelInfo, "Starting Go Chess")

	var game *chess.Game
	var fenString string
	if len(os.Args) > 1 {
		slog.Info("Using provided FEN string", "fen", fenString)
		fenString = os.Args[1]
	} else {
		slog.Info("No FEN string provided, starting with default position")
		fenString = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
	}
	game = chess.NewGame(fenString)

	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Go Chess")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game.PrintBoard()

	if err := ebiten.RunGame(game); err != nil {
		slog.Error("Failed to run game", "error", err)
		log.Fatal(err)
	}
}
