package engine

import "github.com/jgerontis/go-chess/internal/chess"

func FindBestMove(board *chess.Board) chess.Move {
	return chess.NewMove(chess.StringToSquare("e2"), chess.StringToSquare("e4"), 0)
}
