package uci

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// EngineInfo represents basic engine information
type EngineInfo struct {
	Name   string
	Author string
}

// EngineHandler defines the interface for engine implementations
type EngineHandler interface {
	// GetInfo returns basic engine information
	GetInfo() EngineInfo
	
	// SetPosition sets the current board position
	SetPosition(fen string, moves []string) error
	
	// Search performs a search and returns the best move
	Search(searchParams SearchParams) (string, error)
	
	// IsReady returns true if the engine is ready to receive commands
	IsReady() bool
	
	// Stop stops the current search
	Stop()
}

// SearchParams represents search parameters
type SearchParams struct {
	Depth    int
	Movetime int // milliseconds
	Infinite bool
}

// Server represents a UCI server that handles engine communication
type Server struct {
	engine EngineHandler
	reader io.Reader
	writer io.Writer
	quit   bool
}

// NewServer creates a new UCI server with the given engine handler
func NewServer(engine EngineHandler) *Server {
	return &Server{
		engine: engine,
		reader: os.Stdin,
		writer: os.Stdout,
	}
}

// NewServerWithIO creates a new UCI server with custom I/O
func NewServerWithIO(engine EngineHandler, reader io.Reader, writer io.Writer) *Server {
	return &Server{
		engine: engine,
		reader: reader,
		writer: writer,
	}
}

// Run starts the UCI server loop
func (s *Server) Run() error {
	scanner := bufio.NewScanner(s.reader)
	
	for scanner.Scan() && !s.quit {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		err := s.handleCommand(line)
		if err != nil {
			return err
		}
	}
	
	return scanner.Err()
}

// handleCommand processes a single UCI command
func (s *Server) handleCommand(line string) error {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil
	}
	
	command := parts[0]
	args := parts[1:]
	
	switch command {
	case "uci":
		return s.handleUCI()
	case "isready":
		return s.handleIsReady()
	case "position":
		return s.handlePosition(args)
	case "go":
		return s.handleGo(args)
	case "stop":
		return s.handleStop()
	case "quit":
		s.quit = true
		return nil
	default:
		// Unknown command, ignore per UCI specification
		return nil
	}
}

// handleUCI responds to the "uci" command
func (s *Server) handleUCI() error {
	info := s.engine.GetInfo()
	fmt.Fprintf(s.writer, "id name %s\n", info.Name)
	fmt.Fprintf(s.writer, "id author %s\n", info.Author)
	fmt.Fprintf(s.writer, "uciok\n")
	return nil
}

// handleIsReady responds to the "isready" command
func (s *Server) handleIsReady() error {
	if s.engine.IsReady() {
		fmt.Fprintf(s.writer, "readyok\n")
	}
	return nil
}

// handlePosition processes the "position" command
func (s *Server) handlePosition(args []string) error {
	if len(args) == 0 {
		return nil
	}
	
	var fen string
	var moves []string
	
	if args[0] == "startpos" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		args = args[1:]
	} else if args[0] == "fen" && len(args) > 1 {
		// Reconstruct FEN from args[1] onwards until "moves" or end
		var fenParts []string
		i := 1
		for i < len(args) && args[i] != "moves" {
			fenParts = append(fenParts, args[i])
			i++
		}
		fen = strings.Join(fenParts, " ")
		args = args[i:]
	}
	
	// Check for moves
	if len(args) > 0 && args[0] == "moves" {
		moves = args[1:]
	}
	
	return s.engine.SetPosition(fen, moves)
}

// handleGo processes the "go" command
func (s *Server) handleGo(args []string) error {
	params := SearchParams{}
	
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "depth":
			if i+1 < len(args) {
				if depth, err := strconv.Atoi(args[i+1]); err == nil {
					params.Depth = depth
				}
				i++
			}
		case "movetime":
			if i+1 < len(args) {
				if movetime, err := strconv.Atoi(args[i+1]); err == nil {
					params.Movetime = movetime
				}
				i++
			}
		case "infinite":
			params.Infinite = true
		}
	}
	
	// Start search in a goroutine for proper async behavior
	go func() {
		bestMove, err := s.engine.Search(params)
		if err != nil {
			// In a real implementation, you might want to log this error
			// For now, output a fallback move
			fmt.Fprintf(s.writer, "bestmove (none)\n")
			return
		}
		fmt.Fprintf(s.writer, "bestmove %s\n", bestMove)
	}()
	
	return nil
}

// handleStop processes the "stop" command
func (s *Server) handleStop() error {
	s.engine.Stop()
	return nil
}