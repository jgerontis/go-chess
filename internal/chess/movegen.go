package chess

// gets all legal pawn moves for the current position
func (b *Board) GeneratePawnMoves() []Move {
	if b.WhiteToMove {
		return b.GenerateWhitePawnMoves()
	}
	return b.GenerateBlackPawnMoves()
}

// gets all white pawn moves
func (b *Board) GenerateWhitePawnMoves() []Move {
	wp := Piece(WHITE | PAWN)
	black := Piece(BLACK)
	white := Piece(WHITE)
	moves := make([]Move, 0)

	// single push is when the square in front of the pawn is empty and not a promotion square
	singlePushBoard := (*b.Bitboards[wp] << 8) & ^*b.Bitboards[white] & ^*b.Bitboards[black] & ^Rank8

	// double push is another single push from legal single pushes from the pawns starting rank
	doublePushBoard := (singlePushBoard & Rank3) << 8 & ^*b.Bitboards[white] & ^*b.Bitboards[black] & ^Rank8

	// left capture is when there is an enemy piece on the left diagonal
	leftCaptureBoard := (*b.Bitboards[wp] << 7) & ^FileA & *b.Bitboards[black]

	// right capture is when there is an enemy piece on the right diagonal
	rightCaptureBoard := (*b.Bitboards[wp] << 9) & ^FileH & *b.Bitboards[black]

	// en passant is when there is pawn on the 5th rank and the last move was a double push

	// promotion
	promotionBoard := (*b.Bitboards[wp] << 8) & ^*b.Bitboards[white] & ^*b.Bitboards[black] & Rank8
	promotionLeftCaptureBoard := (*b.Bitboards[wp] << 7) & ^FileA & *b.Bitboards[black] & Rank8
	promotionRightCaptureBoard := (*b.Bitboards[wp] << 9) & ^FileH & *b.Bitboards[black] & Rank8

	// turn the bitboards into moves
	for singlePushBoard != 0 {
		// pop the least significant bit
		square := singlePushBoard.PopLSB()
		// create a move and append it to the moves slice
		moves = append(moves, NewMove(square-8, square, 0))
	}
	for doublePushBoard != 0 {
		square := doublePushBoard.PopLSB()
		moves = append(moves, NewMove(square-16, square, 0))
	}
	for leftCaptureBoard != 0 {
		square := leftCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square-7, square, 0))
	}
	for rightCaptureBoard != 0 {
		square := rightCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square-9, square, 0))
	}
	for promotionBoard != 0 {
		square := promotionBoard.PopLSB()
		moves = append(moves, NewMove(square-8, square, PromoteToQueenFlag))
		moves = append(moves, NewMove(square-8, square, PromoteToRookFlag))
		moves = append(moves, NewMove(square-8, square, PromoteToBishopFlag))
		moves = append(moves, NewMove(square-8, square, PromoteToKnightFlag))
	}
	for promotionLeftCaptureBoard != 0 {
		square := promotionLeftCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square-7, square, PromoteToQueenFlag))
		moves = append(moves, NewMove(square-7, square, PromoteToRookFlag))
		moves = append(moves, NewMove(square-7, square, PromoteToBishopFlag))
		moves = append(moves, NewMove(square-7, square, PromoteToKnightFlag))
	}
	for promotionRightCaptureBoard != 0 {
		square := promotionRightCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square-9, square, PromoteToQueenFlag))
		moves = append(moves, NewMove(square-9, square, PromoteToRookFlag))
		moves = append(moves, NewMove(square-9, square, PromoteToBishopFlag))
		moves = append(moves, NewMove(square-9, square, PromoteToKnightFlag))
	}
	return moves
}

// gets all black pawn moves
func (b *Board) GenerateBlackPawnMoves() []Move {
	bp := Piece(BLACK | PAWN)
	black := Piece(BLACK)
	white := Piece(WHITE)
	moves := make([]Move, 0)

	// single push is when the square in front of the pawn is empty and not a promotion square
	singlePushBoard := (*b.Bitboards[bp] >> 8) & ^*b.Bitboards[white] & ^*b.Bitboards[black] & ^Rank1

	// double push is another single push from legal single pushes from the pawns starting rank
	doublePushBoard := (singlePushBoard & Rank6) >> 8 & ^*b.Bitboards[white] & ^*b.Bitboards[black] & ^Rank1

	// left capture is when there is an enemy piece on the left diagonal
	leftCaptureBoard := (*b.Bitboards[bp] >> 9) & ^FileH & *b.Bitboards[white]

	// right capture is when there is an enemy piece on the right diagonal
	rightCaptureBoard := (*b.Bitboards[bp] >> 7) & ^FileA & *b.Bitboards[white]

	// en passant is when there is pawn on the 5th rank and the last move was a double push

	// promotion
	promotionBoard := (*b.Bitboards[bp] >> 8) & ^*b.Bitboards[white] & ^*b.Bitboards[black] & Rank1
	promotionLeftCaptureBoard := (*b.Bitboards[bp] >> 9) & ^FileH & *b.Bitboards[white] & Rank1
	promotionRightCaptureBoard := (*b.Bitboards[bp] >> 7) & ^FileA & *b.Bitboards[white] & Rank1

	for singlePushBoard != 0 {
		square := singlePushBoard.PopLSB()
		moves = append(moves, NewMove(square+8, square, 0))
	}
	for doublePushBoard != 0 {
		square := doublePushBoard.PopLSB()
		moves = append(moves, NewMove(square+16, square, 0))
	}
	for leftCaptureBoard != 0 {
		square := leftCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square+9, square, 0))
	}
	for rightCaptureBoard != 0 {
		square := rightCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square+7, square, 0))
	}
	for promotionBoard != 0 {
		square := promotionBoard.PopLSB()
		moves = append(moves, NewMove(square+8, square, PromoteToQueenFlag))
		moves = append(moves, NewMove(square+8, square, PromoteToRookFlag))
		moves = append(moves, NewMove(square+8, square, PromoteToBishopFlag))
		moves = append(moves, NewMove(square+8, square, PromoteToKnightFlag))
	}
	for promotionLeftCaptureBoard != 0 {
		square := promotionLeftCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square+9, square, PromoteToQueenFlag))
		moves = append(moves, NewMove(square+9, square, PromoteToRookFlag))
		moves = append(moves, NewMove(square+9, square, PromoteToBishopFlag))
		moves = append(moves, NewMove(square+9, square, PromoteToKnightFlag))
	}
	for promotionRightCaptureBoard != 0 {
		square := promotionRightCaptureBoard.PopLSB()
		moves = append(moves, NewMove(square+7, square, PromoteToQueenFlag))
		moves = append(moves, NewMove(square+7, square, PromoteToRookFlag))
		moves = append(moves, NewMove(square+7, square, PromoteToBishopFlag))
		moves = append(moves, NewMove(square+7, square, PromoteToKnightFlag))
	}

	return moves
}

// gets all knight moves for the current position
func (b *Board) GenerateKnightMoves() []Move {
	return nil
}

// gets all bishop moves for the current position
func (b *Board) GenerateBishopMoves() []Move {
	return nil
}

// gets all rook moves for the current position
func (b *Board) GenerateRookMoves() []Move {
	return nil
}

// gets all queen moves for the current position
func (b *Board) GenerateQueenMoves() []Move {
	return append(b.GenerateBishopMoves(), b.GenerateRookMoves()...)
}

// gets all king moves for the current position
func (b *Board) GenerateKingMoves() []Move {
	return nil
}
