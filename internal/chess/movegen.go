package chess

import (
	"strings"
)

/*
	All move generation is done using Bitboards.
	How exactly we implement this varies by piece type.
	But I can explain the concept easiest with the white King. All dots represent zeroes, the dots just look nicer.
	If the white king is on e4, then the bitboard for the white king would be:
		. . . . . . . .
		. . . . . . . .
		. . . . . . . .
		. . . . . . . .
		. . . . 1 . . .
		. . . . . . . .
		. . . . . . . .
		. . . . . . . .
	The bitboard of legal king moves would be all of the squares around the king
		. . . . . . . .
		. . . . . . . .
		. . . . . . . .
		. . . 1 1 1 . .
		. . . 1 . 1 . .
		. . . 1 1 1 . .
		. . . . . . . .
		. . . . . . . .
	But the king can't move to a square that has a friendly piece on it 😲.
	If there was a white pawn on D4, we have to make sure D4 is not in the legal moves bitboard.
	So we AND the white king's move bitboard with the inverse of the white bitboard.
		1 1 1 1 1 1 1 1
		1 1 1 1 1 1 1 1
		1 1 1 1 1 1 1 1
		1 1 1 . . 1 1 1
		1 1 1 1 1 1 1 1
		1 1 1 1 1 1 1 1
		1 1 1 1 1 1 1 1
	Then the result of ANDing the two bitboards is:
		. . . . . . . .
		. . . . . . . .
		. . . . . . . .
		. . . 1 1 1 . .
		. . . . . 1 . .
		. . . 1 1 1 . .
		. . . . . . . .
		. . . . . . . .
	Now, to turn this bitboard of legal moves into actual moves, we start with the original index
	of the king as the move source and we just keep popping off the least significant bit from the legal
	movoes board as a target square for a move until there are no more bits left.

	This is the general idea behind move generation for all pieces.
	Me have some idea for move masks and then we AND them with some other bitboards to get the legal moves.
*/

// gets all legal pawn moves for the current position
func (b *Board) GeneratePawnMoves() []Move {
	if b.WhiteToMove {
		return b.GenerateWhitePawnMoves()
	}
	return b.GenerateBlackPawnMoves()
}

