package main

import (
	"frontend/ui"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	if err := os.Setenv("FYNE_THEME", "light"); err != nil {
		log.Println("Warning: unable to set theme:", err)
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("ECB Test Monitor")
	myWindow.Resize(fyne.NewSize(1024, 768))

	dashboard := ui.NewDashboard()
	myWindow.SetContent(dashboard.Build())

	dashboard.StartAutoRefresh()
	myWindow.ShowAndRun()
}
