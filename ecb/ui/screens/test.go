package screens

import (
	"context"
	"fmt"
	"sync"
	"time"

	"deeply/ecb/config"
	"deeply/ecb/state"
	"deeply/ecb/test"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TestScreen struct {
	svc      *test.Service
	cfg      *config.Service
	appState *state.AppState

	content     fyne.CanvasObject
	once        sync.Once
	statusLabel *widget.Label
	history     *widget.List
	input       *widget.Entry
	numRows     []test.HistoryRow
	cancelPoll  context.CancelFunc
}

func NewTestScreen(svc *test.Service, cfg *config.Service, appState *state.AppState) *TestScreen {
	return &TestScreen{
		svc:      svc,
		cfg:      cfg,
		appState: appState,
	}
}

func (s *TestScreen) Route() string {
	return "/test"
}

func (s *TestScreen) Title() string {
	return "Test Station"
}

func (s *TestScreen) Icon() fyne.Resource {
	return theme.MediaPlayIcon()
}

func (s *TestScreen) Content(ctx context.Context) fyne.CanvasObject {
	s.ensureContent(ctx)
	return s.content
}

func (s *TestScreen) OnShow(ctx context.Context) {
	s.ensureContent(ctx)
	s.reloadHistory(ctx)
	s.startPolling(ctx)
}

func (s *TestScreen) ensureContent(ctx context.Context) {
	s.once.Do(func() {
		s.statusLabel = widget.NewLabel("Status: idle")
		s.input = widget.NewEntry()
		s.input.SetPlaceHolder("Scan barcode...")
		s.input.OnSubmitted = func(value string) {
			record := test.ScanRecord{
				SN:     value,
				Line:   fmt.Sprintf("%d", s.appState.LineActive),
				Result: "PENDING",
			}
			_ = s.svc.StoreScan(ctx, record)
			s.input.SetText("")
			s.reloadHistory(ctx)
		}

		s.history = widget.NewList(
			func() int { return len(s.numRows) },
			func() fyne.CanvasObject {
				return widget.NewLabel("item")
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				if i >= len(s.numRows) {
					return
				}
				row := s.numRows[i]
				o.(*widget.Label).SetText(fmt.Sprintf("%s | %s | %s", row.Created.Format(time.RFC822), row.Line, row.SN))
			},
		)

		toolbar := widget.NewToolbar(
			widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
				_ = s.svc.StartTest(ctx, s.appState.LineActive)
			}),
			widget.NewToolbarAction(theme.MediaStopIcon(), func() {
				_ = s.svc.ResetRig(ctx, s.appState.LineActive)
			}),
			widget.NewToolbarSpacer(),
			widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
				s.reloadHistory(ctx)
			}),
		)

		inputBox := container.NewVBox(widget.NewLabel("Scan Input"), s.input)
		layout := container.NewBorder(
			container.NewVBox(toolbar, s.statusLabel),
			inputBox,
			nil,
			nil,
			s.history,
		)

		s.content = layout
	})
}

func (s *TestScreen) reloadHistory(ctx context.Context) {
	settings, _ := s.cfg.Load(ctx)
	records, err := s.svc.LoadHistory(ctx, settings.LineName, 20)
	if err != nil {
		return
	}
	s.numRows = records
	s.history.Refresh()
}

func (s *TestScreen) startPolling(ctx context.Context) {
	if s.cancelPoll != nil {
		s.cancelPoll()
	}
	ctx, cancel := context.WithCancel(ctx)
	s.cancelPoll = cancel
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				status, err := s.svc.PollStatus(context.Background())
				if err != nil {
					continue
				}
				statusText := fmt.Sprintf("Status: line %d | pass=%t fail=%t testing=%t", status.Line, status.Pass, status.Fail, status.UnderTest)
				label := s.statusLabel
				if label != nil {
					fyne.CurrentApp().Driver().RunOnMain(func() {
						label.SetText(statusText)
					})
				}
			}
		}
	}()
}