// gets all white pawn moves
func (b *Board) GenerateWhitePawnMoves() []Move {
	wp := WHITE | PAWN
	moves := make([]Move, 0)
	enPassantSquare := NewBitboard()
	if b.EnPassantSquare != -1 {
		enPassantSquare.Set(b.EnPassantSquare)
	}
	// single push is when the square in front of the pawn is empty and not a promotion square
	singlePushBoard := (*b.Bitboards[wp] << 8) & ^*b.Bitboards[WHITE] & ^*b.Bitboards[BLACK] & ^Rank8

	// double push is another single push from legal single pushes from the pawns starting rank
	doublePushBoard := (singlePushBoard & Rank3) << 8 & ^*b.Bitboards[WHITE] & ^*b.Bitboards[BLACK] & ^Rank8

	// left capture is when there is an enemy piece on the left diagonal or the en passant square
	leftCaptureBoard := ((*b.Bitboards[wp] << 7) & ^FileA & *b.Bitboards[BLACK])

	// right capture is when there is an enemy piece on the right diagonal or the en passant square
	rightCaptureBoard := ((*b.Bitboards[wp] << 9) & ^FileH & *b.Bitboards[BLACK])

	// en passant
	leftEnPassantBoard := ((*b.Bitboards[wp] << 7) & ^FileA & *enPassantSquare)
	rightEnPassantBoard := ((*b.Bitboards[wp] << 9) & ^FileH & *enPassantSquare)

	// promotion
	promotionBoard := (*b.Bitboards[wp] << 8) & ^*b.Bitboards[WHITE] & ^*b.Bitboards[BLACK] & Rank8
	promotionLeftCaptureBoard := (*b.Bitboards[wp] << 7) & ^FileA & *b.Bitboards[BLACK] & Rank8
	promotionRightCaptureBoard := (*b.Bitboards[wp] << 9) & ^FileH & *b.Bitboards[BLACK] & Rank8

	// turn the bitboards into moves
	for singlePushBoard != 0 {
		// pop the least significant bit
		square := singlePushBoard.PopLSB()
		// create a move and append it to the moves slice
		moves = append(moves, NewMove(square-8, square, 0))
	}
	for doublePushBoard != 0 {
		square := doublePushBoard.PopLSB()
		moves = append(moves, NewMove(square-16, square, PAWN_DOUBLE_FLAG))
	}
	for leftCaptureBoard != 0 {
		square := leftCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square-7, square, 0))
	}
	for rightCaptureBoard != 0 {
		square := rightCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square-9, square, 0))
	}
	for leftEnPassantBoard != 0 {
		square := leftEnPassantBoard.PopLSB()
		moves = append(moves, NewMove(square-7, square, EN_PASSANT_FLAG))
	}
	for rightEnPassantBoard != 0 {
		square := rightEnPassantBoard.PopLSB()
		moves = append(moves, NewMove(square-9, square, EN_PASSANT_FLAG))
	}
	for promotionBoard != 0 {
		square := promotionBoard.PopLSB()
		moves = append(moves, NewMove(square-8, square, PROMOTE_QUEEN_FLAG))
		moves = append(moves, NewMove(square-8, square, PROMOTE_ROOK_FLAG))
		moves = append(moves, NewMove(square-8, square, PROMOTE_BISHOP_FLAG))
		moves = append(moves, NewMove(square-8, square, PROMOTE_KNIGHT_FLAG))
	}
	for promotionLeftCaptureBoard != 0 {
		square := promotionLeftCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square-7, square, PROMOTE_QUEEN_FLAG))
		moves = append(moves, NewMove(square-7, square, PROMOTE_ROOK_FLAG))
		moves = append(moves, NewMove(square-7, square, PROMOTE_BISHOP_FLAG))
		moves = append(moves, NewMove(square-7, square, PROMOTE_KNIGHT_FLAG))
	}
	for promotionRightCaptureBoard != 0 {
		square := promotionRightCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square-9, square, PROMOTE_QUEEN_FLAG))
		moves = append(moves, NewMove(square-9, square, PROMOTE_ROOK_FLAG))
		moves = append(moves, NewMove(square-9, square, PROMOTE_BISHOP_FLAG))
		moves = append(moves, NewMove(square-9, square, PROMOTE_KNIGHT_FLAG))
	}
	return b.FilterLegalMoves(moves)
}

