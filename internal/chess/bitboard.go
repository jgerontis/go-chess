package chess

import "math/bits"

type Bitboard uint64

// a bitboard is a 64-bit integer where each bit represents a square on the board
// a 1 in the nth bit means there is a piece of that type on that square
const (
	Rank8 Bitboard = 0xFF00000000000000 // 0xFF is the leading 8 bits set to 1
	Rank7 Bitboard = 0x00FF000000000000
	Rank6 Bitboard = 0x0000FF0000000000
	Rank5 Bitboard = 0x000000FF00000000
	Rank4 Bitboard = 0x00000000FF000000
	Rank3 Bitboard = 0x0000000000FF0000
	Rank2 Bitboard = 0x000000000000FF00
	Rank1 Bitboard = 0x00000000000000FF
	FileH Bitboard = 0x8080808080808080 // 0x80 is the 8th bit set to 1
	FileG Bitboard = 0x4040404040404040
	FileF Bitboard = 0x2020202020202020
	FileE Bitboard = 0x1010101010101010
	FileD Bitboard = 0x0808080808080808
	FileC Bitboard = 0x0404040404040404
	FileB Bitboard = 0x0202020202020202
	FileA Bitboard = 0x0101010101010101
)

func NewBitboard() *Bitboard {
	return new(Bitboard)
}

// set a piece at a given square
func (b *Bitboard) Set(square int) {
	*b |= 1 << square
}

// remove a piece at a given square
func (b *Bitboard) Clear(square int) {
	*b &= ^(1 << square)
}

// check if a square is occupied
func (b *Bitboard) Occupied(square int) bool {
	return (*b & (1 << square)) != 0
}

// get the least significant bit
func (b *Bitboard) PopLSB() int {
	lsb := b.BitScanForward()
	b.Clear(lsb)
	return lsb
}

func (b *Bitboard) BitScanForward() int {
	return bits.TrailingZeros64(uint64(*b))
}

// print the bitboard (usually for debug)
func (b Bitboard) Print() {
	for i := range 64 {
		if i%8 == 0 {
			println()
		}
		if b&(1<<i) != 0 {
			print("1")
		} else {
			print("0")
		}
	}
	println()
}
