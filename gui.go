package main

import (
    "github.com/gizak/termui"
    "os"
)

func Gui(){

    err := termui.Init()
    if err != nil {
        panic(err)
    }
    defer termui.Close()

    Alerts := []string{}
    Log_lines := []string{}

    LastAlert := termui.NewPar("")
    LastAlert.Height = 3
    LastAlert.BorderLabel = "Last Alert"

    Statistics := termui.NewList()
    Statistics.Height = 6
    Statistics.BorderLabel = "Statistics"

    Logs := termui.NewList()
    Logs.Items = Alerts
    Logs.Height = 12
    Logs.BorderLabel = "Logs"
    //Logs.BorderFg = termui.ColorYellow

    AlertsHistoric := termui.NewList()
    AlertsHistoric.Height = 7
    AlertsHistoric.Y = Logs.Y + Logs.Height
    AlertsHistoric.BorderLabel = "Alerts Historic"

    Info := termui.NewPar("Press 'Q' to quit.")
    Info.Height = 3
    Info.BorderLabel = "Info."
    Info.BorderFg = termui.ColorCyan


    termui.Body.AddRows(
        termui.NewRow(
            termui.NewCol(12, 0, LastAlert)),
        termui.NewRow(
            termui.NewCol(12, 0, Statistics)),
        termui.NewRow(
            termui.NewCol(12, 0, Logs)),
        termui.NewRow(
            termui.NewCol(12, 0, AlertsHistoric)),
        termui.NewRow(
            termui.NewCol(12, 0, Info)))

    // calculate layout
    termui.Body.Align()

    go func(){

        for {
            select{
                case s := <- Stats_chan:
                    Statistics.Items = s

                case s := <- HighTrafficAlert_chan:
                    LastAlert.Text = s
                    if len(Alerts) > 4 {
                        Alerts = Alerts[1:]
                    }
                    Alerts = append(Alerts, s)
                    AlertsHistoric.Items = Alerts

                case s := <- lineLog_chan:
                    if len(Log_lines) > 9 {
                        Log_lines = Log_lines[1:]
                    }
                    Log_lines = append(Log_lines, s)
                    Logs.Items = Log_lines
            }

            // Every time something changes -> update the GUI
            termui.Render(termui.Body)
        }

    }()

    termui.Handle("/sys/kbd/q", func(termui.Event) {
        termui.StopLoop()
        termui.Clear()
    })

    termui.Handle("/sys/wnd/resize", func(e termui.Event) {
        termui.Body.Width = termui.TermWidth()
        termui.Body.Align()
        termui.Clear()
        termui.Render(termui.Body)
    })

    termui.Loop()
    os.Exit(0)
}