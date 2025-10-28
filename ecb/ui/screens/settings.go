package screens

import (
	"context"
	"strconv"
	"sync"

	"deeply/ecb/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SettingsScreen struct {
	cfg *config.Service

	content fyne.CanvasObject
	form    *widget.Form
	once    sync.Once
}

func NewSettingsScreen(cfg *config.Service) *SettingsScreen {
	return &SettingsScreen{cfg: cfg}
}

func (s *SettingsScreen) Route() string {
	return "/pengaturan"
}

func (s *SettingsScreen) Title() string {
	return "Settings"
}

func (s *SettingsScreen) Icon() fyne.Resource {
	return theme.SettingsIcon()
}

func (s *SettingsScreen) Content(ctx context.Context) fyne.CanvasObject {
	s.ensureContent(ctx)
	return s.content
}

func (s *SettingsScreen) OnShow(ctx context.Context) {
	s.ensureContent(ctx)
	settings, _ := s.cfg.Load(ctx)
	ipLocal := settings.LocalIP
	ipSimo := settings.SimoIP
	useWlan := settings.UseWLAN
	lineName := settings.LineName
	tack := strconv.Itoa(settings.Tacktime)
	s.form.Items[0].Widget.(*widget.Entry).SetText(ipLocal)
	s.form.Items[1].Widget.(*widget.Entry).SetText(ipSimo)
	s.form.Items[2].Widget.(*widget.Check).SetChecked(useWlan)
	s.form.Items[3].Widget.(*widget.Entry).SetText(lineName)
	s.form.Items[4].Widget.(*widget.Entry).SetText(tack)
}

func (s *SettingsScreen) ensureContent(ctx context.Context) {
	s.once.Do(func() {
		ipLocal := widget.NewEntry()
		ipLocal.SetPlaceHolder("192.168.0.10")
		ipSimo := widget.NewEntry()
		ipSimo.SetPlaceHolder("10.30.1.5")
		useWlan := widget.NewCheck("Use WLAN", nil)
		lineName := widget.NewEntry()
		lineName.SetPlaceHolder("Line A")
		tacktime := widget.NewEntry()
		tacktime.SetPlaceHolder("0")

		s.form = &widget.Form{
			OnSubmit: func() {
				value, _ := strconv.Atoi(tacktime.Text)
				_ = s.cfg.Save(ctx, config.Settings{
					LocalIP:  ipLocal.Text,
					SimoIP:   ipSimo.Text,
					UseWLAN:  useWlan.Checked,
					LineName: lineName.Text,
					Tacktime: value,
				})
			},
			OnCancel: func() {
				s.OnShow(ctx)
			},
		}
		s.form.Append("Local IP", ipLocal)
		s.form.Append("SIMO3 IP", ipSimo)
		s.form.Append("Wireless", useWlan)
		s.form.Append("Line Name", lineName)
		s.form.Append("Tacktime (s)", tacktime)

		card := widget.NewCard("ECB Settings", "Mirror of Ecbconfig entries", s.form)
		s.content = container.NewVBox(card)
	})
}
