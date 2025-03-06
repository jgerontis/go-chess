package chess

import (
	"testing"
)

func TestPieceLogic(t *testing.T) {
	p := Piece(PAWN | WHITE)
	if p.Color() != WHITE {
		t.Errorf("Expected WHITE, got %d", p.Color())
	}
	if p.Type() != PAWN {
		t.Errorf("Expected PAWN, got %d", p.Type())
	}
	if p.IsNone() {
		t.Errorf("Expected false, got true")
	}
	if p.String() != "♟" {
		t.Errorf("Expected ♟, got %s", p.String())
	}
}
