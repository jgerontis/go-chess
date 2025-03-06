package gui

import (
	"bytes"
	"image"
	"image/color"
	"os"
	"strconv"

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
	// get the legal targets for the selected piece
	legalTargetsInts := []int{}
	if g.Selected != -1 {
		for _, move := range g.Board.LegalMoves {
			if move.Source() == g.Selected {
				legalTargetsInts = append(legalTargetsInts, move.Target())
			}
		}
	}

	// get the margin for math
	margin := (WindowWidth - BoardSize) / 2

	// draw the background
	bgOpts := &ebiten.DrawImageOptions{}
	bgOpts.GeoM.Translate(float64(margin), float64(margin))
	screen.DrawImage(g.Background, bgOpts)

	// we're going to make a bunch of tiles, then draw them on the screen

	for index := range 64 {
		tile := ebiten.NewImage(SquareSize, SquareSize)
		tileOpts := &ebiten.DrawImageOptions{}

		// red highlight selected piece
		if g.Selected == index {
			square, opts := makeSquare(color.RGBA{255, 0, 0, 200})
			tile.DrawImage(square, opts)
		}
		// yellow highlight previous move
		if g.PrevMove != 0 {
			source := g.PrevMove.Source()
			target := g.PrevMove.Target()
			if index == source || index == target {
				square, opts := makeSquare(color.RGBA{100, 100, 0, 100})
				tile.DrawImage(square, opts)
			}
		}
		// blue highlight legal targets
		if g.Selected != -1 && len(legalTargetsInts) > 0 {
			for _, targetIndex := range legalTargetsInts {
				if index == targetIndex {
					square, opts := makeSquare(color.RGBA{0, 0, 255, 100})
					tile.DrawImage(square, opts)
				}
			}
		}
		// draw the piece on the square if it's not being dragged
		if index != g.Selected || !g.Dragging {
			piece := g.Board.GetPieceAtIndex(index)
			if !piece.IsNone() {
				pieceImg := g.PieceImages[piece.FenChar()]
				pieceOpts := &ebiten.DrawImageOptions{}
				tile.DrawImage(pieceImg, pieceOpts)
			}
		}
		// draw the tile to the Background
		x := (index%8)*SquareSize + margin
		y := (7-index/8)*SquareSize + margin
		tileOpts.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(tile, tileOpts)
	}

	// draw the dragged piece
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

func makeSquare(c color.RGBA) (*ebiten.Image, *ebiten.DrawImageOptions) {
	square := ebiten.NewImage(SquareSize, SquareSize)
	square.Fill(c)
	opts := &ebiten.DrawImageOptions{}
	return square, opts
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
	dark := color.RGBA{181, 136, 99, 255}
	light := color.RGBA{240, 217, 181, 255}
	// draw the squares to the background
	for rank := range 8 {
		for file := range 8 {
			x := file * SquareSize
			y := (7 - rank) * SquareSize
			index := rank*8 + file
			var square *ebiten.Image
			var c color.RGBA
			if (rank+file)%2 == 0 {
				c = light
			} else {
				c = dark
			}
			square, opts := makeSquare(c)
			opts.GeoM.Translate(float64(x), float64(y))
			// add the index to the top left of the square for debugging
			ebitenutil.DebugPrint(square, strconv.Itoa(index))
			background.DrawImage(square, opts)
		}
	}
	return background, nil
}
