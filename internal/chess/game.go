package chess

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
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

	// if we're holding down the mouse, turn dragging on
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.Dragging {
			piece := g.Board.GetPiece(rank, file)
			if !piece.IsNone() {
				g.Selected = rank*8 + file
				g.Dragging = true
			} else {
				g.Selected = -1
			}
		}
		// if we just clicked on another square with a valid move, mack the move

	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.Dragging {
			// move the piece
			target := rank*8 + file
			if g.Selected != rank*8+file {
				g.Board.MovePiece(g.Selected, target)
				g.WhiteToMove = !g.WhiteToMove
			}

			g.Dragging = false
			g.Selected = -1
		}
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
