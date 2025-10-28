package screens

import (
	"context"
	"sync"
	"time"

	"deeply/ecb/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MaintenanceScreen struct {
	appState *state.AppState

	content fyne.CanvasObject
	mode    *widget.RadioGroup
	once    sync.Once
}

func NewMaintenanceScreen(appState *state.AppState) *MaintenanceScreen {
	return &MaintenanceScreen{
		appState: appState,
	}
}

func (s *MaintenanceScreen) Route() string {
	return "/maintenance"
}

func (s *MaintenanceScreen) Title() string {
	return "Maintenance"
}

func (s *MaintenanceScreen) Icon() fyne.Resource {
	return theme.SettingsIcon()
}

func (s *MaintenanceScreen) Content(ctx context.Context) fyne.CanvasObject {
	s.ensureContent(ctx)
	return s.content
}

func (s *MaintenanceScreen) OnShow(ctx context.Context) {
	s.ensureContent(ctx)
	s.appState.WithRead(ctx, func(state *state.AppState) {
		s.mode.SetSelected(state.Mode)
	})
}

func (s *MaintenanceScreen) ensureContent(ctx context.Context) {
	s.once.Do(func() {
		s.mode = widget.NewRadioGroup([]string{"LIVE", "simulateHW", "simulateDB", "simulateAll"}, func(value string) {
			_ = s.appState.WithLock(ctx, func(state *state.AppState) error {
				state.Mode = value
				state.LastModeChange = time.Now()
				return nil
			})
		})
		s.mode.Required = true
		card := widget.NewCard("ECB Mode", "Switch between live hardware and simulation", s.mode)
		s.content = container.NewVBox(card)
	})
}
