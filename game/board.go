package game

import (
	"fmt"
	"strconv"
	"strings"
)

// For sanity, Board[0] is going to be A1, Board[1] is B1, etc.
type Board [64]Piece

func (b *Board) GetPiece(rank, file int) Piece {
	return b[rank*8+file]
}

func (b *Board) SetPiece(rank, file int, piece Piece) {
	b[rank*8+file] = piece
}

func (b *Board) MovePiece(oldIndex, newIndex int) {
	b[newIndex] = b[oldIndex]
	b[oldIndex] = Piece(NONE)
}

func (b *Board) ToFEN() string {
	var sb strings.Builder

	for rank := 7; rank >= 0; rank-- {
		emptyCount := 0
		for file := range 8 {
			idx := rank*8 + file
			piece := b[idx]

			if piece.IsNone() {
				emptyCount++
			} else {
				if emptyCount > 0 {
					sb.WriteString(strconv.Itoa(emptyCount))
					emptyCount = 0
				}
				sb.WriteString(piece.FenChar())
			}
		}
		if emptyCount > 0 {
			sb.WriteString(strconv.Itoa(emptyCount))
		}
		if rank > 0 {
			sb.WriteString("/")
		}
	}
	return sb.String()
}

// FromFEN initializes the board from a FEN string (supports full or partial FEN)
func NewBoardFromFEN(fen string) (*Board, error) {
	if fen == "" {
		return nil, fmt.Errorf("empty FEN string")
	}
	// Handle full FEN by taking just the board part (before first space)
	parts := strings.Split(fen, " ")
	boardFen := parts[0]

	ranks := strings.Split(boardFen, "/")
	if len(ranks) != 8 {
		return nil, fmt.Errorf("invalid FEN: expected 8 ranks, got %d", len(ranks))
	}

	b := &Board{}
	for rankIdx, rankStr := range ranks {
		rank := 7 - rankIdx // FEN starts from rank 8 (index 7)
		file := 0

		for _, c := range rankStr {
			if c >= '1' && c <= '8' {
				// skip empty squares
				file += int(c - '0')
			} else {
				var piece Piece

				switch c {
				case 'P':
					piece = Piece(WHITE | PAWN)
				case 'N':
					piece = Piece(WHITE | KNIGHT)
				case 'B':
					piece = Piece(WHITE | BISHOP)
				case 'R':
					piece = Piece(WHITE | ROOK)
				case 'Q':
					piece = Piece(WHITE | QUEEN)
				case 'K':
					piece = Piece(WHITE | KING)
				case 'p':
					piece = Piece(BLACK | PAWN)
				case 'n':
					piece = Piece(BLACK | KNIGHT)
				case 'b':
					piece = Piece(BLACK | BISHOP)
				case 'r':
					piece = Piece(BLACK | ROOK)
				case 'q':
					piece = Piece(BLACK | QUEEN)
				case 'k':
					piece = Piece(BLACK | KING)
				default:
					return nil, fmt.Errorf("invalid FEN character: %c", c)
				}

				if file >= 8 {
					return nil, fmt.Errorf("invalid FEN: too many pieces in rank %d", rank+1)
				}

				idx := rank*8 + file
				b[idx] = piece
				file++
			}
		}
		if file != 8 {
			return nil, fmt.Errorf("invalid FEN: rank %d has wrong number of squares", rank+1)
		}
	}
	return b, nil
}
