package screens

import (
	"context"
	"fmt"
	"sync"

	"deeply/ecb/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AboutScreen struct {
	cfg     *config.Service
	content fyne.CanvasObject
	once    sync.Once
}

func NewAboutScreen(cfg *config.Service) *AboutScreen {
	return &AboutScreen{cfg: cfg}
}

func (s *AboutScreen) Route() string {
	return "/about"
}

func (s *AboutScreen) Title() string {
	return "About"
}

func (s *AboutScreen) Icon() fyne.Resource {
	return theme.InfoIcon()
}

func (s *AboutScreen) Content(ctx context.Context) fyne.CanvasObject {
	s.ensureContent(ctx)
	return s.content
}

func (s *AboutScreen) OnShow(ctx context.Context) {
	s.ensureContent(ctx)
}

func (s *AboutScreen) ensureContent(ctx context.Context) {
	s.once.Do(func() {
		settings, _ := s.cfg.Load(ctx)
		info := widget.NewRichTextFromMarkdown(fmt.Sprintf(
			"**ECB Test Station**\n\nVersion: `%s`\nLine: `%s`\nLocation: `%s`\n\nProgrammed by `rudy.gunawan@polytron.co.id`",
			"201808221", settings.LineName, settings.Location,
		))
		info.Wrapping = fyne.TextWrapWord
		s.content = container.NewBorder(nil, nil, nil, nil, info)
	})
}
