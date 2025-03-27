package chess

type Move uint16

// Move is a 16-bit integer with the following format:
// ffff tttttt ssssss
// s bits 0-5 = source square (0-63)
// t bits 6-11 = target square (0-63)
// f bits 12-15 = flag (0-7)

// flag consts, luckily 3 bits is enough for us
const (
	NO_FLAG             = 0 // 000
	PAWN_DOUBLE_FLAG    = 1 // 011
	EN_PASSANT_FLAG     = 2 // 001
	CASTLE_FLAG         = 3 // 010
	PROMOTE_KNIGHT_FLAG = 4 // 100
	PROMOTE_BISHOP_FLAG = 5 // 101
	PROMOTE_ROOK_FLAG   = 6 // 110
	PROMOTE_QUEEN_FLAG  = 7 // 111
)

func NewMove(source, target, flag int) Move {
	return Move(uint16(source) | uint16(target<<6) | uint16(flag<<12))
}

// get the string representation of a move e.g. "e2e4"
func (m *Move) String() string {
	return SquareToString(m.Source()) + SquareToString(m.Target())
}

// get the index of a move's source square
func (m *Move) Source() int {
	return int(*m & 63) // gets the last 6 bits
}

// get the index of a move's target square
func (m *Move) Target() int {
	return int((*m >> 6) & 63) // gets the next 6 bits
}

// get the flag of a move
func (m *Move) Flag() int {
	return int((*m >> 12) & 7) // gets the leading 3 bits
}

// TODO find a better file for this function
// convert square index (0-63) to a chess coordinate (e.g., "e2")
func SquareToString(square int) string {
	file := square % 8
	rank := square / 8
	// create strings using rune arrays to avoid compiler warnings
	fileStr := string([]rune{'a' + rune(file)})
	rankStr := string([]rune{'1' + rune(rank)})
	return fileStr + rankStr
}

func StringToSquare(s string) int {
	file := int(s[0] - 'a')
	rank := int(s[1] - '1')
	return rank*8 + file
}
