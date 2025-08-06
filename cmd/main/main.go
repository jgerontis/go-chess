package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jgerontis/go-chess/gui"
	"github.com/jgerontis/go-chess/internal/chess"
)

func main() {
	log.Println("Starting Go Chess")

	// Command line flags
	var (
		debugMode   = flag.Bool("debug", false, "Skip to debug mode")
		humanVsAI   = flag.Bool("human-vs-ai", false, "Skip to human vs AI mode")
		aiVsAI      = flag.Bool("ai-vs-ai", false, "Skip to AI vs AI mode")
		fenString   = flag.String("fen", chess.START_FEN, "FEN string for initial position")
	)
	flag.Parse()

	// Determine starting mode
	var startMode gui.GameMode = gui.MainMenu
	if *debugMode {
		startMode = gui.DebugMode
		log.Println("Starting in debug mode")
	} else if *humanVsAI {
		startMode = gui.HumanVsAI
		log.Println("Starting in human vs AI mode")
	} else if *aiVsAI {
		startMode = gui.AIvsAI
		log.Println("Starting in AI vs AI mode")
	} else {
		log.Println("Starting with main menu")
	}

	// Handle legacy positional FEN argument for backwards compatibility
	if len(flag.Args()) > 0 {
		*fenString = flag.Args()[0]
		if startMode == gui.MainMenu {
			startMode = gui.DebugMode // Default to debug mode if FEN provided without flag
		}
		log.Println("Using provided FEN string:", *fenString)
	}

	app := gui.NewApp(startMode, *fenString)

	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Go Chess")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(160)

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal("Failed to run game, error: ", err)
	}
}
