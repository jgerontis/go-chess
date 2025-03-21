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

	// en passant is when there is pawn on the 4th rank and the last move was a double push

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
	var colorToMove byte
	if b.WhiteToMove {
		colorToMove = WHITE
	} else {
		colorToMove = BLACK
	}
	knigthtBitboard := *b.Bitboards[Piece(colorToMove|KNIGHT)]
	moves := make([]Move, 0)
	for knigthtBitboard != 0 {
		fromSquare := knigthtBitboard.PopLSB()
		knightMoves := KnightMasks[fromSquare] & ^*b.Bitboards[Piece(colorToMove)]
		for knightMoves != 0 {
			toSquare := knightMoves.PopLSB()
			moves = append(moves, NewMove(fromSquare, toSquare, 0))
		}
	}
	return moves
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
	allPieces := *b.Bitboards[Piece(WHITE)] | *b.Bitboards[Piece(BLACK)]
	occ := allPieces & RookMasks[pos]
	occ *= RookMagics[pos]
	occ >>= (64 - RookShifts[pos])
	// get the legal moves for the rook at the given position
	legalMoves := RookAttacks[pos][occ] & ^*b.Bitboards[Piece(colorToMove)]
	moves := make([]Move, 0)
	for legalMoves != 0 {
		toSquare := legalMoves.PopLSB()
		// add the move to the moves slice
		moves = append(moves, NewMove(pos, toSquare, 0))
	}
	return moves
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
	allPieces := *b.Bitboards[Piece(WHITE)] | *b.Bitboards[Piece(BLACK)]
	blockers := allPieces & BishopMasks[pos]
	// use magic to get the index
	index := (blockers * BishopMagics[pos]) >> (64 - BishopShifts[pos])
	// get the legal moves for the rook at the given position
	legalMoves := BishopAttacks[pos][index] & ^*b.Bitboards[Piece(colorToMove)]
	moves := make([]Move, 0)
	for legalMoves != 0 {
		toSquare := legalMoves.PopLSB()
		// add the move to the moves slice
		moves = append(moves, NewMove(pos, toSquare, 0))
	}
	return moves
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
	bishopBitboard := *b.Bitboards[Piece(colorToMove|BISHOP)]
	if bishopBitboard == 0 {
		return []Move{}
	}
	moves := make([]Move, 0)
	for bishopBitboard != 0 {
		fromSquare := bishopBitboard.PopLSB()
		bishopMoves := b.GenerateBishopMovesAtPos(fromSquare)
		moves = append(moves, bishopMoves...)
	}
	return moves
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
	rookBitboard := *b.Bitboards[Piece(colorToMove|ROOK)]
	if rookBitboard == 0 {
		return []Move{}
	}
	moves := make([]Move, 0)
	for rookBitboard != 0 {
		fromSquare := rookBitboard.PopLSB()
		rookMoves := b.GenerateRookMovesAtPos(fromSquare)
		moves = append(moves, rookMoves...)
	}
	return moves
}

// gets all queen moves for the current position
func (b *Board) GenerateQueenMoves() []Move {
	var queenBitboard Bitboard
	if b.WhiteToMove {
		queenBitboard = *b.Bitboards[Piece(WHITE|QUEEN)]
	} else {
		queenBitboard = *b.Bitboards[Piece(BLACK|QUEEN)]
	}
	if queenBitboard == 0 {
		return []Move{}
	}
	queenPos := queenBitboard.PopLSB()
	// queen moves are just the combination of bishop and rook moves
	return append(b.GenerateRookMovesAtPos(queenPos), b.GenerateBishopMovesAtPos(queenPos)...)
}

// gets all king moves for the current position
func (b *Board) GenerateKingMoves() []Move {
	var colorToMove byte
	if b.WhiteToMove {
		colorToMove = WHITE
	} else {
		colorToMove = BLACK
	}
	kingBitboard := *b.Bitboards[Piece(colorToMove|KING)]
	kingPos := kingBitboard.PopLSB()
	kingMoves := KingMasks[kingPos] & ^*b.Bitboards[Piece(colorToMove)]
	moves := make([]Move, 0)
	for kingMoves != 0 {
		toSquare := kingMoves.PopLSB()
		moves = append(moves, NewMove(kingPos, toSquare, 0))
	}
	// castling
	if b.WhiteToMove {
		// white king side castle
		K := strings.Contains(b.WhiteCastleRights, "K")
		if K && !b.Bitboards[Piece(WHITE)].Occupied(5) && !b.Bitboards[Piece(WHITE)].Occupied(6) {
			moves = append(moves, NewMove(4, 6, CastleFlag))
		}
		// white queen side castle
		Q := strings.Contains(b.WhiteCastleRights, "Q")
		if Q && !b.Bitboards[Piece(WHITE)].Occupied(1) && !b.Bitboards[Piece(WHITE)].Occupied(2) && !b.Bitboards[Piece(WHITE)].Occupied(3) {
			moves = append(moves, NewMove(4, 2, CastleFlag))
		}
	} else {
		// black king side castle
		k := strings.Contains(b.BlackCastleRights, "k")
		if k && !b.Bitboards[Piece(BLACK)].Occupied(61) && !b.Bitboards[Piece(BLACK)].Occupied(62) {
			moves = append(moves, NewMove(60, 62, CastleFlag))
		}
		// black queen side castle
		q := strings.Contains(b.BlackCastleRights, "q")
		if q && !b.Bitboards[Piece(BLACK)].Occupied(57) && !b.Bitboards[Piece(BLACK)].Occupied(58) && !b.Bitboards[Piece(BLACK)].Occupied(59) {
			moves = append(moves, NewMove(60, 58, CastleFlag))
		}
	}

	return moves
}

