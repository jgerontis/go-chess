package engine

import (
	"fmt"

	"github.com/jgerontis/go-chess/internal/chess"
	"github.com/jgerontis/go-chess/internal/uci"
)

/*
	See: https://www.chessprogramming.org/UCI
	GoChess Engine - A chess engine implementation using the UCI protocol.
	The engine implements the uci.EngineHandler interface to work with
	the shared UCI server infrastructure.
*/

// GoChessEngine represents our chess engine implementation
type GoChessEngine struct {
	board      *chess.Board
	stopChan   chan struct{}
	searching  bool
	currentBest chess.Move
}

// NewGoChessEngine creates a new instance of our chess engine
func NewGoChessEngine() *GoChessEngine {
	return &GoChessEngine{
		board:    chess.NewBoard(),
		stopChan: make(chan struct{}),
	}
}

// GetInfo returns basic engine information
func (e *GoChessEngine) GetInfo() uci.EngineInfo {
	return uci.EngineInfo{
		Name:   "GoChess",
		Author: "Josh Gerontis",
	}
}

// SetPosition sets the current board position
func (e *GoChessEngine) SetPosition(fen string, moves []string) error {
	// Load the position
	if fen == "" {
		fen = chess.START_FEN
	}

	e.board.LoadFEN(fen)

	// Apply moves if provided
	for _, moveStr := range moves {
		move, err := e.parseMove(moveStr)
		if err != nil {
			return fmt.Errorf("invalid move %s: %v", moveStr, err)
		}

		// Generate legal moves to validate the move
		e.board.GenerateLegalMoves()
		isLegal := false
		for _, legalMove := range e.board.LegalMoves {
			if legalMove == move {
				isLegal = true
				break
			}
		}

		if !isLegal {
			return fmt.Errorf("illegal move: %s", moveStr)
		}

		// Make the move
		_ = e.board.MakeMove(move)
	}

	return nil
}

// Search performs a search and returns the best move
func (e *GoChessEngine) Search(params uci.SearchParams) (string, error) {
	e.searching = true
	e.currentBest = 0
	
	// Generate legal moves first
	e.board.GenerateLegalMoves()
	if len(e.board.LegalMoves) == 0 {
		e.searching = false
		return "(none)", fmt.Errorf("no legal moves available")
	}
	
	// Start with the first legal move as a fallback
	e.currentBest = e.board.LegalMoves[0]
	
	// Create a new stop channel for this search
	e.stopChan = make(chan struct{})
	
	// Run search in goroutine and wait for result or stop signal
	resultChan := make(chan chess.Move, 1)
	
	go func() {
		// Your actual search logic would go here
		// For now, simulate with the simple search
		bestMove := FindBestMove(e.board)
		
		select {
		case resultChan <- bestMove:
			// Search completed normally
		case <-e.stopChan:
			// Search was stopped, send current best
			resultChan <- e.currentBest
		}
	}()
	
	// Wait for either search completion or stop signal
	var bestMove chess.Move
	select {
	case bestMove = <-resultChan:
		// Search completed
	case <-e.stopChan:
		// Search was stopped
		bestMove = e.currentBest
	}
	
	e.searching = false
	
	if bestMove == 0 {
		bestMove = e.currentBest
	}
	
	return bestMove.String(), nil
}

// IsReady returns true if the engine is ready to receive commands
func (e *GoChessEngine) IsReady() bool {
	return true
}

// Stop stops the current search
func (e *GoChessEngine) Stop() {
	if e.searching {
		close(e.stopChan)
	}
}

// parseMove converts a UCI move string to a chess.Move
func (e *GoChessEngine) parseMove(moveStr string) (chess.Move, error) {
	if len(moveStr) < 4 {
		return 0, fmt.Errorf("move string too short: %s", moveStr)
	}

	fromSquare := chess.StringToSquare(moveStr[:2])
	toSquare := chess.StringToSquare(moveStr[2:4])

	if fromSquare < 0 || fromSquare >= 64 || toSquare < 0 || toSquare >= 64 {
		return 0, fmt.Errorf("invalid squares in move: %s", moveStr)
	}

	// Handle promotion
	flags := chess.NO_FLAG
	if len(moveStr) == 5 {
		switch moveStr[4] {
		case 'q':
			flags = chess.PROMOTE_QUEEN_FLAG
		case 'r':
			flags = chess.PROMOTE_ROOK_FLAG
		case 'b':
			flags = chess.PROMOTE_BISHOP_FLAG
		case 'n':
			flags = chess.PROMOTE_KNIGHT_FLAG
		}
	}

	return chess.NewMove(fromSquare, toSquare, flags), nil
}

// StartEngine runs the UCI engine loop using the new infrastructure
func StartEngine() {
	engine := NewGoChessEngine()
	server := uci.NewServer(engine)
	server.Run()
}
