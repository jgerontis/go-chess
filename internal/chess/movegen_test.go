package chess

import (
	"testing"
)

func TestGeneratePawnMoves(t *testing.T) {
	b := NewBoard()
	b.LoadFEN(START_FEN)
	b.GeneratePawnMoves()
}