// all possible knight moves for each square
var KnightMasks = [64]Bitboard{
	0x20400, 0x50800, 0xa1100, 0x142200, 0x284400, 0x508800, 0xa01000, 0x402000,
	0x2040004, 0x5080008, 0xa110011, 0x14220022, 0x28440044, 0x50880088, 0xa0100010, 0x40200020,
	0x204000402, 0x508000805, 0xa1100110a, 0x1422002214, 0x2844004428, 0x5088008850, 0xa0100010a0, 0x4020002040,
	0x20400040200, 0x50800080500, 0xa1100110a00, 0x142200221400, 0x284400442800, 0x508800885000, 0xa0100010a000, 0x402000204000,
	0x2040004020000, 0x5080008050000, 0xa1100110a0000, 0x14220022140000, 0x28440044280000, 0x50880088500000, 0xa0100010a00000, 0x40200020400000,
	0x204000402000000, 0x508000805000000, 0xa1100110a000000, 0x1422002214000000, 0x2844004428000000, 0x5088008850000000, 0xa0100010a0000000, 0x4020002040000000,
	0x400040200000000, 0x800080500000000, 0x1100110a00000000, 0x2200221400000000, 0x4400442800000000, 0x8800885000000000, 0x100010a000000000, 0x2000204000000000,
	0x4020000000000, 0x8050000000000, 0x110a0000000000, 0x22140000000000, 0x44280000000000, 0x88500000000000, 0x10a00000000000, 0x20400000000000,
}

// all possible king moves for each square
var KingMasks = [64]Bitboard{
	0x302, 0x705, 0xE0A, 0x1C14, 0x3828, 0x7050, 0xE0A0, 0xC040,
	0x30203, 0x70507, 0xE0A0E, 0x1C141C, 0x382838, 0x705070, 0xE0A0E0, 0xC040C0,
	0x3020300, 0x7050700, 0xE0A0E00, 0x1C141C00, 0x38283800, 0x70507000, 0xE0A0E000, 0xC040C000,
	0x302030000, 0x705070000, 0xE0A0E0000, 0x1C141C0000, 0x3828380000, 0x7050700000, 0xE0A0E00000, 0xC040C00000,
	0x30203000000, 0x70507000000, 0xE0A0E000000, 0x1C141C000000, 0x382838000000, 0x705070000000, 0xE0A0E0000000, 0xC040C0000000,
	0x3020300000000, 0x7050700000000, 0xE0A0E00000000, 0x1C141C00000000, 0x38283800000000, 0x70507000000000, 0xE0A0E000000000, 0xC040C000000000,
	0x302030000000000, 0x705070000000000, 0xE0A0E0000000000, 0x1C141C0000000000, 0x3828380000000000, 0x7050700000000000, 0xE0A0E00000000000, 0xC040C00000000000,
	0x203000000000000, 0x507000000000000, 0xA0E000000000000, 0x141C000000000000, 0x2838000000000000, 0x5070000000000000, 0xA0E0000000000000, 0x40C0000000000000,
}

var RookMasks [64]Bitboard = [64]Bitboard{
	282578800148862, 565157600297596, 1130315200595066, 2260630401190006, 4521260802379886, 9042521604759646, 18085043209519166, 36170086419038334,
	282578800180736, 565157600328704, 1130315200625152, 2260630401218048, 4521260802403840, 9042521604775424, 18085043209518592, 36170086419037696,
	282578808340736, 565157608292864, 1130315208328192, 2260630408398848, 4521260808540160, 9042521608822784, 18085043209388032, 36170086418907136,
	282580897300736, 565159647117824, 1130317180306432, 2260632246683648, 4521262379438080, 9042522644946944, 18085043175964672, 36170086385483776,
	283115671060736, 565681586307584, 1130822006735872, 2261102847592448, 4521664529305600, 9042787892731904, 18085034619584512, 36170077829103616,
	420017753620736, 699298018886144, 1260057572672512, 2381576680245248, 4624614895390720, 9110691325681664, 18082844186263552, 36167887395782656,
	35466950888980736, 34905104758997504, 34344362452452352, 33222877839362048, 30979908613181440, 26493970160820224, 17522093256097792, 35607136465616896,
	9079539427579068672, 8935706818303361536, 8792156787827803136, 8505056726876686336, 7930856604974452736, 6782456361169985536, 4485655873561051136, 9115426935197958144,
}

var BishopMasks [64]Bitboard = [64]Bitboard{
	18049651735527936, 70506452091904, 275415828992, 1075975168, 38021120, 8657588224, 2216338399232, 567382630219776,
	9024825867763712, 18049651735527424, 70506452221952, 275449643008, 9733406720, 2216342585344, 567382630203392, 1134765260406784,
	4512412933816832, 9024825867633664, 18049651768822272, 70515108615168, 2491752130560, 567383701868544, 1134765256220672, 2269530512441344,
	2256206450263040, 4512412900526080, 9024834391117824, 18051867805491712, 637888545440768, 1135039602493440, 2269529440784384, 4539058881568768,
	1128098963916800, 2256197927833600, 4514594912477184, 9592139778506752, 19184279556981248, 2339762086609920, 4538784537380864, 9077569074761728,
	562958610993152, 1125917221986304, 2814792987328512, 5629586008178688, 11259172008099840, 22518341868716544, 9007336962655232, 18014673925310464,
	2216338399232, 4432676798464, 11064376819712, 22137335185408, 44272556441600, 87995357200384, 35253226045952, 70506452091904,
	567382630219776, 1134765260406784, 2832480465846272, 5667157807464448, 11333774449049600, 22526811443298304, 9024825867763712, 18049651735527936,
}
