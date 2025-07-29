package chess

import (
	"testing"
)

func TestFENLoadExport(t *testing.T) {
	board := NewBoard()
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	board.LoadFEN(fen)
	actual := board.ExportFEN()
	if actual != fen {
		t.Errorf("Expected %s, got %s", fen, actual)
	}
}

func TestFENLoadExportVariousPositions(t *testing.T) {
	testCases := []string{
		// Starting positionF
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		// After 1.e4
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		// After 1.e4 e5
		"rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
		// Position with no castling rights
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1",
		// Position with partial castling rights
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Kq - 0 1",
		// Empty board
		"8/8/8/8/8/8/8/8 w - - 0 1",
		// Complex position
		"r1bqkb1r/pppp1ppp/2n2n2/4p3/2B1P3/3P1N2/PPP2PPP/RNBQK2R b KQkq - 0 4",
	}

	for _, fen := range testCases {
		board := NewBoard()
		board.LoadFEN(fen)
		exported := board.ExportFEN()
		if exported != fen {
			t.Errorf("FEN round-trip failed:\nInput:    %s\nExported: %s", fen, exported)
		}
	}
}

func TestFENPiecePlacement(t *testing.T) {
	board := NewBoard()

	// Test individual piece placement
	board.LoadFEN("8/8/8/8/8/8/8/R7 w - - 0 1") // White rook on a1
	piece := board.GetPieceAtIndex(0)           // a1
	if piece.Type() != ROOK || piece.Color() != WHITE {
		t.Error("Failed to load white rook on a1")
	}

	board.LoadFEN("7r/8/8/8/8/8/8/8 w - - 0 1") // Black rook on h8
	piece = board.GetPieceAtIndex(63)           // h8
	if piece.Type() != ROOK || piece.Color() != BLACK {
		t.Error("Failed to load black rook on h8")
	}

	// Test all piece types
	board.LoadFEN("rnbqkbnr/8/8/8/8/8/8/RNBQKBNR w - - 0 1")

	// Check white pieces on first rank
	expectedWhitePieces := []struct {
		square int
		piece  byte
	}{
		{0, ROOK}, {1, KNIGHT}, {2, BISHOP}, {3, QUEEN},
		{4, KING}, {5, BISHOP}, {6, KNIGHT}, {7, ROOK},
	}

	for _, expected := range expectedWhitePieces {
		piece := board.GetPieceAtIndex(expected.square)
		if piece.Type() != expected.piece || piece.Color() != WHITE {
			t.Errorf("Expected white %d on square %d, got %v", expected.piece, expected.square, piece)
		}
	}

	// Check black pieces on eighth rank
	expectedBlackPieces := []struct {
		square int
		piece  byte
	}{
		{56, ROOK}, {57, KNIGHT}, {58, BISHOP}, {59, QUEEN},
		{60, KING}, {61, BISHOP}, {62, KNIGHT}, {63, ROOK},
	}

	for _, expected := range expectedBlackPieces {
		piece := board.GetPieceAtIndex(expected.square)
		if piece.Type() != expected.piece || piece.Color() != BLACK {
			t.Errorf("Expected black %d on square %d, got %v", expected.piece, expected.square, piece)
		}
	}
}

func TestFENActiveColor(t *testing.T) {
	board := NewBoard()

	// Test white to move
	board.LoadFEN("8/8/8/8/8/8/8/8 w - - 0 1")
	if !board.WhiteToMove {
		t.Error("Expected white to move")
	}

	// Test black to move
	board.LoadFEN("8/8/8/8/8/8/8/8 b - - 0 1")
	if board.WhiteToMove {
		t.Error("Expected black to move")
	}
}

func TestFENCastlingRights(t *testing.T) {
	board := NewBoard()

	// Test all castling rights
	board.LoadFEN("8/8/8/8/8/8/8/8 w KQkq - 0 1")
	if board.WhiteCastleRights != "KQ" {
		t.Errorf("Expected white castle rights 'KQ', got '%s'", board.WhiteCastleRights)
	}
	if board.BlackCastleRights != "kq" {
		t.Errorf("Expected black castle rights 'kq', got '%s'", board.BlackCastleRights)
	}

	// Test partial rights
	board.LoadFEN("8/8/8/8/8/8/8/8 w Kq - 0 1")
	if board.WhiteCastleRights != "K" {
		t.Errorf("Expected white castle rights 'K', got '%s'", board.WhiteCastleRights)
	}
	if board.BlackCastleRights != "q" {
		t.Errorf("Expected black castle rights 'q', got '%s'", board.BlackCastleRights)
	}

	// Test no castling rights
	board.LoadFEN("8/8/8/8/8/8/8/8 w - - 0 1")
	if board.WhiteCastleRights != "" {
		t.Errorf("Expected no white castle rights, got '%s'", board.WhiteCastleRights)
	}
	if board.BlackCastleRights != "" {
		t.Errorf("Expected no black castle rights, got '%s'", board.BlackCastleRights)
	}
}

func TestFENEnPassant(t *testing.T) {
	board := NewBoard()

	// Test en passant square with a real position
	board.LoadFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
	if board.EnPassantSquare != StringToSquare("e3") {
		t.Errorf("Expected en passant square e3 (%d), got %d", StringToSquare("e3"), board.EnPassantSquare)
	}

	// Test no en passant square
	board.LoadFEN(START_FEN)
	if board.EnPassantSquare != -1 {
		t.Errorf("Expected no en passant square (-1), got %d", board.EnPassantSquare)
	}
}

func TestFENMoveCounts(t *testing.T) {
	board := NewBoard()

	// Test move counts
	board.LoadFEN("8/8/8/8/8/8/8/8 w - - 5 10")
	if board.HalfMoves != 5 {
		t.Errorf("Expected half moves 5, got %d", board.HalfMoves)
	}
	if board.FullMoves != 10 {
		t.Errorf("Expected full moves 10, got %d", board.FullMoves)
	}
}

func TestFENEdgeCases(t *testing.T) {
	board := NewBoard()

	// Test position with many consecutive empty squares
	board.LoadFEN("8/8/8/8/8/8/8/8 w - - 0 1")
	for i := 0; i < 64; i++ {
		if !board.GetPieceAtIndex(i).IsNone() {
			t.Errorf("Expected empty square at %d", i)
		}
	}

	// Test position with mixed empty squares and pieces
	board.LoadFEN("r6r/8/8/8/8/8/8/R6R w - - 0 1")

	// Should have rooks on corners
	corners := []int{0, 7, 56, 63} // a1, h1, a8, h8
	for _, corner := range corners {
		piece := board.GetPieceAtIndex(corner)
		if piece.Type() != ROOK {
			t.Errorf("Expected rook on corner %d", corner)
		}
	}

	// All other squares should be empty
	for i := 0; i < 64; i++ {
		isCorner := false
		for _, corner := range corners {
			if i == corner {
				isCorner = true
				break
			}
		}
		if !isCorner && !board.GetPieceAtIndex(i).IsNone() {
			t.Errorf("Expected empty square at %d", i)
		}
	}
}
