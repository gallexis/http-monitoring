package main

import (
    "github.com/gizak/termui"
    "log"
)

type UiLayouts struct {
    LastAlert      *termui.Par
    MonitoringData *termui.List
    Log            *termui.List
    AlertsHistory  *termui.List
    Info           *termui.Par
}

func NewUiLayout() UiLayouts {
    return UiLayouts{
        LastAlert:      termui.NewPar(""),
        MonitoringData: termui.NewList(),
        Log:            termui.NewList(),
        AlertsHistory:  termui.NewList(),
        Info:           termui.NewPar("Press 'Q' to quit."),
    }
}

func (ui UiLayouts) display() {
    ui.LastAlert.Height = 3
    ui.LastAlert.BorderLabel = "Last Alert"

    ui.MonitoringData.Height = 6
    ui.MonitoringData.BorderLabel = "Monitoring Data"

    ui.Log.Items = []string{}
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
            termui.NewCol(12, 0, ui.MonitoringData)),
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

        case metrics := <-DisplayMonitoringDataChan:
            ui.MonitoringData.Items = metrics.Display()

        case metrics := <-DisplayTrafficAlertChan:
            ui.LastAlert.Text = metrics.LastAlertMessage
            ui.AlertsHistory.Items = metrics.ringToStringArray(metrics.Alerts)

        case metrics := <-DisplayLogLineChan:
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
}
