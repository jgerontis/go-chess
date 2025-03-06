package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jgerontis/go-chess/internal/chess"
)

/*
	UCI (Universal Chess Interface) is a protocol for communicating with chess engines.
	The engine will receive commands from a GUI and respond accordingly.
	An engine should not maintain a state between commands.
	Further reading: https://www.chessprogramming.org/UCI
	Full UCI spec: https://gist.github.com/DOBRO/2592c6dad754ba67e6dcaec8c90165bf
*/

// runs the UCI engine loop
func StartEngine() {
	scanner := bufio.NewScanner(os.Stdin)
	board := chess.NewBoard()

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "uci":
			fmt.Println("id name GoChessJG")
			fmt.Println("id author JoshGerontis")
			fmt.Println("uciok")

		case "isready":
			fmt.Println("readyok")

		case "position":
			board.LoadFEN(parts[2])

		case "go":
			bestMove := FindBestMove(board)
			fmt.Println("bestmove", bestMove)

		case "quit":
			return
		}
	}
}
