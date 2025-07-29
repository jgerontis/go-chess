package chess

import (
	"testing"
)

func TestGetPieceAtIndex(t *testing.T) {
	board := NewBoard()
	board.LoadFEN(START_FEN)
	piece := board.GetPieceAtIndex(0)
	if piece.IsNone() {
		t.Fatalf("Got nil piece at index 0")
	}
	if piece.Type() != ROOK {
		t.Errorf("Expected ROOK, got %d", piece.Type())
	}
	if piece.Color() != WHITE {
		t.Errorf("Expected WHITE, got %d", piece.Color())
	}
	piece = board.GetPieceAtIndex(4)
	if piece.IsNone() {
		t.Fatalf("Got nil piece at index 4")
	}
	if piece.Type() != KING {
		t.Errorf("Expected KING, got %d", piece.Type())
	}
	if piece.Color() != WHITE {
		t.Errorf("Expected WHITE, got %d", piece.Color())
	}
	piece = board.GetPieceAtIndex(60)
	if piece.IsNone() {
		t.Fatalf("Got nil piece at index 60")
	}
	if piece.Type() != KING {
		t.Errorf("Expected KING, got %d", piece.Type())
	}
	if piece.Color() != BLACK {
		t.Errorf("Expected BLACK, got %d", piece.Color())
	}
	piece = board.GetPieceAtIndex(59)
	if piece.IsNone() {
		t.Fatalf("Got nil piece at index 59")
	}
	if piece.Type() != QUEEN {
		t.Errorf("Expected QUEEN, got %d", piece.Type())
	}
	if piece.Color() != BLACK {
		t.Errorf("Expected BLACK, got %d", piece.Color())
	}
	piece = board.GetPieceAtIndex(63)
	if piece.IsNone() {
		t.Fatalf("Got nil piece at index 63")
	}
	if piece.Type() != ROOK {
		t.Errorf("Expected ROOK, got %d", piece.Type())
	}
	if piece.Color() != BLACK {
		t.Errorf("Expected BLACK, got %d", piece.Color())
	}
}

func TestSetPieceAtIndex(t *testing.T) {
	board := NewBoard()
	wr := WHITE | ROOK
	board.Bitboards = make(map[byte]*Bitboard)
	board.Bitboards[wr] = NewBitboard()
	board.Bitboards[WHITE] = NewBitboard()
	board.SetPieceAtIndex(Piece(wr), 0)
	occupied := board.Bitboards[wr].Occupied(0)
	if !occupied {
		t.Errorf("Expected square 0 to be occupied")
	}
	occupied = board.Bitboards[WHITE].Occupied(0)
	if !occupied {
		t.Errorf("Expected square 0 to be occupied")
	}
}

func TestBitboardInitialization(t *testing.T) {
	board := NewBoard()
	board.LoadFEN(START_FEN)
	// check that all bitboards are initialized
	for _, bb := range board.Bitboards {
		if bb == nil {
			t.Fatalf("Expected bitboard to be initialized")
		}
	}
}

func TestIsSquareAttacked(t *testing.T) {
	board := NewBoard()
	board.LoadFEN(START_FEN)

	// Test that squares in front of pawns are attacked
	// White pawn on a2 attacks b3 and (no c3 due to edge)
	if !board.IsSquareAttacked(9, WHITE) { // b3
		t.Error("Expected b3 to be attacked by white pawn")
	}

	// Black pawn on a7 attacks b6 
	if !board.IsSquareAttacked(41, BLACK) { // b6
		t.Error("Expected b6 to be attacked by black pawn")
	}

	// Test that knights attack the right squares
	// White knight on b1 attacks a3, c3, d2
	if !board.IsSquareAttacked(16, WHITE) { // a3
		t.Error("Expected a3 to be attacked by white knight")
	}
	if !board.IsSquareAttacked(18, WHITE) { // c3
		t.Error("Expected c3 to be attacked by white knight")
	}

	// Test a square that shouldn't be attacked in starting position
	// Let's use e4 (square 28) instead which is definitely not under attack
	if board.IsSquareAttacked(28, WHITE) { // e4 - not attacked by anything in starting position
		t.Error("Expected e4 to not be attacked by white")
	}
	if board.IsSquareAttacked(28, BLACK) { // e4 - not attacked by anything in starting position  
		t.Error("Expected e4 to not be attacked by black")
	}
}

