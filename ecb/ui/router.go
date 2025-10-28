package ui

import (
	"context"
	"fmt"

	"deeply/ecb/navigation"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// Screen draws a route inside the main window.
type Screen interface {
	Route() string
	Title() string
	Icon() fyne.Resource
	Content(ctx context.Context) fyne.CanvasObject
	OnShow(ctx context.Context)
}

type Router struct {
	navService *navigation.Service
	screens    map[string]Screen
	stack      *container.AppTabs
	openTabs   map[string]*container.TabItem
}

func NewRouter(navService *navigation.Service) *Router {
	return &Router{
		navService: navService,
		screens:    make(map[string]Screen),
		stack:      container.NewAppTabs(),
		openTabs:   make(map[string]*container.TabItem),
	}
}

func (r *Router) Register(screen Screen) {
	r.screens[screen.Route()] = screen
}

func (r *Router) Content() *container.AppTabs {
	return r.stack
}

func (r *Router) Breadcrumb(url string) []navigation.BreadcrumbEntry {
	if r.navService == nil {
		return nil
	}
	return r.navService.Breadcrumb(url)
}

func (r *Router) Navigate(ctx context.Context, route string) error {
	screen, ok := r.screens[route]
	if !ok {
		return fmt.Errorf("screen not registered for route %s", route)
	}
	screen.OnShow(ctx)
	if existing := r.openTabs[route]; existing != nil {
		existing.Text = screen.Title()
		existing.Icon = screen.Icon()
		existing.Content = screen.Content(ctx)
		r.stack.Refresh()
		r.stack.Select(existing)
		return nil
	}
	tab := container.NewTabItemWithIcon(screen.Title(), screen.Icon(), screen.Content(ctx))
	r.stack.Append(tab)
	r.stack.Select(tab)
	r.openTabs[route] = tab
	return nil
}

func (r *Router) Close(route string) {
	if tab, ok := r.openTabs[route]; ok {
		r.stack.Remove(tab)
		delete(r.openTabs, route)
	}
}
