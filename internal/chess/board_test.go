package chess

import (
	"testing"
)

func TestGetPieceAtIndex(t *testing.T) {
	board := NewBoard()
	board.LoadFEN(START_FEN)
	piece := board.GetPieceAtIndex(0)
	if piece.IsNone() {
		t.Fatalf("Got nil piece at index 0")
	}
	if piece.Type() != ROOK {
		t.Errorf("Expected ROOK, got %d", piece.Type())
	}
	if piece.Color() != WHITE {
		t.Errorf("Expected WHITE, got %d", piece.Color())
	}
	piece = board.GetPieceAtIndex(4)
	if piece.IsNone() {
		t.Fatalf("Got nil piece at index 4")
	}
	if piece.Type() != KING {
		t.Errorf("Expected KING, got %d", piece.Type())
	}
	if piece.Color() != WHITE {
		t.Errorf("Expected WHITE, got %d", piece.Color())
	}
	piece = board.GetPieceAtIndex(60)
	if piece.IsNone() {
		t.Fatalf("Got nil piece at index 60")
	}
	if piece.Type() != KING {
		t.Errorf("Expected KING, got %d", piece.Type())
	}
	if piece.Color() != BLACK {
		t.Errorf("Expected BLACK, got %d", piece.Color())
	}
	piece = board.GetPieceAtIndex(59)
	if piece.IsNone() {
		t.Fatalf("Got nil piece at index 59")
	}
	if piece.Type() != QUEEN {
		t.Errorf("Expected QUEEN, got %d", piece.Type())
	}
	if piece.Color() != BLACK {
		t.Errorf("Expected BLACK, got %d", piece.Color())
	}
	piece = board.GetPieceAtIndex(63)
	if piece.IsNone() {
		t.Fatalf("Got nil piece at index 63")
	}
	if piece.Type() != ROOK {
		t.Errorf("Expected ROOK, got %d", piece.Type())
	}
	if piece.Color() != BLACK {
		t.Errorf("Expected BLACK, got %d", piece.Color())
	}
}

func TestSetPieceAtIndex(t *testing.T) {
	board := NewBoard()
	wr := WHITE | ROOK
	board.Bitboards = make(map[byte]*Bitboard)
	board.Bitboards[wr] = NewBitboard()
	board.Bitboards[WHITE] = NewBitboard()
	board.SetPieceAtIndex(Piece(wr), 0)
	occupied := board.Bitboards[wr].Occupied(0)
	if !occupied {
		t.Errorf("Expected square 0 to be occupied")
	}
	occupied = board.Bitboards[WHITE].Occupied(0)
	if !occupied {
		t.Errorf("Expected square 0 to be occupied")
	}
}

func TestBitboardInitialization(t *testing.T) {
	board := NewBoard()
	board.LoadFEN(START_FEN)
	// check that all bitboards are initialized
	for _, bb := range board.Bitboards {
		if bb == nil {
			t.Fatalf("Expected bitboard to be initialized")
		}
	}
}
