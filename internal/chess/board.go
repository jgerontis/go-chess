package chess

type Board struct {
	// bitboards for each piece type
	Bitboards map[byte]*Bitboard

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
	b.Bitboards[byte(piece)].Set(index)
	b.Bitboards[piece.Color()].Set(index)
}

func (b *Board) ClearPieceAtIndex(piece Piece, index int) {
	b.Bitboards[byte(piece)].Clear(index)
	b.Bitboards[piece.Color()].Clear(index)
}

// get the type of piece at a given square
func (b *Board) GetPieceAtIndex(square int) Piece {
	switch {
	case b.Bitboards[PAWN|WHITE].Occupied(square):
		return Piece(PAWN | WHITE)
	case b.Bitboards[KNIGHT|WHITE].Occupied(square):
		return Piece(KNIGHT | WHITE)
	case b.Bitboards[BISHOP|WHITE].Occupied(square):
		return Piece(BISHOP | WHITE)
	case b.Bitboards[ROOK|WHITE].Occupied(square):
		return Piece(ROOK | WHITE)
	case b.Bitboards[QUEEN|WHITE].Occupied(square):
		return Piece(QUEEN | WHITE)
	case b.Bitboards[KING|WHITE].Occupied(square):
		return Piece(KING | WHITE)
	case b.Bitboards[PAWN|BLACK].Occupied(square):
		return Piece(PAWN | BLACK)
	case b.Bitboards[KNIGHT|BLACK].Occupied(square):
		return Piece(KNIGHT | BLACK)
	case b.Bitboards[BISHOP|BLACK].Occupied(square):
		return Piece(BISHOP | BLACK)
	case b.Bitboards[ROOK|BLACK].Occupied(square):
		return Piece(ROOK | BLACK)
	case b.Bitboards[QUEEN|BLACK].Occupied(square):
		return Piece(QUEEN | BLACK)
	case b.Bitboards[KING|BLACK].Occupied(square):
		return Piece(KING | BLACK)
	default:
		return Piece(0)
	}
}

// updates the board state from a Move object and returns the previous state
func (b *Board) MakeMove(move Move) BoardState {
	// Save the current state before making changes
	state := b.SaveState()

	// get the original piece
	piece := b.GetPieceAtIndex(move.Source())
	// update relevant enemy bitboards if it was a capture
	enemyPiece := b.GetPieceAtIndex(move.Target())
	state.CapturedPiece = enemyPiece // Save captured piece for unmake
	if !enemyPiece.IsNone() {
		b.ClearPieceAtIndex(enemyPiece, move.Target())
	}

	// increment move counter
	b.HalfMoves++
	// increment full move counter on black moves
	if !b.WhiteToMove {
		b.FullMoves++
	}

	// clear the source square from the piece's board and the color board
	b.ClearPieceAtIndex(piece, move.Source())
	// set the original piece on the target square
	b.SetPieceAtIndex(piece, move.Target())

	switch move.Flag() {
	case CASTLE_FLAG:
		// move the rook
		if b.WhiteToMove {
			// kingside
			if move.Target() == 6 {
				b.ClearPieceAtIndex(Piece(ROOK|WHITE), 7)
				b.SetPieceAtIndex(Piece(ROOK|WHITE), 5)
			} else {
				// queenside
				b.ClearPieceAtIndex(Piece(ROOK|WHITE), 0)
				b.SetPieceAtIndex(Piece(ROOK|WHITE), 3)
			}
		} else {
			// kingside
			if move.Target() == 62 {
				b.ClearPieceAtIndex(Piece(ROOK|BLACK), 63)
				b.SetPieceAtIndex(Piece(ROOK|BLACK), 61)
			} else {
				// queenside
				b.ClearPieceAtIndex(Piece(ROOK|BLACK), 56)
				b.SetPieceAtIndex(Piece(ROOK|BLACK), 59)
			}
		}
	case PROMOTE_KNIGHT_FLAG:
		b.ClearPieceAtIndex(piece, move.Target())
		b.SetPieceAtIndex(Piece(KNIGHT|piece.Color()), move.Target())
	case PROMOTE_BISHOP_FLAG:
		b.ClearPieceAtIndex(piece, move.Target())
		b.SetPieceAtIndex(Piece(BISHOP|piece.Color()), move.Target())
	case PROMOTE_ROOK_FLAG:
		b.ClearPieceAtIndex(piece, move.Target())
		b.SetPieceAtIndex(Piece(ROOK|piece.Color()), move.Target())
	case PROMOTE_QUEEN_FLAG:
		b.ClearPieceAtIndex(piece, move.Target())
		b.SetPieceAtIndex(Piece(QUEEN|piece.Color()), move.Target())
	case EN_PASSANT_FLAG:
		// clear the captured pawn
		if state.WhiteToMove { // Use the original turn state
			b.ClearPieceAtIndex(Piece(PAWN|BLACK), move.Target()-8)
		} else {
			b.ClearPieceAtIndex(Piece(PAWN|WHITE), move.Target()+8)
		}
	case PAWN_DOUBLE_FLAG:
		break
	default:
		break
	}

	// set the en passant square if it was a double pawn move
	if move.Flag() == PAWN_DOUBLE_FLAG {
		if state.WhiteToMove { // Use the original turn state
			b.EnPassantSquare = move.Target() - 8
		} else {
			b.EnPassantSquare = move.Target() + 8
		}
	} else {
		b.EnPassantSquare = -1
	}

	// change turns
	b.WhiteToMove = !b.WhiteToMove

	return state
}