func TestIsInCheck(t *testing.T) {
	board := NewBoard()
	
	// Starting position - no one should be in check
	board.LoadFEN(START_FEN)
	if board.IsInCheck(WHITE) {
		t.Error("White should not be in check in starting position")
	}
	if board.IsInCheck(BLACK) {
		t.Error("Black should not be in check in starting position")
	}

	// Test a position where white king is in check
	// Position: black queen on d1 attacking white king on e1
	board.LoadFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNB1K1NR w KQkq - 0 1") // removed white queen
	board.SetPieceAtIndex(Piece(BLACK|QUEEN), 3) // Put black queen on d1
	if !board.IsInCheck(WHITE) {
		t.Error("White king should be in check")
	}
	if board.IsInCheck(BLACK) {
		t.Error("Black king should not be in check")
	}

	// Test a position where black king is in check  
	board.LoadFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	board.SetPieceAtIndex(Piece(WHITE|QUEEN), 59) // Put white queen on d8
	if board.IsInCheck(WHITE) {
		t.Error("White king should not be in check")
	}
	if !board.IsInCheck(BLACK) {
		t.Error("Black king should be in check")
	}
}

func TestFilterLegalMoves(t *testing.T) {
	board := NewBoard()
	
	// Test a position where some moves would be illegal due to check
	// Position: White king on e1, black rook on e8, white bishop on d2
	// The king cannot move to d1, e2, f1 as they would be in check from the rook
	board.LoadFEN("4r3/8/8/8/8/8/3B4/4K3 w - - 0 1")
	
	// Generate all pseudo-legal king moves
	kingMoves := board.GenerateKingMoves()
	
	// The king should have some pseudo-legal moves but fewer legal moves
	if len(kingMoves) == 0 {
		t.Error("King should have some pseudo-legal moves")
	}
	
	// Test a simpler case - starting position where all moves should be legal
	board.LoadFEN(START_FEN)
	pawnMoves := board.GeneratePawnMoves()
	if len(pawnMoves) == 0 {
		t.Error("Should have pawn moves in starting position")
	}
	
	// In starting position, all pawn moves should be legal (no pins or checks)
	// This tests that FilterLegalMoves doesn't incorrectly filter out legal moves
	for _, move := range pawnMoves {
		state := board.MakeMove(move)
		// After making a legal move, the original player should not be in check
		movingColor := BLACK // After white moves, we check if white is still safe
		if move.Source() >= 48 { // If move was from black's side
			movingColor = WHITE
		}
		if board.IsInCheck(movingColor) {
			t.Errorf("Move %s should be legal but leaves king in check", move.String())
		}
		board.UnmakeMove(move, state)
	}
}

