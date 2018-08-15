package main

import (
    "github.com/gizak/termui"
    "log"
    "os"
)

type UiLayouts struct {
    LastAlert     *termui.Par
    MetricsData   *termui.List
    Log           *termui.List
    AlertsHistory *termui.List
    Info          *termui.Par
}

func NewUiLayout() UiLayouts {
    return UiLayouts{
        LastAlert:     termui.NewPar(""),
        MetricsData:   termui.NewList(),
        Log:           termui.NewList(),
        AlertsHistory: termui.NewList(),
        Info:          termui.NewPar("Press 'Q' to quit."),
    }
}

func (ui UiLayouts) display() {
    ui.LastAlert.Height = 3
    ui.LastAlert.BorderLabel = "Last Alert"

    ui.MetricsData.Height = 6
    ui.MetricsData.BorderLabel = "Metrics"

    ui.Log.Items = []string{} //ringToStringArray(Alerts)
    ui.Log.Height = 12
    ui.Log.BorderLabel = "Log"

    ui.AlertsHistory.Height = 7
    ui.AlertsHistory.Y = ui.Log.Y + ui.Log.Height
    ui.AlertsHistory.BorderLabel = "Alerts Historic"

    ui.Info.Height = 3
    ui.Info.BorderLabel = "Info."
    ui.Info.BorderFg = termui.ColorCyan

    termui.Body.AddRows(
        termui.NewRow(
            termui.NewCol(12, 0, ui.LastAlert)),
        termui.NewRow(
            termui.NewCol(12, 0, ui.MetricsData)),
        termui.NewRow(
            termui.NewCol(12, 0, ui.Log)),
        termui.NewRow(
            termui.NewCol(12, 0, ui.AlertsHistory)),
        termui.NewRow(
            termui.NewCol(12, 0, ui.Info)))

    // calculate layout
    termui.Body.Align()
}

// Every time we receive something from a channel, we update the UI
func (ui UiLayouts) EventLoop() {
    for {
        select {

        case metrics := <-DisplayMetricsChan:
            ui.MetricsData.Items = metrics.Display()

        case metrics := <-DisplayTrafficAlertChan:
            ui.LastAlert.Text = metrics.LastAlertMessage                       // Update the "LastAlert" ...
            ui.AlertsHistory.Items = metrics.ringToStringArray(metrics.Alerts) // ... and append it to "AlertsHistory"

        case metrics := <-DisplayLogLineChan:
            // Append LogLines with the latest line received from the log file
            ui.Log.Items = metrics.ringToStringArray(metrics.LogLines)
        }

        // Refresh the UI with the updated data
        termui.Render(termui.Body)
    }
}

func RunUI() {
    err := termui.Init()
    if err != nil {
        log.Panic(err.Error())
    }

    ui := NewUiLayout()

    ui.display()
    go ui.EventLoop()

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
