package screens

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ErrorScreen struct {
	route string
	title string
	text  string
	icon  fyne.Resource
}

func NewErrorScreen(route, title, text string) *ErrorScreen {
	return &ErrorScreen{
		route: route,
		title: title,
		text:  text,
		icon:  theme.WarningIcon(),
	}
}

func (s *ErrorScreen) Route() string {
	return s.route
}

func (s *ErrorScreen) Title() string {
	return s.title
}

func (s *ErrorScreen) Icon() fyne.Resource {
	return s.icon
}

func (s *ErrorScreen) Content(ctx context.Context) fyne.CanvasObject {
	return container.NewCenter(widget.NewLabel(s.text))
}

func (s *ErrorScreen) OnShow(ctx context.Context) {}