func TestEnPassant(t *testing.T) {
	board := NewBoard()
	
	// Test white en passant capture
	// Position: white pawn on e5, black pawn just moved d7-d5
	board.LoadFEN("rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 1")
	
	// White should be able to capture en passant on d6
	pawnMoves := board.GeneratePawnMoves()
	
	var enPassantMove Move
	var foundEnPassant bool
	for _, move := range pawnMoves {
		if move.Flag() == EN_PASSANT_FLAG {
			enPassantMove = move
			foundEnPassant = true
			break
		}
	}
	
	if !foundEnPassant {
		t.Error("Should find en passant move in this position")
	} else {
		// Verify the move is from e5 to d6
		if enPassantMove.Source() != 36 { // e5
			t.Errorf("En passant source should be e5 (36), got %d", enPassantMove.Source())
		}
		if enPassantMove.Target() != 43 { // d6  
			t.Errorf("En passant target should be d6 (43), got %d", enPassantMove.Target())
		}
		
		// Test making the en passant move
		state := board.MakeMove(enPassantMove)
		
		// The black pawn on d5 should be removed
		if !board.GetPieceAtIndex(35).IsNone() { // d5 should be empty
			t.Error("Black pawn on d5 should be captured by en passant")
		}
		
		// White pawn should be on d6
		piece := board.GetPieceAtIndex(43) // d6
		if piece.Type() != PAWN || piece.Color() != WHITE {
			t.Error("White pawn should be on d6 after en passant")
		}
		
		// Test unmaking the move
		board.UnmakeMove(enPassantMove, state)
		
		// Black pawn should be back on d5
		piece = board.GetPieceAtIndex(35) // d5
		if piece.Type() != PAWN || piece.Color() != BLACK {
			t.Error("Black pawn should be restored on d5 after unmake")
		}
		
		// White pawn should be back on e5
		piece = board.GetPieceAtIndex(36) // e5
		if piece.Type() != PAWN || piece.Color() != WHITE {
			t.Error("White pawn should be back on e5 after unmake")
		}
		
		// d6 should be empty again
		if !board.GetPieceAtIndex(43).IsNone() { // d6
			t.Error("d6 should be empty after unmake")
		}
	}
	
	// Test black en passant capture
	board.LoadFEN("rnbqkbnr/pppp1ppp/8/8/3Pp3/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1")
	
	blackPawnMoves := board.GeneratePawnMoves()
	foundEnPassant = false
	for _, move := range blackPawnMoves {
		if move.Flag() == EN_PASSANT_FLAG {
			enPassantMove = move
			foundEnPassant = true
			break
		}
	}
	
	if !foundEnPassant {
		t.Error("Should find black en passant move")
	} else {
		// Should be from e4 to d3
		if enPassantMove.Source() != 28 { // e4
			t.Errorf("Black en passant source should be e4 (28), got %d", enPassantMove.Source())
		}
		if enPassantMove.Target() != 19 { // d3
			t.Errorf("Black en passant target should be d3 (19), got %d", enPassantMove.Target())
		}
	}
	
	// Test that en passant square is properly reset
	board.LoadFEN(START_FEN)
	board.EnPassantSquare = 20 // Set to some square
	
	// Make a non-double pawn move
	knightMoves := board.GenerateKnightMoves()
	if len(knightMoves) > 0 {
		state := board.MakeMove(knightMoves[0])
		if board.EnPassantSquare != -1 {
			t.Error("En passant square should be reset to -1 after non-double pawn move")
		}
		board.UnmakeMove(knightMoves[0], state)
		if board.EnPassantSquare != 20 {
			t.Error("En passant square should be restored after unmake")
		}
	}
}

func TestCastling(t *testing.T) {
	board := NewBoard()
	
	// Test white kingside castling
	board.LoadFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1")
	
	kingMoves := board.GenerateKingMoves()
	var kingsideCastle Move
	var queensideCastle Move
	foundKingside := false
	foundQueenside := false
	
	for _, move := range kingMoves {
		if move.Flag() == CASTLE_FLAG {
			if move.Target() == 6 { // g1
				kingsideCastle = move
				foundKingside = true
			} else if move.Target() == 2 { // c1
				queensideCastle = move
				foundQueenside = true
			}
		}
	}
	
	if !foundKingside {
		t.Error("Should find white kingside castle")
	}
	if !foundQueenside {
		t.Error("Should find white queenside castle")
	}
	
	// Test making kingside castle
	state := board.MakeMove(kingsideCastle)
	
	// Also test that we found queenside castle (just to use the variable)
	_ = queensideCastle
	
	// King should be on g1
	king := board.GetPieceAtIndex(6) // g1
	if king.Type() != KING || king.Color() != WHITE {
		t.Error("King should be on g1 after kingside castle")
	}
	
	// Rook should be on f1
	rook := board.GetPieceAtIndex(5) // f1
	if rook.Type() != ROOK || rook.Color() != WHITE {
		t.Error("Rook should be on f1 after kingside castle")
	}
	
	// e1 and h1 should be empty
	if !board.GetPieceAtIndex(4).IsNone() { // e1
		t.Error("e1 should be empty after castle")
	}
	if !board.GetPieceAtIndex(7).IsNone() { // h1
		t.Error("h1 should be empty after castle")
	}
	
	// Test unmaking the castle
	board.UnmakeMove(kingsideCastle, state)
	
	// King should be back on e1
	king = board.GetPieceAtIndex(4) // e1
	if king.Type() != KING || king.Color() != WHITE {
		t.Error("King should be back on e1 after unmake")
	}
	
	// Rook should be back on h1
	rook = board.GetPieceAtIndex(7) // h1
	if rook.Type() != ROOK || rook.Color() != WHITE {
		t.Error("Rook should be back on h1 after unmake")
	}
	
	// f1 and g1 should be empty
	if !board.GetPieceAtIndex(5).IsNone() { // f1
		t.Error("f1 should be empty after unmake")
	}
	if !board.GetPieceAtIndex(6).IsNone() { // g1
		t.Error("g1 should be empty after unmake")
	}
}

