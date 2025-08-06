package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type App struct {
	CurrentState GameMode
	MainMenu     *MainMenuState
	Game         *Game
	InitialFEN   string
}

func NewApp(startMode GameMode, fenString string) *App {
	app := &App{
		CurrentState: startMode,
		MainMenu:     NewMainMenuState(),
		InitialFEN:   fenString,
	}
	
	if startMode != MainMenu {
		app.initializeGame()
	}
	
	return app
}

func (a *App) initializeGame() {
	if a.Game == nil {
		a.Game = NewGame(a.InitialFEN)
	}
}

func (a *App) Update() error {
	switch a.CurrentState {
	case MainMenu:
		newMode, err := a.MainMenu.Update()
		if err != nil {
			return err
		}
		if newMode != MainMenu {
			a.CurrentState = newMode
			a.initializeGame()
		}
	case DebugMode, HumanVsAI, AIvsAI:
		if a.Game != nil {
			err := a.Game.Update()
			if err == ErrReturnToMenu {
				a.CurrentState = MainMenu
				a.Game = nil // Reset game state when returning to menu
				return nil
			}
			return err
		}
	}
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	switch a.CurrentState {
	case MainMenu:
		a.MainMenu.Draw(screen)
	case DebugMode, HumanVsAI, AIvsAI:
		if a.Game != nil {
			a.Game.Draw(screen)
		}
	}
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 1000
}