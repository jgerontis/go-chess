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
