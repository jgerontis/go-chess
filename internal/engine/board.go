package engine

type Board struct {
	// bitboards for each piece type
	Pawns   Bitboard
	Knights Bitboard
	Bishops Bitboard
	Rooks   Bitboard
	Queens  Bitboard
	Kings   Bitboard
	Black   Bitboard
	White   Bitboard
}