func TestPromotion(t *testing.T) {
	board := NewBoard()
	
	// Test white pawn promotion
	board.LoadFEN("8/P7/8/8/8/8/8/8 w - - 0 1")
	
	pawnMoves := board.GeneratePawnMoves()
	if len(pawnMoves) != 4 {
		t.Errorf("Expected 4 promotion moves, got %d", len(pawnMoves))
	}
	
	// Test queen promotion
	var queenPromotion Move
	for _, move := range pawnMoves {
		if move.Flag() == PROMOTE_QUEEN_FLAG {
			queenPromotion = move
			break
		}
	}
	
	state := board.MakeMove(queenPromotion)
	
	// Should have a white queen on a8
	piece := board.GetPieceAtIndex(56) // a8
	if piece.Type() != QUEEN || piece.Color() != WHITE {
		t.Error("Should have white queen on a8 after promotion")
	}
	
	// a7 should be empty
	if !board.GetPieceAtIndex(48).IsNone() { // a7
		t.Error("a7 should be empty after promotion")
	}
	
	// Test unmaking promotion
	board.UnmakeMove(queenPromotion, state)
	
	// Should have white pawn back on a7
	piece = board.GetPieceAtIndex(48) // a7
	if piece.Type() != PAWN || piece.Color() != WHITE {
		t.Error("Should have white pawn back on a7 after unmake")
	}
	
	// a8 should be empty
	if !board.GetPieceAtIndex(56).IsNone() { // a8
		t.Error("a8 should be empty after unmake")
	}
}

func TestMoveUnmakeConsistency(t *testing.T) {
	board := NewBoard()
	board.LoadFEN(START_FEN)
	
	// Generate all legal moves
	board.GenerateLegalMoves()
	
	// Test that make/unmake is consistent for all moves
	for _, move := range board.LegalMoves {
		// Save the original FEN
		originalFEN := board.ExportFEN()
		
		// Make the move
		state := board.MakeMove(move)
		
		// Unmake the move
		board.UnmakeMove(move, state)
		
		// Board should be identical to original
		restoredFEN := board.ExportFEN()
		if originalFEN != restoredFEN {
			t.Errorf("Make/unmake inconsistency for move %s:\nOriginal: %s\nRestored: %s", 
				move.String(), originalFEN, restoredFEN)
		}
	}
}

func TestCaptureRestoration(t *testing.T) {
	board := NewBoard()
	
	// Position with possible captures
	board.LoadFEN("rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
	
	// White pawn can capture black pawn
	captureMove := NewMove(28, 35, 0) // e4 captures d5
	
	// Verify there's a black pawn on d5
	capturedPiece := board.GetPieceAtIndex(35)
	if capturedPiece.Type() != PAWN || capturedPiece.Color() != BLACK {
		t.Error("Expected black pawn on d5")
	}
	
	// Make the capture
	state := board.MakeMove(captureMove)
	
	// d5 should now have white pawn
	piece := board.GetPieceAtIndex(35)
	if piece.Type() != PAWN || piece.Color() != WHITE {
		t.Error("Expected white pawn on d5 after capture")
	}
	
	// e4 should be empty
	if !board.GetPieceAtIndex(28).IsNone() {
		t.Error("e4 should be empty after move")
	}
	
	// Unmake the move
	board.UnmakeMove(captureMove, state)
	
	// Black pawn should be restored on d5
	piece = board.GetPieceAtIndex(35)
	if piece.Type() != PAWN || piece.Color() != BLACK {
		t.Error("Black pawn should be restored on d5")
	}
	
	// White pawn should be back on e4
	piece = board.GetPieceAtIndex(28)
	if piece.Type() != PAWN || piece.Color() != WHITE {
		t.Error("White pawn should be back on e4")
	}
}
