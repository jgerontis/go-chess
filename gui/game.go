package gui

import (
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jgerontis/go-chess/internal/chess"
)

// Game is going to manage the game state and interface with ebitengine and the chess engine.
type Game struct {
	AudioPlayer *AudioPlayer
	Background  *ebiten.Image
	Board       *chess.Board
	Dragging    bool
	PieceImages map[string]*ebiten.Image
	Player1     string
	Player2     string
	Selected    int
	PrevMove    chess.Move
}

func NewGame(FEN string) *Game {
	background, err := generateBackground()
	if err != nil {
		panic(err)
	}

	images, err := loadPieceImages()
	if err != nil {
		panic(err)
	}

	board := chess.NewBoard()
	board.LoadFEN(FEN)
	board.GenerateLegalMoves()

	return &Game{
		Board:       board,
		PieceImages: images,
		Selected:    -1,
		Dragging:    false,
		Background:  background,
		PrevMove:    0,
	}
}

func (g *Game) MakeMove(move chess.Move) {
	g.Board.MakeMove(move)
	// only generate legal moves when a move is made
	g.Board.GenerateLegalMoves()
}

func (g *Game) Update() error {
	// start by getting the mouse coordinates
	x, y := ebiten.CursorPosition()
	rank, file := g.mouseCoordsToBoardCoords(x, y)
	// do nothing if the mouse is off the board
	if rank < 0 || file < 0 && g.Dragging {
		g.Dragging = false
		return nil
	}

	// see what index we're hovering over
	hovIdx := rank*8 + file
	// if there's a piece at the hov idx, change the cursor to a pointer
	if !g.Board.GetPieceAtIndex(hovIdx).IsNone() && !g.Dragging {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else if g.Dragging {
		ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	// if we just clicked
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Dragging = true
		// let's see what we clicked
		piece := g.Board.GetPieceAtIndex(hovIdx)
		switch {
		// we clicked on an empty square, or an enemy square with a piece already selected
		case piece.IsNone() || (!piece.CanMove(g.Board.WhiteToMove) && g.Selected != -1):
			// we're trying to make a move
			moveAttempt := chess.NewMove(g.Selected, hovIdx, 0)
			// if moveAttempt in legalmoves
			if g.IsLegalMove(moveAttempt) {
				g.MakeMove(moveAttempt)
				return nil
			}
			// if it's an illegal move, unselect the piece
			g.Selected = -1
			return nil
		// we clicked on an ally piece, select it
		case piece.CanMove(g.Board.WhiteToMove):
			g.Selected = hovIdx
			return nil
		// we clicked on an enemy piece, select it
		case !piece.CanMove(g.Board.WhiteToMove):
			g.Selected = hovIdx
			return nil
		}
	}
	// if we just let go
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// since we can only drag our own pieces, we can just try to move the piece
		if g.Dragging && hovIdx != g.Selected {
			// we're trying to make a move
			moveAttempt := chess.NewMove(g.Selected, hovIdx, 0)
			if g.IsLegalMove(moveAttempt) {
				g.MakeMove(moveAttempt)
				g.Selected = -1
			}
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
			fmt.Printf("%s│", g.Board.GetPieceAtIndex(idx))
		}
		if row > 0 {
			fmt.Println("\n │─┼─┼─┼─┼─┼─┼─┼─│")
		}
	}
	fmt.Println("\n └─┴─┴─┴─┴─┴─┴─┴─┘")
	fmt.Println("  a b c d e f g h")
}

func (g *Game) IsLegalMove(move chess.Move) bool {
	return slices.Contains(g.Board.LegalMoves, move)
}
