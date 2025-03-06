package gui

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
	// if g.Board.WhiteToMove {
	// 	turn = "White to move"
	// } else {
	// 	turn = "Black to move"
	// }
	legalMoves := ""
	for _, move := range g.Board.LegalMoves {
		legalMoves = legalMoves + move.String() + " "
	}
	ebitenutil.DebugPrint(screen, legalMoves)
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
			index := rank*8 + file

			// check if the square is selected
			squareSelected := g.Selected == index

			// highlight selected square or the squares from the prior move
			if squareSelected {
				square := ebiten.NewImage(SquareSize, SquareSize)
				squareColor := color.RGBA{255, 0, 0, 100}
				square.Fill(squareColor)
				squareOpts := &ebiten.DrawImageOptions{}
				squareOpts.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(square, squareOpts)
			}
			if g.PrevMove != 0 {
				source := g.PrevMove.Source()
				target := g.PrevMove.Target()
				if index == source || index == target {
					square := ebiten.NewImage(SquareSize, SquareSize)
					squareColor := color.RGBA{0, 255, 0, 100}
					square.Fill(squareColor)
					squareOpts := &ebiten.DrawImageOptions{}
					squareOpts.GeoM.Translate(float64(x), float64(y))
					screen.DrawImage(square, squareOpts)
				}
			}

			// get the piece on the square, skip if there is no piece
			piece := g.Board.GetPieceAtIndex(index)
			if piece.IsNone() {
				continue
			}

			// don't draw the piece in the square if it's being dragged
			if squareSelected && g.Dragging {
				continue
			}

			// draw the piece on the square
			pieceImg := g.PieceImages[piece.FenChar()]
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
		pieceImage := g.PieceImages[g.Board.GetPieceAtIndex(g.Selected).FenChar()]
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(pieceImage, opts)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WindowWidth, WindowHeight
}

// loadSVG("assets/images/"+"bp"+".svg", SquareSize, SquareSize)

func loadPieceImages() (map[string]*ebiten.Image, error) {
	pieceImages := make(map[string]*ebiten.Image)
	fenChars := []string{"p", "n", "b", "r", "q", "k", "P", "N", "B", "R", "Q", "K"}
	fileNames := []string{"bp", "bn", "bb", "br", "bq", "bk", "wp", "wn", "wb", "wr", "wq", "wk"}
	for i, fenChar := range fenChars {
		img, err := loadSVG("assets/images/"+fileNames[i]+".svg", SquareSize, SquareSize)
		if err != nil {
			return nil, err
		}
		pieceImages[fenChar] = img
	}
	return pieceImages, nil
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

// translates where the mouse is to board rank and file, returns -1, -1 if off the board
func (g *Game) mouseCoordsToBoardCoords(x, y int) (int, int) {
	margin := (WindowWidth - BoardSize) / 2
	if x < margin || x > margin+BoardSize || y < margin || y > margin+BoardSize {
		return -1, -1
	}

	// convert the x and y to board coordinates
	file := (x - margin) / SquareSize
	rank := 7 - ((y - margin) / SquareSize)
	return rank, file
}

// pre-generates the background image so we only do it once
func generateBackground() (*ebiten.Image, error) {
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
