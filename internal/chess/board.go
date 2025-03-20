package chess

type Board struct {
	// bitboards for each piece type
	Bitboards map[Piece]*Bitboard

	// other board state info, used by FEN
	EnPassantSquare   int
	HalfMoves         int
	FullMoves         int
	BlackCastleRights string
	WhiteCastleRights string
	WhiteToMove       bool

	LegalMoves []Move
}

// create an empty Board object
func NewBoard() *Board {
	return &Board{}
}

// get legal moves for the current position
func (b *Board) GenerateLegalMoves() {
	b.LegalMoves = b.GeneratePawnMoves()
	b.LegalMoves = append(b.LegalMoves, b.GenerateKnightMoves()...)
	b.LegalMoves = append(b.LegalMoves, b.GenerateBishopMoves()...)
	b.LegalMoves = append(b.LegalMoves, b.GenerateRookMoves()...)
	b.LegalMoves = append(b.LegalMoves, b.GenerateQueenMoves()...)
	b.LegalMoves = append(b.LegalMoves, b.GenerateKingMoves()...)
}

// sets a piece in relevant bitboards at the given index
func (b *Board) SetPieceAtIndex(piece Piece, index int) {
	b.Bitboards[piece].Set(index)
	b.Bitboards[Piece(piece.Color())].Set(index)
}

// get the piece at a given square
func (b *Board) GetPieceAtIndex(square int) Piece {
	switch {
	case b.Bitboards[Piece(PAWN|WHITE)].Occupied(square):
		return Piece(PAWN | WHITE)
	case b.Bitboards[Piece(KNIGHT|WHITE)].Occupied(square):
		return Piece(KNIGHT | WHITE)
	case b.Bitboards[Piece(BISHOP|WHITE)].Occupied(square):
		return Piece(BISHOP | WHITE)
	case b.Bitboards[Piece(ROOK|WHITE)].Occupied(square):
		return Piece(ROOK | WHITE)
	case b.Bitboards[Piece(QUEEN|WHITE)].Occupied(square):
		return Piece(QUEEN | WHITE)
	case b.Bitboards[Piece(KING|WHITE)].Occupied(square):
		return Piece(KING | WHITE)
	case b.Bitboards[Piece(PAWN|BLACK)].Occupied(square):
		return Piece(PAWN | BLACK)
	case b.Bitboards[Piece(KNIGHT|BLACK)].Occupied(square):
		return Piece(KNIGHT | BLACK)
	case b.Bitboards[Piece(BISHOP|BLACK)].Occupied(square):
		return Piece(BISHOP | BLACK)
	case b.Bitboards[Piece(ROOK|BLACK)].Occupied(square):
		return Piece(ROOK | BLACK)
	case b.Bitboards[Piece(QUEEN|BLACK)].Occupied(square):
		return Piece(QUEEN | BLACK)
	case b.Bitboards[Piece(KING|BLACK)].Occupied(square):
		return Piece(KING | BLACK)
	default:
		return Piece(0)
	}
}

// updates the board state from a Move object
func (b *Board) MakeMove(move Move) {
	// increment move counter
	b.HalfMoves++
	// increment full move counter on black moves
	if !b.WhiteToMove {
		b.FullMoves++
	}
	// get the original piece
	piece := b.GetPieceAtIndex(move.Source())
	// clear the source square from the piece's board and the color board
	b.Bitboards[piece].Clear(move.Source())
	b.Bitboards[Piece(piece.Color())].Clear(move.Source())
	// update relevant enemy bitboards if it was a capture
	enemyPiece := b.GetPieceAtIndex(move.Target())
	if !enemyPiece.IsNone() {
		b.Bitboards[enemyPiece].Clear(move.Target())
		b.Bitboards[Piece(enemyPiece.Color())].Clear(move.Target())
	}
	// set the original piece on the target square
	b.SetPieceAtIndex(piece, move.Target())
	// change turns
	b.WhiteToMove = !b.WhiteToMove
}