// gets all black pawn moves
func (b *Board) GenerateBlackPawnMoves() []Move {
	bp := BLACK | PAWN
	moves := make([]Move, 0)
	enPassantSquare := NewBitboard()
	if b.EnPassantSquare != -1 {
		enPassantSquare.Set(b.EnPassantSquare)
	}

	// single push is when the square in front of the pawn is empty and not a promotion square
	singlePushBoard := (*b.Bitboards[bp] >> 8) & ^*b.Bitboards[WHITE] & ^*b.Bitboards[BLACK] & ^Rank1

	// double push is another single push from legal single pushes from the pawns starting rank
	doublePushBoard := (singlePushBoard & Rank6) >> 8 & ^*b.Bitboards[WHITE] & ^*b.Bitboards[BLACK] & ^Rank1

	// left capture is when there is an enemy piece on the left diagonal
	leftCaptureBoard := (*b.Bitboards[bp] >> 9) & ^FileH & *b.Bitboards[WHITE]

	// right capture is when there is an enemy piece on the right diagonal
	rightCaptureBoard := (*b.Bitboards[bp] >> 7) & ^FileA & *b.Bitboards[WHITE]

	// the en passant board is when the en passant square is the left or right diagonal of the pawn
	leftEnPassantBoard := (*b.Bitboards[bp] >> 9) & ^FileH & *enPassantSquare
	rightEnPassantBoard := (*b.Bitboards[bp] >> 7) & ^FileA & *enPassantSquare

	// promotion
	promotionBoard := (*b.Bitboards[bp] >> 8) & ^*b.Bitboards[WHITE] & ^*b.Bitboards[BLACK] & Rank1
	promotionLeftCaptureBoard := (*b.Bitboards[bp] >> 9) & ^FileH & *b.Bitboards[WHITE] & Rank1
	promotionRightCaptureBoard := (*b.Bitboards[bp] >> 7) & ^FileA & *b.Bitboards[WHITE] & Rank1

	for singlePushBoard != 0 {
		square := singlePushBoard.PopLSB()
		moves = append(moves, NewMove(square+8, square, 0))
	}
	for doublePushBoard != 0 {
		square := doublePushBoard.PopLSB()
		moves = append(moves, NewMove(square+16, square, PAWN_DOUBLE_FLAG))
	}
	for leftCaptureBoard != 0 {
		square := leftCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square+9, square, 0))
	}
	for rightCaptureBoard != 0 {
		square := rightCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square+7, square, 0))
	}
	for leftEnPassantBoard != 0 {
		square := leftEnPassantBoard.PopLSB()
		moves = append(moves, NewMove(square+9, square, EN_PASSANT_FLAG))
	}
	for rightEnPassantBoard != 0 {
		square := rightEnPassantBoard.PopLSB()
		moves = append(moves, NewMove(square+7, square, EN_PASSANT_FLAG))
	}
	for promotionBoard != 0 {
		square := promotionBoard.PopLSB()
		moves = append(moves, NewMove(square+8, square, PROMOTE_QUEEN_FLAG))
		moves = append(moves, NewMove(square+8, square, PROMOTE_ROOK_FLAG))
		moves = append(moves, NewMove(square+8, square, PROMOTE_BISHOP_FLAG))
		moves = append(moves, NewMove(square+8, square, PROMOTE_KNIGHT_FLAG))
	}
	for promotionLeftCaptureBoard != 0 {
		square := promotionLeftCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square+9, square, PROMOTE_QUEEN_FLAG))
		moves = append(moves, NewMove(square+9, square, PROMOTE_ROOK_FLAG))
		moves = append(moves, NewMove(square+9, square, PROMOTE_BISHOP_FLAG))
		moves = append(moves, NewMove(square+9, square, PROMOTE_KNIGHT_FLAG))
	}
	for promotionRightCaptureBoard != 0 {
		square := promotionRightCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square+7, square, PROMOTE_QUEEN_FLAG))
		moves = append(moves, NewMove(square+7, square, PROMOTE_ROOK_FLAG))
		moves = append(moves, NewMove(square+7, square, PROMOTE_BISHOP_FLAG))
		moves = append(moves, NewMove(square+7, square, PROMOTE_KNIGHT_FLAG))
	}

	return b.FilterLegalMoves(moves)
}

// gets all knight moves for the current position
func (b *Board) GenerateKnightMoves() []Move {
	var colorToMove byte
	if b.WhiteToMove {
		colorToMove = WHITE
	} else {
		colorToMove = BLACK
	}
	knigthtBitboard := *b.Bitboards[colorToMove|KNIGHT]
	moves := make([]Move, 0)
	for knigthtBitboard != 0 {
		fromSquare := knigthtBitboard.PopLSB()
		knightMoves := KnightMasks[fromSquare] & ^*b.Bitboards[colorToMove]
		for knightMoves != 0 {
			toSquare := knightMoves.PopLSB()
			moves = append(moves, NewMove(fromSquare, toSquare, 0))
		}
	}
	return b.FilterLegalMoves(moves)
}

// get all orthogonal moves for a piece at the given index
func (b *Board) GenerateRookMovesAtPos(pos int) []Move {
	var colorToMove byte
	if b.WhiteToMove {
		colorToMove = WHITE
	} else {
		colorToMove = BLACK
	}
	// get the blockers for the rook mask at the given position
	allPieces := *b.Bitboards[WHITE] | *b.Bitboards[BLACK]
	occ := allPieces & RookMasks[pos]
	occ *= RookMagics[pos]
	occ >>= (64 - RookShifts[pos])
	// get the legal moves for the rook at the given position
	legalMoves := RookAttacks[pos][occ] & ^*b.Bitboards[colorToMove]
	moves := make([]Move, 0)
	for legalMoves != 0 {
		toSquare := legalMoves.PopLSB()
		// add the move to the moves slice
		moves = append(moves, NewMove(pos, toSquare, 0))
	}
	return b.FilterLegalMoves(moves)
}

