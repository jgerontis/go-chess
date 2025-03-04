package chess

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Board       *Board
	PieceImages map[Piece]*ebiten.Image
	WhiteToMove bool
	Selected    int
	Dragging    bool
	Background  *ebiten.Image
}

func NewGame(fen string) *Game {
	images, err := loadPieceImages()
	if err != nil {
		panic(err)
	}

	background, err := GenerateBackground()
	if err != nil {
		panic(err)
	}

	board, err := NewBoardFromFEN(fen)
	if err != nil {
		panic(err)
	}
	return &Game{
		Board:       board,
		PieceImages: images,
		WhiteToMove: true,
		Selected:    -1,
		Dragging:    false,
		Background:  background,
	}
}

func (g *Game) Update() error {
	// start by getting the mouse coordinates
	x, y := ebiten.CursorPosition()
	rank, file := g.MouseCoordsToBoardCoords(x, y)
	// do nothing if the mouse is off the board
	if rank < 0 || file < 0 {
		g.Dragging = false
		return nil
	}

	// see what index we're hovering over
	hovIdx := rank*8 + file
	// if there's a piece at the hov idx, change the cursor to a pointer
	if !g.Board[hovIdx].IsNone() && !g.Dragging {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else if g.Dragging {
		ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	// if we just clicked
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// see what we clicked on
		newPiece := g.Board.GetPiece(rank, file)
		// if I clicked on an empty square, I'm either moving a piece or clearing the selection
		if newPiece.IsNone() {
			// if we already had a piece selected and it's our piece
			if g.Selected != -1 && g.Board[g.Selected].CanMove(g.WhiteToMove) {
				// move the piece
				g.Board.MovePiece(g.Selected, hovIdx)
				g.Selected = -1
				g.WhiteToMove = !g.WhiteToMove
				return nil
			}
			// if we had nothing selected, then clear the selection
			g.Selected = -1
			return nil
		}
		// if we clicked on a piece
		if !newPiece.IsNone() {
			// if it's an enemy piece, 3 scenarios
			if !newPiece.CanMove(g.WhiteToMove) {
				// 1. we had nothing selected, just select the enemy piece
				if g.Selected == -1 {
					g.Selected = hovIdx
					return nil
				}
				// 2. we had an enemy piece selected, just select the new enemy piece
				if !g.Board[g.Selected].CanMove(g.WhiteToMove) {
					g.Selected = hovIdx
					return nil
				}
				// 3. we had our piece selected, capture the enemy piece
				if g.Board[g.Selected].CanMove(g.WhiteToMove) {
					g.Board.MovePiece(g.Selected, hovIdx)
					g.Selected = -1
					g.WhiteToMove = !g.WhiteToMove
					return nil
				}
			}
			if newPiece.CanMove(g.WhiteToMove) {
				// we clicked on our own piece
				g.Selected = hovIdx
				g.Dragging = true
				return nil
			}
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// since we can only drag our own pieces, we can just move the piece
		if g.Dragging && hovIdx != g.Selected {
			g.Board.MovePiece(g.Selected, hovIdx)
			g.Selected = -1
			g.WhiteToMove = !g.WhiteToMove
		}
		g.Dragging = false
	}
	return nil
}

func (g *Game) PrintBoard() {
	fmt.Println("  a b c d e f g h")
	fmt.Println(" ┌─┬─┬─┬─┬─┬─┬─┬─┐")
	for row := 7; row >= 0; row-- {
		fmt.Printf("%d│", row+1)
		for col := range 8 {
			idx := row*8 + col
			fmt.Printf("%s│", g.Board[idx])
		}
		if row > 0 {
			fmt.Println("\n │─┼─┼─┼─┼─┼─┼─┼─│")
		}
	}
	fmt.Println("\n └─┴─┴─┴─┴─┴─┴─┴─┘")
	fmt.Println("  a b c d e f g h")
}

func (g *Game) GetLegalMoves() {
	// TODO
}
