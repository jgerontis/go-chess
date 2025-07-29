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
	if bb<<1 != 0 {
		t.Errorf("Expected shift left to be 0")
	}
	if bb>>1 != 0x4000000000000000 {
		t.Errorf("Expected shift right to be 0x4000000000000000")
	}
	newBB := NewBitboard()
	newBB.Set(63)
	if *newBB<<1 != 0 {
		t.Errorf("Expected shift left to be 0")
	}
	if *newBB>>1 != 0x4000000000000000 {
		t.Errorf("Expected shift right to be 0x4000000000000000")
	}
}

func TestBitboardNot(t *testing.T) {
	bb := Bitboard(0x8000000000000000)
	if ^bb != 0x7FFFFFFFFFFFFFFF {
		t.Errorf("Expected not to be 0x7FFFFFFFFFFFFFFF")
	}
	newBB := NewBitboard()
	newBB.Set(0)
	if ^*newBB != 0xFFFFFFFFFFFFFFFE {
		t.Errorf("Expected not to be 0xFFFFFFFFFFFFFFFE")
	}
}

func TestBitboardClear(t *testing.T) {
	bb := NewBitboard()
	
	// Set a bit and then clear it
	bb.Set(15)
	if !bb.Occupied(15) {
		t.Error("Bit 15 should be set")
	}
	
	bb.Clear(15)
	if bb.Occupied(15) {
		t.Error("Bit 15 should be cleared")
	}
	
	// Test clearing a bit that's already clear
	bb.Clear(20)
	if bb.Occupied(20) {
		t.Error("Bit 20 should remain clear")
	}
}

func TestBitboardPopLSB(t *testing.T) {
	bb := NewBitboard()
	
	// Set multiple bits
	bb.Set(0)
	bb.Set(5)
	bb.Set(63)
	
	// PopLSB should return the lowest set bit and clear it
	lsb := bb.PopLSB()
	if lsb != 0 {
		t.Errorf("Expected LSB to be 0, got %d", lsb)
	}
	if bb.Occupied(0) {
		t.Error("Bit 0 should be cleared after PopLSB")
	}
	
	// Next PopLSB should return 5
	lsb = bb.PopLSB()
	if lsb != 5 {
		t.Errorf("Expected LSB to be 5, got %d", lsb)
	}
	if bb.Occupied(5) {
		t.Error("Bit 5 should be cleared after PopLSB")
	}
	
	// Last PopLSB should return 63
	lsb = bb.PopLSB()
	if lsb != 63 {
		t.Errorf("Expected LSB to be 63, got %d", lsb)
	}
	if bb.Occupied(63) {
		t.Error("Bit 63 should be cleared after PopLSB")
	}
	
	// Bitboard should now be empty
	if *bb != 0 {
		t.Error("Bitboard should be empty after all PopLSB operations")
	}
}

func TestBitboardGetLSB(t *testing.T) {
	bb := NewBitboard()
	
	// Set multiple bits
	bb.Set(3)
	bb.Set(10)
	bb.Set(45)
	
	// GetLSB should return the lowest set bit without clearing it
	lsb := bb.GetLSB()
	if lsb != 3 {
		t.Errorf("Expected LSB to be 3, got %d", lsb)
	}
	if !bb.Occupied(3) {
		t.Error("Bit 3 should still be set after GetLSB")
	}
	
	// Multiple calls should return the same value
	lsb2 := bb.GetLSB()
	if lsb2 != 3 {
		t.Errorf("Expected LSB to still be 3, got %d", lsb2)
	}
}

func TestBitboardAnd(t *testing.T) {
	bb1 := NewBitboard()
	bb2 := NewBitboard()
	
	bb1.Set(0)
	bb1.Set(5)
	bb1.Set(10)
	
	bb2.Set(5)
	bb2.Set(10)
	bb2.Set(15)
	
	// AND operation
	result := *bb1 & *bb2
	
	// Result should have bits 5 and 10 set
	if (result & (1 << 5)) == 0 {
		t.Error("Bit 5 should be set in AND result")
	}
	if (result & (1 << 10)) == 0 {
		t.Error("Bit 10 should be set in AND result")
	}
	if (result & (1 << 0)) != 0 {
		t.Error("Bit 0 should not be set in AND result")
	}
	if (result & (1 << 15)) != 0 {
		t.Error("Bit 15 should not be set in AND result")
	}
}

func TestBitboardOr(t *testing.T) {
	bb1 := NewBitboard()
	bb2 := NewBitboard()
	
	bb1.Set(0)
	bb1.Set(5)
	
	bb2.Set(5)
	bb2.Set(10)
	
	// OR operation
	result := *bb1 | *bb2
	
	// Result should have bits 0, 5, and 10 set
	if (result & (1 << 0)) == 0 {
		t.Error("Bit 0 should be set in OR result")
	}
	if (result & (1 << 5)) == 0 {
		t.Error("Bit 5 should be set in OR result")
	}
	if (result & (1 << 10)) == 0 {
		t.Error("Bit 10 should be set in OR result")
	}
}

func TestBitboardConstants(t *testing.T) {
	// Test rank constants
	if Rank1 != 0x00000000000000FF {
		t.Errorf("Rank1 constant incorrect: got 0x%016X", Rank1)
	}
	if Rank8 != 0xFF00000000000000 {
		t.Errorf("Rank8 constant incorrect: got 0x%016X", Rank8)
	}
	
	// Test file constants
	if FileA != 0x0101010101010101 {
		t.Errorf("FileA constant incorrect: got 0x%016X", FileA)
	}
	if FileH != 0x8080808080808080 {
		t.Errorf("FileH constant incorrect: got 0x%016X", FileH)
	}
	
	// Test that rank constants cover the correct squares
	for i := 0; i < 8; i++ {
		if (Rank1 & (1 << i)) == 0 {
			t.Errorf("Rank1 should include square %d", i)
		}
		if (Rank8 & (1 << (56+i))) == 0 {
			t.Errorf("Rank8 should include square %d", 56+i)
		}
	}
	
	// Test that file constants cover the correct squares
	for i := 0; i < 8; i++ {
		if (FileA & (1 << (i*8))) == 0 {
			t.Errorf("FileA should include square %d", i*8)
		}
		if (FileH & (1 << (i*8+7))) == 0 {
			t.Errorf("FileH should include square %d", i*8+7)
		}
	}
}
