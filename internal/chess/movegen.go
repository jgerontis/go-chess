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
	bp := Piece(BLACK | PAWN)
	moves := make([]Move, 0)
	// let's start with the single push
	// single push is when the square in front of the pawn is empty and not a promotion square
	singlePushBoard := b.Bitboards[wp].SLeft(8) & (b.Bitboards[wp].Not() & b.Bitboards[bp].Not()) & ^Rank8
	// double push is when the pawn is on it's starting square and neither square in front is occupied
	doublePushBoard := b.Bitboards[wp].SLeft(16) & (b.Bitboards[wp].Not() & b.Bitboards[bp].Not()) & ^Rank2
	for square := range 64 {
		if singlePushBoard.Occupied(square) {
			moves = append(moves, NewMove(square-8, square, 0))
		}
		if doublePushBoard.Occupied(square) {
			moves = append(moves, NewMove(square-16, square, 0))
		}
	}
	return moves
}

// gets all black pawn moves
func (b *Board) GenerateBlackPawnMoves() []Move {
	wp := Piece(WHITE | PAWN)
	bp := Piece(BLACK | PAWN)
	moves := make([]Move, 0)
	// single push
	singlePushBoard := b.Bitboards[bp].SRight(8) & (b.Bitboards[wp].Not() & b.Bitboards[bp].Not()) & ^Rank1
	// double push
	doublePushBoard := b.Bitboards[bp].SRight(16) & (b.Bitboards[wp].Not() & b.Bitboards[bp].Not()) & ^Rank7
	for square := range 64 {
		if singlePushBoard.Occupied(square) {
			moves = append(moves, NewMove(square+8, square, 0))
		}
		if doublePushBoard.Occupied(square) {
			moves = append(moves, NewMove(square+16, square, 0))
		}
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
