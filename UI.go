package main

import (
	"github.com/gizak/termui"
	"log"
	"os"
)

/*
   Arrays of strings used to display the contents of the alert messages and the lines received
   from the log file
*/
var (
	Alerts   = []string{}
	LogLines = []string{}
)

// UI Layouts
var (
	LastAlert      = termui.NewPar("")
	MetricsData    = termui.NewList()
	Log            = termui.NewList()
	AlertsHistoric = termui.NewList()
	Info           = termui.NewPar("Press 'Q' to quit.")
)

func display() {
	LastAlert.Height = 3
	LastAlert.BorderLabel = "Last Alert"

	MetricsData.Height = 6
	MetricsData.BorderLabel = "Metrics"

	Log.Items = Alerts
	Log.Height = 12
	Log.BorderLabel = "Log"

	AlertsHistoric.Height = 7
	AlertsHistoric.Y = Log.Y + Log.Height
	AlertsHistoric.BorderLabel = "Alerts Historic"

	Info.Height = 3
	Info.BorderLabel = "Info."
	Info.BorderFg = termui.ColorCyan

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, LastAlert)),
		termui.NewRow(
			termui.NewCol(12, 0, MetricsData)),
		termui.NewRow(
			termui.NewCol(12, 0, Log)),
		termui.NewRow(
			termui.NewCol(12, 0, AlertsHistoric)),
		termui.NewRow(
			termui.NewCol(12, 0, Info)))

	// calculate layout
	termui.Body.Align()
}

// Every time we receive something from a channel, we update the UI
func EventLoop() {
	for {
		select {

		case metrics := <-MetricsToUI_chan:
			MetricsData.Items = metrics.Display()

		case alertMessage := <-TrafficAlertToUI_chan:
			// Update the "LastAlert" ...
			LastAlert.Text = alertMessage

			// ... and append it to "AlertsHistoric"
			if len(Alerts) > 4 {
				Alerts = Alerts[1:]
			}
			Alerts = append(Alerts, alertMessage)
			AlertsHistoric.Items = Alerts

		case line := <-CommonLogToUI_chan:
			// Append LogLines with the latest line received from the log file
			if len(LogLines) > 9 {
				LogLines = LogLines[1:]
			}
			LogLines = append(LogLines, line)
			Log.Items = LogLines
		}

		// Refresh the UI with the updated data
		termui.Render(termui.Body)
	}
}

func RunUI() {
	err := termui.Init()
	if err != nil {
		log.Panic(err)
	}

	display()
	go EventLoop()

	// Exit when press 'Q'
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
		termui.Clear()
	})

	// Resize window
	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Clear()
		termui.Render(termui.Body)
	})

	// Loop until 'Q' is pressed
	termui.Loop()

	// Exit
	termui.Close()
	os.Exit(0)
}
