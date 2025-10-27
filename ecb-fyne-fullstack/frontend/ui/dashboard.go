package ui

import (
	"fmt"
	"frontend/api"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Dashboard struct {
	statusBinding binding.String
	valueBinding  binding.String
	logs          []string
	logList       *widget.List
}

func NewDashboard() *Dashboard {
	status := binding.NewString()
	_ = status.Set("⏳ Menunggu data...")

	value := binding.NewString()
	_ = value.Set("0.00")

	return &Dashboard{
		statusBinding: status,
		valueBinding:  value,
		logs:          []string{},
	}
}

func (d *Dashboard) Build() fyne.CanvasObject {
	// Status Card
	statusLabel := widget.NewLabelWithData(d.statusBinding)
	statusLabel.TextStyle = fyne.TextStyle{Bold: true}

	valueLabel := widget.NewLabelWithData(d.valueBinding)
	valueLabel.TextStyle = fyne.TextStyle{Bold: true}

	statusCard := widget.NewCard("📊 Status Sensor", "", container.NewVBox(
		widget.NewLabel("Status:"),
		statusLabel,
		widget.NewSeparator(),
		widget.NewLabel("Value:"),
		valueLabel,
	))

	// Control Buttons
	refreshBtn := widget.NewButton("🔄 Refresh", d.refreshData)
	clearBtn := widget.NewButton("🗑️ Clear Logs", d.clearLogs)
	controls := container.NewHBox(refreshBtn, clearBtn)

	// Log List
	d.logList = widget.NewList(
		func() int { return len(d.logs) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(d.logs[id])
		},
	)

	logCard := widget.NewCard("📝 Log Activity", "", d.logList)

	// Layout
	content := container.NewBorder(
		container.NewVBox(statusCard, controls),
		nil, nil, nil,
		logCard,
	)

	return content
}

func (d *Dashboard) refreshData() {
	sensor, err := api.GetLatestSensor()
	if err != nil {
		d.addLog("❌ Error: " + err.Error())
		_ = d.statusBinding.Set("⚠️ Error mengambil data")
		return
	}

	_ = d.statusBinding.Set(sensor.Status)
	_ = d.valueBinding.Set(fmt.Sprintf("%.2f", sensor.Value))

	timestamp := sensor.Timestamp.Format("15:04:05")
	d.addLog(fmt.Sprintf("[%s] ✅ Status: %s | Value: %.2f",
		timestamp, sensor.Status, sensor.Value))
}

func (d *Dashboard) addLog(message string) {
	d.logs = append([]string{message}, d.logs...)
	if len(d.logs) > 100 {
		d.logs = d.logs[:100]
	}
	d.logList.Refresh()
}

func (d *Dashboard) clearLogs() {
	d.logs = []string{}
	d.logList.Refresh()
}

func (d *Dashboard) StartAutoRefresh() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			d.refreshData()
		}
	}()
}
