package gui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type GameMode int

const (
	MainMenu GameMode = iota
	DebugMode
	HumanVsAI
	AIvsAI
)

type MenuOption struct {
	Text string
	Mode GameMode
	Rect image.Rectangle
}

type MainMenuState struct {
	Background *ebiten.Image
	Options    []MenuOption
	Font       font.Face
}

func NewMainMenuState() *MainMenuState {
	background, err := generateBackground()
	if err != nil {
		panic(err)
	}

	options := []MenuOption{
		{Text: "Debug Mode (Manual Play)", Mode: DebugMode},
		{Text: "Human vs AI", Mode: HumanVsAI},
		{Text: "AI vs AI", Mode: AIvsAI},
	}

	return &MainMenuState{
		Background: background,
		Options:    options,
		Font:       basicfont.Face7x13,
	}
}

func (m *MainMenuState) Update() (GameMode, error) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for _, option := range m.Options {
			if x >= option.Rect.Min.X && x <= option.Rect.Max.X &&
				y >= option.Rect.Min.Y && y <= option.Rect.Max.Y {
				return option.Mode, nil
			}
		}
	}
	return MainMenu, nil
}

func (m *MainMenuState) Draw(screen *ebiten.Image) {
	// Center the background like in gameplay  
	margin := (WindowWidth - BoardSize) / 2
	bgOpts := &ebiten.DrawImageOptions{}
	bgOpts.GeoM.Translate(float64(margin), float64(margin))
	screen.DrawImage(m.Background, bgOpts)
	
	// Draw title with background
	titleText := "Go Chess"
	titleBounds := text.BoundString(m.Font, titleText)
	titleX := (screen.Bounds().Dx() - titleBounds.Dx()) / 2
	titleY := 150
	titlePadding := 20
	
	// Title background
	vector.DrawFilledRect(screen, 
		float32(titleX-titlePadding), 
		float32(titleY-titleBounds.Dy()-titlePadding/2), 
		float32(titleBounds.Dx()+titlePadding*2), 
		float32(titleBounds.Dy()+titlePadding), 
		color.RGBA{0, 0, 0, 180}, false)
	
	text.Draw(screen, titleText, m.Font, titleX, titleY, color.White)
	
	// Draw help text with background
	helpText := "Press ESC in any mode to return to this menu"
	helpBounds := text.BoundString(m.Font, helpText)
	helpX := (screen.Bounds().Dx() - helpBounds.Dx()) / 2
	helpY := 250
	helpPadding := 10
	
	// Help text background
	vector.DrawFilledRect(screen, 
		float32(helpX-helpPadding), 
		float32(helpY-helpBounds.Dy()-helpPadding/2), 
		float32(helpBounds.Dx()+helpPadding*2), 
		float32(helpBounds.Dy()+helpPadding), 
		color.RGBA{0, 0, 0, 120}, false)
	
	text.Draw(screen, helpText, m.Font, helpX, helpY, color.RGBA{220, 220, 220, 255})
	
	// Draw menu options
	startY := 320
	spacing := 80
	buttonWidth := 300
	buttonHeight := 50
	
	for i, option := range m.Options {
		y := startY + i*spacing
		x := (screen.Bounds().Dx() - buttonWidth) / 2
		
		// Update option rect for click detection
		m.Options[i].Rect = image.Rect(x, y-buttonHeight/2, x+buttonWidth, y+buttonHeight/2)
		
		// Draw button background
		vector.DrawFilledRect(screen, float32(x), float32(y-buttonHeight/2), float32(buttonWidth), float32(buttonHeight), color.RGBA{40, 40, 40, 255}, false)
		vector.StrokeRect(screen, float32(x), float32(y-buttonHeight/2), float32(buttonWidth), float32(buttonHeight), 2, color.RGBA{100, 100, 100, 255}, false)
		
		// Draw button text
		textBounds := text.BoundString(m.Font, option.Text)
		textX := x + (buttonWidth-textBounds.Dx())/2
		textY := y + textBounds.Dy()/2
		text.Draw(screen, option.Text, m.Font, textX, textY, color.White)
	}
}