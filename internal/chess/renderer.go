package chess

import (
	"bytes"
	"image"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

const (
	WindowWidth  = 900
	WindowHeight = 900
	BoardSize    = 800
	SquareSize   = BoardSize / 8
)

func (g *Game) Draw(screen *ebiten.Image) {
	// write player turn in the top left
	var turn string
	if g.WhiteToMove {
		turn = "White to move"
	} else {
		turn = "Black to move"
	}
	ebitenutil.DebugPrint(screen, turn)
	margin := (WindowWidth - BoardSize) / 2
	// draw the Background within the margins
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(margin), float64(margin))
	screen.DrawImage(g.Background, opts)
	for rank := range 8 {
		for file := range 8 {
			// get the graphical coordinates of the square
			x := margin + file*SquareSize
			y := margin + (7-rank)*SquareSize

			// check if the square is selected
			squareSelected := g.Selected == rank*8+file

			// draw a red square if the square is selected
			if squareSelected {
				square := ebiten.NewImage(SquareSize, SquareSize)
				squareColor := color.RGBA{255, 0, 0, 255}
				square.Fill(squareColor)
				squareOpts := &ebiten.DrawImageOptions{}
				squareOpts.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(square, squareOpts)
			}

			// get the piece on the square, skip if there is no piece
			piece := g.Board.GetPiece(rank, file)
			if piece.IsNone() {
				continue
			}

			// don't draw the piece in the square if it's being dragged
			if squareSelected && g.Dragging {
				continue
			}

			// draw the piece on the square
			pieceImg := g.PieceImages[piece]
			pieceOpts := &ebiten.DrawImageOptions{}
			// if the piece is being dragged, move it with the mouse
			pieceOpts.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(pieceImg, pieceOpts)
		}
	}
	if g.Selected != -1 && g.Dragging {
		x, y := ebiten.CursorPosition()
		// adjust so we grab the middle of the piece, not the top left
		x -= SquareSize / 2
		y -= SquareSize / 2
		// draw the piece being dragged
		pieceImage := g.PieceImages[g.Board[g.Selected]]
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(pieceImage, opts)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WindowWidth, WindowHeight
}

func loadPieceImages() (map[Piece]*ebiten.Image, error) {
	pieces := make(map[Piece]*ebiten.Image)
	// load the images for the pieces
	pieces[Piece(WHITE|KING)], _ = loadSVG("assets/images/wk.svg", SquareSize, SquareSize)
	pieces[Piece(BLACK|KING)], _ = loadSVG("assets/images/bk.svg", SquareSize, SquareSize)
	pieces[Piece(WHITE|QUEEN)], _ = loadSVG("assets/images/wq.svg", SquareSize, SquareSize)
	pieces[Piece(BLACK|QUEEN)], _ = loadSVG("assets/images/bq.svg", SquareSize, SquareSize)
	pieces[Piece(WHITE|ROOK)], _ = loadSVG("assets/images/wr.svg", SquareSize, SquareSize)
	pieces[Piece(BLACK|ROOK)], _ = loadSVG("assets/images/br.svg", SquareSize, SquareSize)
	pieces[Piece(WHITE|BISHOP)], _ = loadSVG("assets/images/wb.svg", SquareSize, SquareSize)
	pieces[Piece(BLACK|BISHOP)], _ = loadSVG("assets/images/bb.svg", SquareSize, SquareSize)
	pieces[Piece(WHITE|KNIGHT)], _ = loadSVG("assets/images/wn.svg", SquareSize, SquareSize)
	pieces[Piece(BLACK|KNIGHT)], _ = loadSVG("assets/images/bn.svg", SquareSize, SquareSize)
	pieces[Piece(WHITE|PAWN)], _ = loadSVG("assets/images/wp.svg", SquareSize, SquareSize)
	pieces[Piece(BLACK|PAWN)], _ = loadSVG("assets/images/bp.svg", SquareSize, SquareSize)

	return pieces, nil
}

func loadSVG(path string, width, height int) (*ebiten.Image, error) {
	// Read the SVG file
	svgBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse the SVG
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svgBytes))
	if err != nil {
		return nil, err
	}

	// Set the SVG dimensions
	icon.SetTarget(0, 0, float64(width), float64(height))

	// Create a raster image to render to
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	// Rasterize the SVG
	scanner := rasterx.NewScannerGV(width, height, rgba, rgba.Bounds())
	raster := rasterx.NewDasher(width, height, scanner)
	icon.Draw(raster, 1.0)

	// Convert to Ebiten image
	ebitenImg := ebiten.NewImageFromImage(rgba)

	return ebitenImg, nil
}

func (g *Game) MouseCoordsToBoardCoords(x, y int) (int, int) {
	margin := (WindowWidth - BoardSize) / 2
	if x < margin || x > margin+BoardSize || y < margin || y > margin+BoardSize {
		return -1, -1
	}

	// convert the x and y to board coordinates
	file := (x - margin) / SquareSize
	rank := 7 - ((y - margin) / SquareSize)
	return rank, file
}

func GenerateBackground() (*ebiten.Image, error) {
	// make a background image, and light and dark squares
	background := ebiten.NewImage(BoardSize, BoardSize)
	darkSquare := ebiten.NewImage(SquareSize, SquareSize)
	darkSquare.Fill(color.RGBA{181, 136, 99, 255})
	lightSquare := ebiten.NewImage(SquareSize, SquareSize)
	lightSquare.Fill(color.RGBA{240, 217, 181, 255})
	// draw the squares to the background
	for rank := range 8 {
		for file := range 8 {
			var square *ebiten.Image
			if (rank+file)%2 == 0 {
				square = lightSquare
			} else {
				square = darkSquare
			}
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(file*SquareSize), float64(rank*SquareSize))
			background.DrawImage(square, opts)
		}
	}
	return background, nil
}
