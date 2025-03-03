package chess

import "strings"

type Piece byte

// each piece is a byte with the format 000WBTTT
const (
	NONE   byte = 0  // 00000000
	PAWN   byte = 1  // 00000001
	KNIGHT byte = 2  // 00000010
	BISHOP byte = 3  // 00000011
	ROOK   byte = 4  // 00000100
	QUEEN  byte = 5  // 00000101
	KING   byte = 6  // 00000110
	BLACK  byte = 8  // 00001000
	WHITE  byte = 16 // 00010000
)

func (p Piece) Color() byte {
	return byte(p) & 24
}

func (p Piece) Type() byte {
	return byte(p) & 7
}

// IsEmpty returns true if the piece is empty (no piece)
func (p Piece) IsNone() bool {
	return p == Piece(NONE) || (p&0x07) == 0
}

func (p Piece) String() string {
	switch p {
	case Piece(PAWN | BLACK):
		return "♙"
	case Piece(KNIGHT | BLACK):
		return "♘"
	case Piece(BISHOP | BLACK):
		return "♗"
	case Piece(ROOK | BLACK):
		return "♖"
	case Piece(QUEEN | BLACK):
		return "♕"
	case Piece(KING | BLACK):
		return "♔"
	case Piece(PAWN | WHITE):
		return "♟"
	case Piece(KNIGHT | WHITE):
		return "♞"
	case Piece(BISHOP | WHITE):
		return "♝"
	case Piece(ROOK | WHITE):
		return "♜"
	case Piece(QUEEN | WHITE):
		return "♛"
	case Piece(KING | WHITE):
		return "♚"
	default:
		return " "
	}
}

func (p Piece) FenChar() string {
	if p.IsNone() {
		return ""
	}

	// FEN piece characters (lowercase for black, uppercase for white)
	pieceChars := []string{"?", "p", "n", "b", "r", "q", "k"}
	char := pieceChars[p.Type()]

	if p.Color() == WHITE {
		return strings.ToUpper(char)
	}
	return char
}