// gets all diagonal moves for a piece at the given index
func (b *Board) GenerateBishopMovesAtPos(pos int) []Move {
	var colorToMove byte
	if b.WhiteToMove {
		colorToMove = WHITE
	} else {
		colorToMove = BLACK
	}
	// get the blockers for the rook mask at the given position
	allPieces := *b.Bitboards[WHITE] | *b.Bitboards[BLACK]
	blockers := allPieces & BishopMasks[pos]
	// use magic to get the index
	index := (blockers * BishopMagics[pos]) >> (64 - BishopShifts[pos])
	// get the legal moves for the rook at the given position
	legalMoves := BishopAttacks[pos][index] & ^*b.Bitboards[colorToMove]
	moves := make([]Move, 0)
	for legalMoves != 0 {
		toSquare := legalMoves.PopLSB()
		// add the move to the moves slice
		moves = append(moves, NewMove(pos, toSquare, 0))
	}
	return b.FilterLegalMoves(moves)
}

// gets all bishop moves for the current position
func (b *Board) GenerateBishopMoves() []Move {
	// get all diagonal moves for each bishop
	var colorToMove byte
	if b.WhiteToMove {
		colorToMove = WHITE
	} else {
		colorToMove = BLACK
	}
	bishopBitboard := *b.Bitboards[colorToMove|BISHOP]
	if bishopBitboard == 0 {
		return []Move{}
	}
	moves := make([]Move, 0)
	for bishopBitboard != 0 {
		fromSquare := bishopBitboard.PopLSB()
		bishopMoves := b.GenerateBishopMovesAtPos(fromSquare)
		moves = append(moves, bishopMoves...)
	}
	return b.FilterLegalMoves(moves)
}

// gets all rook moves for the current position
func (b *Board) GenerateRookMoves() []Move {
	// get all sliding moves for each rook
	var colorToMove byte
	if b.WhiteToMove {
		colorToMove = WHITE
	} else {
		colorToMove = BLACK
	}
	rookBitboard := *b.Bitboards[colorToMove|ROOK]
	if rookBitboard == 0 {
		return []Move{}
	}
	moves := make([]Move, 0)
	for rookBitboard != 0 {
		fromSquare := rookBitboard.PopLSB()
		rookMoves := b.GenerateRookMovesAtPos(fromSquare)
		moves = append(moves, rookMoves...)
	}
	return b.FilterLegalMoves(moves)
}

// gets all queen moves for the current position
func (b *Board) GenerateQueenMoves() []Move {
	var queenBitboard Bitboard
	if b.WhiteToMove {
		queenBitboard = *b.Bitboards[WHITE|QUEEN]
	} else {
		queenBitboard = *b.Bitboards[BLACK|QUEEN]
	}
	if queenBitboard == 0 {
		return []Move{}
	}
	queenPos := queenBitboard.PopLSB()
	// queen moves are just the combination of bishop and rook moves
	return b.FilterLegalMoves(append(b.GenerateRookMovesAtPos(queenPos), b.GenerateBishopMovesAtPos(queenPos)...))
}

// gets all king moves for the current position
func (b *Board) GenerateKingMoves() []Move {
	var colorToMove byte
	if b.WhiteToMove {
		colorToMove = WHITE
	} else {
		colorToMove = BLACK
	}
	kingBitboard := *b.Bitboards[colorToMove|KING]
	kingPos := kingBitboard.PopLSB()
	kingMoves := KingMasks[kingPos] & ^*b.Bitboards[colorToMove]
	moves := make([]Move, 0)
	for kingMoves != 0 {
		toSquare := kingMoves.PopLSB()
		moves = append(moves, NewMove(kingPos, toSquare, 0))
	}
	// castling
	if b.WhiteToMove {
		// white king side castle
		K := strings.Contains(b.WhiteCastleRights, "K")
		if K && !b.Bitboards[WHITE].Occupied(5) && !b.Bitboards[WHITE].Occupied(6) {
			moves = append(moves, NewMove(4, 6, CASTLE_FLAG))
		}
		// white queen side castle
		Q := strings.Contains(b.WhiteCastleRights, "Q")
		if Q && !b.Bitboards[WHITE].Occupied(1) && !b.Bitboards[WHITE].Occupied(2) && !b.Bitboards[WHITE].Occupied(3) {
			moves = append(moves, NewMove(4, 2, CASTLE_FLAG))
		}
	} else {
		// black king side castle
		k := strings.Contains(b.BlackCastleRights, "k")
		if k && !b.Bitboards[BLACK].Occupied(61) && !b.Bitboards[BLACK].Occupied(62) {
			moves = append(moves, NewMove(60, 62, CASTLE_FLAG))
		}
		// black queen side castle
		q := strings.Contains(b.BlackCastleRights, "q")
		if q && !b.Bitboards[BLACK].Occupied(57) && !b.Bitboards[BLACK].Occupied(58) && !b.Bitboards[BLACK].Occupied(59) {
			moves = append(moves, NewMove(60, 58, CASTLE_FLAG))
		}
	}

	return b.FilterLegalMoves(moves)
}

