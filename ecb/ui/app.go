package ui

import (
	"context"

	"deeply/ecb/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

type Application struct {
	app       fyne.App
	mainWin   fyne.Window
	router    *Router
	appState  *state.AppState
	rootStack *container.Split
}

func NewApplication(app fyne.App, router *Router, appState *state.AppState) *Application {
	win := app.NewWindow("ECB Station")
	root := container.NewHSplit(router.Content(), container.NewMax())
	root.SetOffset(0.2)
	return &Application{
		app:       app,
		mainWin:   win,
		router:    router,
		appState:  appState,
		rootStack: root,
	}
}

func (a *Application) RegisterScreens(screens ...Screen) {
	for _, screen := range screens {
		a.router.Register(screen)
	}
}

func (a *Application) Run(ctx context.Context, initialRoute string) error {
	a.mainWin.SetContent(a.rootStack)
	a.mainWin.Resize(fyne.NewSize(1280, 720))
	a.app.Settings().SetTheme(theme.DarkTheme())

	if err := a.router.Navigate(ctx, initialRoute); err != nil {
		return err
	}
	a.mainWin.ShowAndRun()
	return nil
}

func (a *Application) Navigate(ctx context.Context, route string) error {
	return a.router.Navigate(ctx, route)
}

func (a *Application) Close(route string) {
	a.router.Close(route)
}
