package chess

import (
	"testing"
)

func TestBitboardSet(t *testing.T) {
	bb := Bitboard(0)
	bb.Set(0)
	if !bb.Occupied(0) {
		t.Errorf("Expected square 0 to be occupied")
	}
	bb.Set(63)
	if !bb.Occupied(63) {
		t.Errorf("Expected square 63 to be occupied")
	}
	newBB := NewBitboard()
	newBB.Set(0)
	if !newBB.Occupied(0) {
		t.Errorf("Expected square 0 to be occupied")
	}
	newBB.Set(63)
	if !newBB.Occupied(63) {
		t.Errorf("Expected square 63 to be occupied")
	}
}

func TestBitboardShift(t *testing.T) {
	bb := Bitboard(0x8000000000000000)
	if bb.SLeft(1) != 0 {
		t.Errorf("Expected shift left to be 0")
	}
	if bb.SRight(1) != 0x4000000000000000 {
		t.Errorf("Expected shift right to be 0x4000000000000000")
	}
	newBB := NewBitboard()
	newBB.Set(63)
	if newBB.SLeft(1) != 0 {
		t.Errorf("Expected shift left to be 0")
	}
	if newBB.SRight(1) != 0x4000000000000000 {
		t.Errorf("Expected shift right to be 0x4000000000000000")
	}
}

func TestBitboardNot(t *testing.T) {
	bb := Bitboard(0x8000000000000000)
	if bb.Not() != 0x7FFFFFFFFFFFFFFF {
		t.Errorf("Expected not to be 0x7FFFFFFFFFFFFFFF")
	}
	newBB := NewBitboard()
	newBB.Set(0)
	if newBB.Not() != 0xFFFFFFFFFFFFFFFE {
		t.Errorf("Expected not to be 0xFFFFFFFFFFFFFFFE")
	}
}