// BoardState represents the board state that needs to be restored when unmaking a move
type BoardState struct {
	WhiteToMove       bool
	EnPassantSquare   int
	HalfMoves         int
	FullMoves         int
	BlackCastleRights string
	WhiteCastleRights string
	CapturedPiece     Piece
}

// SaveState saves the current board state before making a move
func (b *Board) SaveState() BoardState {
	return BoardState{
		WhiteToMove:       b.WhiteToMove,
		EnPassantSquare:   b.EnPassantSquare,
		HalfMoves:         b.HalfMoves,
		FullMoves:         b.FullMoves,
		BlackCastleRights: b.BlackCastleRights,
		WhiteCastleRights: b.WhiteCastleRights,
	}
}

// RestoreState restores the board state after unmaking a move
func (b *Board) RestoreState(state BoardState) {
	b.WhiteToMove = state.WhiteToMove
	b.EnPassantSquare = state.EnPassantSquare
	b.HalfMoves = state.HalfMoves
	b.FullMoves = state.FullMoves
	b.BlackCastleRights = state.BlackCastleRights
	b.WhiteCastleRights = state.WhiteCastleRights
}

// UnmakeMove reverses a move that was previously made
func (b *Board) UnmakeMove(move Move, state BoardState) {
	// Get the piece that was moved (now at target square)
	piece := b.GetPieceAtIndex(move.Target())

	// Handle special moves first
	switch move.Flag() {
	case CASTLE_FLAG:
		// Undo castling - move king back and rook back
		if !state.WhiteToMove { // Was white's move
			if move.Target() == 6 { // Kingside
				b.ClearPieceAtIndex(Piece(ROOK|WHITE), 5)
				b.SetPieceAtIndex(Piece(ROOK|WHITE), 7)
			} else { // Queenside
				b.ClearPieceAtIndex(Piece(ROOK|WHITE), 3)
				b.SetPieceAtIndex(Piece(ROOK|WHITE), 0)
			}
		} else { // Was black's move
			if move.Target() == 62 { // Kingside
				b.ClearPieceAtIndex(Piece(ROOK|BLACK), 61)
				b.SetPieceAtIndex(Piece(ROOK|BLACK), 63)
			} else { // Queenside
				b.ClearPieceAtIndex(Piece(ROOK|BLACK), 59)
				b.SetPieceAtIndex(Piece(ROOK|BLACK), 56)
			}
		}
	case PROMOTE_KNIGHT_FLAG, PROMOTE_BISHOP_FLAG, PROMOTE_ROOK_FLAG, PROMOTE_QUEEN_FLAG:
		// Clear the promoted piece and restore the original pawn
		b.ClearPieceAtIndex(piece, move.Target())
		piece = Piece(PAWN | piece.Color())
	case EN_PASSANT_FLAG:
		// Restore the captured pawn
		if state.WhiteToMove { // Was white's move
			b.SetPieceAtIndex(Piece(PAWN|BLACK), move.Target()-8)
		} else { // Was black's move
			b.SetPieceAtIndex(Piece(PAWN|WHITE), move.Target()+8)
		}
	}

	// Move the piece back to its original square
	b.ClearPieceAtIndex(piece, move.Target())
	b.SetPieceAtIndex(piece, move.Source())

	// Restore captured piece if there was one
	if !state.CapturedPiece.IsNone() {
		b.SetPieceAtIndex(state.CapturedPiece, move.Target())
	}

	// Restore the board state
	b.RestoreState(state)
}
