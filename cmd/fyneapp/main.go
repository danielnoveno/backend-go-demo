package main

import (
	"context"
	"log"

	"deeply/ecb/config"
	"deeply/ecb/hardware"
	"deeply/ecb/navigation"
	"deeply/ecb/state"
	"deeply/ecb/test"
	"deeply/ecb/ui"
	"deeply/ecb/ui/screens"

	"fyne.io/fyne/v2/app"
)

func main() {
	ctx := context.Background()

	navRepo := navigation.NewStaticRepository([]navigation.Record{
		{ID: 1, Title: "Test Station", Icon: "media-play", URL: "/test", Order: 1},
		{ID: 2, Title: "Maintenance", Icon: "settings", URL: "/maintenance", Order: 2},
		{ID: 3, Title: "Settings", Icon: "preferences", URL: "/pengaturan", Order: 3},
		{ID: 4, Title: "About", Icon: "info", URL: "/about", Order: 4},
	})
	navService := navigation.NewService(navRepo)
	if err := navService.Refresh(ctx); err != nil {
		log.Fatalf("init navigation: %v", err)
	}

	cfgRepo := config.NewInMemoryRepository(config.Settings{
		LocalIP:  "192.168.0.10",
		SimoIP:   "10.30.1.5",
		UseWLAN:  false,
		Theme:    "layout1",
		LineType: "refrig-single",
		Location: "Neverland",
		LineName: "Line 1",
		Tacktime: 0,
	})
	cfgService := config.NewService(cfgRepo)

	appState := state.New()

	gpioAdapter := hardware.NewSimulatedAdapter()
	if err := gpioAdapter.Initialize(ctx); err != nil {
		log.Fatalf("init gpio: %v", err)
	}

	testRepo := test.NewInMemoryRepository()
	testService := test.NewService(testRepo, cfgService, gpioAdapter, appState)

	router := ui.NewRouter(navService)
	fyneApp := app.NewWithID("ecb-station")
	application := ui.NewApplication(fyneApp, router, appState)

	application.RegisterScreens(
		screens.NewTestScreen(testService, cfgService, appState),
		screens.NewMaintenanceScreen(appState),
		screens.NewSettingsScreen(cfgService),
		screens.NewAboutScreen(cfgService),
		screens.NewErrorScreen("/neverland", "Not Found", "Halaman tidak ditemukan."),
		screens.NewErrorScreen("/noaccess", "Forbidden", "Akses ditolak."),
	)

	if err := application.Run(ctx, "/test"); err != nil {
		log.Fatalf("run application: %v", err)
	}
}