// IsSquareAttacked checks if a square is attacked by the given color
func (b *Board) IsSquareAttacked(square int, byColor byte) bool {
	// Check for pawn attacks
	if byColor == WHITE {
		// White pawns attack diagonally upward (from white's perspective)
		// So from the target square, we look diagonally downward for white pawns
		if square >= 16 { // Make sure we don't go below rank 2
			// Check diagonal down-left (from white pawn's perspective)
			if square%8 != 0 && b.Bitboards[WHITE|PAWN].Occupied(square-9) {
				return true
			}
			// Check diagonal down-right (from white pawn's perspective)
			if square%8 != 7 && b.Bitboards[WHITE|PAWN].Occupied(square-7) {
				return true
			}
		}
	} else {
		// Black pawns attack diagonally downward (from black's perspective)
		// So from the target square, we look diagonally upward for black pawns
		if square < 48 { // Make sure we don't go above rank 7
			// Check diagonal up-left (from black pawn's perspective)
			if square%8 != 7 && b.Bitboards[BLACK|PAWN].Occupied(square+7) {
				return true
			}
			// Check diagonal up-right (from black pawn's perspective)
			if square%8 != 0 && b.Bitboards[BLACK|PAWN].Occupied(square+9) {
				return true
			}
		}
	}

	// Check for knight attacks
	knightMoves := KnightMasks[square]
	if (knightMoves & *b.Bitboards[byColor|KNIGHT]) != 0 {
		return true
	}

	// Check for king attacks
	kingMoves := KingMasks[square]
	if (kingMoves & *b.Bitboards[byColor|KING]) != 0 {
		return true
	}

	// Check for sliding piece attacks (rook, bishop, queen)
	allPieces := *b.Bitboards[WHITE] | *b.Bitboards[BLACK]

	// Check for rook/queen attacks (orthogonal)
	rookMask := RookMasks[square]
	blockers := allPieces & rookMask
	rookAttacks := RookAttacks[square][(blockers*RookMagics[square])>>(64-RookShifts[square])]
	if (rookAttacks & (*b.Bitboards[byColor|ROOK] | *b.Bitboards[byColor|QUEEN])) != 0 {
		return true
	}

	// Check for bishop/queen attacks (diagonal)
	bishopMask := BishopMasks[square]
	blockers = allPieces & bishopMask
	bishopAttacks := BishopAttacks[square][(blockers*BishopMagics[square])>>(64-BishopShifts[square])]
	return (bishopAttacks & (*b.Bitboards[byColor|BISHOP] | *b.Bitboards[byColor|QUEEN])) != 0
}

// IsInCheck returns true if the specified color's king is in check
func (b *Board) IsInCheck(color byte) bool {
	// Find the king position
	kingBitboard := *b.Bitboards[color|KING]
	if kingBitboard == 0 {
		return false // No king found
	}

	kingPos := kingBitboard.PopLSB()

	// Restore the king position since PopLSB modifies the bitboard
	b.Bitboards[color|KING].Set(kingPos)

	// Check if the king's square is attacked by the opposite color
	oppositeColor := WHITE
	if color == WHITE {
		oppositeColor = BLACK
	}

	return b.IsSquareAttacked(kingPos, oppositeColor)
}

func (b *Board) FilterLegalMoves(moves []Move) []Move {
	legalMoves := make([]Move, 0)

	// Determine which color is moving
	var movingColor byte
	if b.WhiteToMove {
		movingColor = WHITE
	} else {
		movingColor = BLACK
	}

	for _, move := range moves {
		// Make the move temporarily
		state := b.MakeMove(move)

		// Check if this move leaves our own king in check
		// After making the move, it's the opponent's turn, so we check if our king is still safe
		if !b.IsInCheck(movingColor) {
			legalMoves = append(legalMoves, move)
		}

		// Unmake the move to restore the board state
		b.UnmakeMove(move, state)
	}

	return legalMoves
}
