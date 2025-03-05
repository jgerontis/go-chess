package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// StartUCI listens for commands and processes them.
func StartUCI() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("id name GoChess")
	fmt.Println("id author JoshGerontis")
	fmt.Println("uciok")

	for scanner.Scan() {
		input := scanner.Text()
		handleUCICommand(input)
	}
}

func handleUCICommand(input string) {
	if input == "uci" {
		fmt.Println("id name GoChess")
		fmt.Println("id author JoshGerontis")
		fmt.Println("uciok")
	} else if strings.HasPrefix(input, "position") {
		// Handle setting up the board position
		fmt.Println("Position received:", input)
	} else if strings.HasPrefix(input, "go") {
		// Engine needs to compute the best move
		fmt.Println("bestmove e2e4") // Placeholder
	} else if input == "quit" {
		os.Exit(0)
	}
}
