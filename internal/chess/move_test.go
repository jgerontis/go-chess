package chess

import (
	"testing"
)

func TestMoveLogic(t *testing.T) {
	move := NewMove(StringToSquare("e2"), StringToSquare("e4"), 0)
	if move.Source() != StringToSquare("e2") {
		t.Errorf("Expected e2, got %d", move.Source())
	}
	if move.Target() != StringToSquare("e4") {
		t.Errorf("Expected e4, got %d", move.Target())
	}
	if move.Flag() != 0 {
		t.Errorf("Expected 0, got %d", move.Flag())
	}
	if move.String() != "e2e4" {
		t.Errorf("Expected e2e4, got %s", move.String())
	}
}

func TestStringToSquare(t *testing.T) {
	square := StringToSquare("e2")
	if square != 12 {
		t.Errorf("Expected 12, got %d", square)
	}
	square = StringToSquare("a1")
	if square != 0 {
		t.Errorf("Expected 0, got %d", square)
	}
	square = StringToSquare("h8")
	if square != 63 {
		t.Errorf("Expected 63, got %d", square)
	}
}

func TestSquareToString(t *testing.T) {
	str := SquareToString(12)
	if str != "e2" {
		t.Errorf("Expected e2, got %s", str)
	}
	str = SquareToString(0)
	if str != "a1" {
		t.Errorf("Expected a1, got %s", str)
	}
	str = SquareToString(63)
	if str != "h8" {
		t.Errorf("Expected h8, got %s", str)
	}
}
