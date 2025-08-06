package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jgerontis/go-chess/internal/uci"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  TCP mode:  go run ./cmd/engine-cli tcp <address>")
		fmt.Println("  Exec mode: go run ./cmd/engine-cli exec <engine_path> [args...]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run ./cmd/engine-cli tcp localhost:8080")
		fmt.Println("  go run ./cmd/engine-cli exec stockfish")
		fmt.Println("  go run ./cmd/engine-cli exec ./my-engine --threads 4")
		os.Exit(1)
	}

	mode := os.Args[1]
	var client *uci.Client
	var err error

	switch mode {
	case "tcp":
		if len(os.Args) < 3 {
			log.Fatal("TCP mode requires an address (e.g., localhost:8080)")
		}
		address := os.Args[2]
		fmt.Printf("Connecting to UCI server at %s...\n", address)
		client, err = uci.NewClientTCP(address)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connected to TCP server!")

	case "exec":
		if len(os.Args) < 3 {
			log.Fatal("Exec mode requires an engine path")
		}
		enginePath := os.Args[2]
		args := []string{}
		if len(os.Args) > 3 {
			args = os.Args[3:]
		}
		fmt.Printf("Starting UCI engine: %s %s\n", enginePath, strings.Join(args, " "))
		client, err = uci.NewClientExec(enginePath, args...)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Engine process started!")

	default:
		log.Fatalf("Unknown mode: %s. Use 'tcp' or 'exec'", mode)
	}

	defer client.Close()

	fmt.Println("UCI client ready!")
	fmt.Println("Type 'help' for available commands, 'quit' to exit")
	fmt.Println("All UCI commands are supported (uci, position, go, stop, etc.)")
	fmt.Println()

	// Start response listener
	responseChan := make(chan string, 100)
	errorChan := make(chan error, 10)
	client.StartResponseListener(responseChan, errorChan)

	// Handle responses in a separate goroutine
	go func() {
		for {
			select {
			case response, ok := <-responseChan:
				if !ok {
					return
				}
				fmt.Printf("< %s\n", response)
			case err, ok := <-errorChan:
				if !ok {
					return
				}
				log.Printf("Error reading from engine: %v", err)
				return
			}
		}
	}()

	// Interactive command loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		command := strings.TrimSpace(scanner.Text())
		if command == "" {
			continue
		}

		// Handle local commands
		switch command {
		case "help":
			printHelp()
			continue
		case "quit", "exit":
			fmt.Println("Sending quit command to engine...")
			client.Quit()
			time.Sleep(100 * time.Millisecond) // Give time for engine to respond
			return
		}

		// Send command to UCI engine
		err := client.SendCommand(command)
		if err != nil {
			log.Printf("Error sending command: %v", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}

func printHelp() {
	fmt.Println("UCI Client Commands:")
	fmt.Println("  help                    - Show this help message")
	fmt.Println("  quit/exit              - Quit the client")
	fmt.Println()
	fmt.Println("Standard UCI Commands:")
	fmt.Println("  uci                    - Initialize UCI mode")
	fmt.Println("  isready                - Check if engine is ready")
	fmt.Println("  position startpos      - Set starting position")
	fmt.Println("  position fen <fen>     - Set position from FEN string")
	fmt.Println("  position startpos moves <moves> - Set position with moves")
	fmt.Println("  go depth <n>           - Search to specified depth")
	fmt.Println("  go movetime <ms>       - Search for specified time")
	fmt.Println("  go infinite            - Search indefinitely")
	fmt.Println("  stop                   - Stop current search")
	fmt.Println("  quit                   - Quit the engine")
	fmt.Println()
	fmt.Println("Usage Examples:")
	fmt.Println("  TCP mode:  go run ./cmd/engine-cli tcp localhost:8080")
	fmt.Println("  Exec mode: go run ./cmd/engine-cli exec stockfish")
	fmt.Println("  Exec mode: go run ./cmd/engine-cli exec ./my-engine --depth 15")
	fmt.Println()
	fmt.Println("Example UCI session:")
	fmt.Println("  > uci")
	fmt.Println("  < id name Stockfish 15")
	fmt.Println("  < uciok")
	fmt.Println("  > isready")
	fmt.Println("  < readyok")
	fmt.Println("  > position startpos")
	fmt.Println("  > go depth 10")
	fmt.Println("  < info depth 1 score cp 25 pv e2e4")
	fmt.Println("  < bestmove e2e4")
}