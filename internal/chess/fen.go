package chess

import (
	"strconv"
	"strings"

	"log"
)

const START_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

/*
	FEN (Forsyth-Edwards Notation) is a standard notation for describing a particular board position of a chess game.
	https://www.chessprogramming.org/Forsyth-Edwards_Notation
	https://www.chess.com/terms/fen-chess
*/

// convert FEN position to game state
func (b *Board) LoadFEN(fen string) {
	// split the FEN string into parts, 6 total
	parts := strings.Fields(fen)

	// first part is the pieces on the board
	ranks := strings.Split(parts[0], "/")
	// initialize bitboards
	b.Bitboards = make(map[Piece]*Bitboard)
	b.Bitboards[Piece(WHITE|PAWN)] = NewBitboard()
	b.Bitboards[Piece(WHITE|KNIGHT)] = NewBitboard()
	b.Bitboards[Piece(WHITE|BISHOP)] = NewBitboard()
	b.Bitboards[Piece(WHITE|ROOK)] = NewBitboard()
	b.Bitboards[Piece(WHITE|QUEEN)] = NewBitboard()
	b.Bitboards[Piece(WHITE|KING)] = NewBitboard()
	b.Bitboards[Piece(BLACK|PAWN)] = NewBitboard()
	b.Bitboards[Piece(BLACK|KNIGHT)] = NewBitboard()
	b.Bitboards[Piece(BLACK|BISHOP)] = NewBitboard()
	b.Bitboards[Piece(BLACK|ROOK)] = NewBitboard()
	b.Bitboards[Piece(BLACK|QUEEN)] = NewBitboard()
	b.Bitboards[Piece(BLACK|KING)] = NewBitboard()
	// extra bitboards for convenience
	b.Bitboards[Piece(WHITE)] = NewBitboard()
	b.Bitboards[Piece(BLACK)] = NewBitboard()
	// this one exists so that tests don't panic
	b.Bitboards[Piece(NONE)] = NewBitboard()
	// fill them out
	for rank, row := range ranks {
		file := 0
		for _, char := range row {
			square := (7-rank)*8 + file
			var piece Piece

			switch char {
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
				// Skip empty squares
				if char >= '1' && char <= '8' {
					file += int(char - '1') // -1 because we already increment file after the switch
				}
			}
			b.SetPieceAtIndex(piece, square)
			file++
		}
	}

	// second part side to move is 'w' or 'b' (white or black)
	b.WhiteToMove = parts[1] == "w"

	// third is castling rights, capital letters for white, lowercase for black, '-' for none
	// i.e. KQk for white kingside and queenside, black kingside
	castleRights := parts[2]
	if castleRights == "-" {
		b.WhiteCastleRights = ""
		b.BlackCastleRights = ""
	} else {
		for _, char := range castleRights {
			switch char {
			case 'K':
				b.WhiteCastleRights += "K"
			case 'Q':
				b.WhiteCastleRights += "Q"
			case 'k':
				b.BlackCastleRights += "k"
			case 'q':
				b.BlackCastleRights += "q"
			}
		}
	}
	// fourth is en passant target square e.g. "e3" or '-' if none
	enPassantSquare := parts[3]
	if enPassantSquare == "-" {
		b.EnPassantSquare = 0
	} else {
		b.EnPassantSquare = StringToSquare(enPassantSquare)
	}

	// fifth is halfmove clock, number of halfmoves since the last pawn move or capture
	halfMoves, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Println("invalid half move count in fen, using 0, err: ", err)
	}
	b.HalfMoves = halfMoves

	// last is fullmove clock, starts at 1 and is incremented after black moves
	fullMoves, err := strconv.Atoi(parts[5])
	if err != nil {
		log.Println("invalid full move count in fen, using 0, err: ", err)
	}
	b.FullMoves = fullMoves
}

// given a game state, return the FEN string
func (b *Board) ExportFEN() string {
	FEN := ""
	// first part is the pieces on the board
	for rank := 7; rank >= 0; rank-- {
		empty := 0
		for file := range 8 {
			square := rank*8 + file
			piece := b.GetPieceAtIndex(square)
			if piece.IsNone() {
				empty++
			} else {
				if empty > 0 {
					FEN += strconv.Itoa(empty)
					empty = 0
				}
				FEN += piece.FenChar()
			}
		}
		if empty > 0 {
			FEN += strconv.Itoa(empty)
		}
		if rank > 0 {
			FEN += "/"
		}
	}
	FEN += " "
	// second part side to move is 'w' or 'b' (white or black)
	if b.WhiteToMove {
		FEN += "w"
	} else {
		FEN += "b"
	}
	FEN += " "
	// third is castling rights, capital letters for white, lowercase for black, '-' for none
	if b.WhiteCastleRights == "" && b.BlackCastleRights == "" {
		FEN += "-"
	} else {
		FEN += b.WhiteCastleRights + b.BlackCastleRights
	}
	FEN += " "
	// fourth is en passant target square e.g. "e3" or '-' if none
	if b.EnPassantSquare == 0 {
		FEN += "-"
	} else {
		FEN += SquareToString(b.EnPassantSquare)
	}
	FEN += " "
	// fifth is halfmove clock, number of halfmoves since the last pawn move or capture
	FEN += strconv.Itoa(b.HalfMoves)
	FEN += " "
	// last is fullmove clock, starts at 1 and is incremented after black moves
	FEN += strconv.Itoa(b.FullMoves)
	return FEN
}
