package engine

type Bitboard uint64

// a bitboard is a 64-bit integer where each bit represents a square on the board
// a 1 in the nth bit means there is a piece of that type on that square

// set a piece at a given square
func (b *Bitboard) Set(square int) {
	*b |= 1 << square
}

// remove a piece at a given square
func (b *Bitboard) Clear(square int) {
	*b &= ^(1 << square)
}

// check if a square is occupied
func (b Bitboard) Occupied(square int) bool {
	return (b & (1 << square)) != 0
}

// print the bitboard (usually for debug)
func (b *Bitboard) Print() {
	for i := range 64 {
		if i%8 == 0 {
			println()
		}
		if *b&(1<<i) != 0 {
			print("1")
		} else {
			print("0")
		}
	}
	println()
}
