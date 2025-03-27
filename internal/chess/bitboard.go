package chess

import (
	"fmt"
	"math/bits"
)

type Bitboard uint64

// a bitboard is a 64-bit integer where each bit represents a square on the board
// a 1 in the nth bit means there is a piece of that type on that square
// https://www.chessprogramming.org/Bitboards
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

// get the least significant bit and clear it
func (b *Bitboard) PopLSB() int {
	lsb := b.BitScanForward()
	b.Clear(lsb)
	return lsb
}

// get the least significant bit
func (b *Bitboard) GetLSB() int {
	return b.BitScanForward()
}

func (b *Bitboard) BitScanForward() int {
	return bits.TrailingZeros64(uint64(*b))
}

// print the bitboard as an 8x8 grid with the lsb in the bottom left
func (b Bitboard) Print() {
	for rank := 7; rank >= 0; rank-- { // Start from rank 7 (top) down to rank 0 (bottom)
		for file := 0; file < 8; file++ { // File moves left to right
			square := rank*8 + file // Convert (rank, file) to bit index
			if (b & (1 << square)) != 0 {
				fmt.Print("1 ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println() // Newline after each rank
	}
	fmt.Println()
}
