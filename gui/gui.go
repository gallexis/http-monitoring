package gui

import (
    "github.com/gizak/termui"
    "os"
    "http-monitoring/monitoring"
    "bytes"
    "sort"
    "strconv"
    "fmt"
    "http-monitoring/alert"
)

var (
    Alerts   = []string{}
    LogLines = []string{}
)

var (
    LastAlert      = termui.NewPar("")
    Statistics     = termui.NewList()
    Logs           = termui.NewList()
    AlertsHistoric = termui.NewList()
    Info           = termui.NewPar("Press 'Q' to quit.")
)

var(
    TotalHTTPRequests  = "Total HTTP requests : %d"
    TotalSizeEmitted   = "Total Size emitted in bytes : %d"
    MostViewedSections = "Most viewed sections: %s"
)


func getMostViewedSections(mapSection map[string]int) string {
    type sectionByHits struct {
        Section string
        Hits    int
    }
    var buffer bytes.Buffer
    var ss []sectionByHits

    for k, v := range mapSection {
        ss = append(ss, sectionByHits{k, v})
    }

    sort.Slice(ss, func(i, j int) bool {
        return ss[i].Hits > ss[j].Hits
    })

    for i, kv := range ss {
        if i > 2{
            break
        }
        buffer.WriteString("[")
        buffer.WriteString(kv.Section)
        buffer.WriteString(" : ")
        buffer.WriteString(strconv.Itoa(kv.Hits))
        buffer.WriteString("] ")
    }

    return fmt.Sprintf(MostViewedSections, buffer.String())
}

func setUI() {
    LastAlert.Height = 3
    LastAlert.BorderLabel = "Last Alert"

    Statistics.Height = 6
    Statistics.BorderLabel = "Statistics"

    Logs.Items = Alerts
    Logs.Height = 12
    Logs.BorderLabel = "Logs"
    //Logs.BorderFg = termui.ColorYellow

    AlertsHistoric.Height = 7
    AlertsHistoric.Y = Logs.Y + Logs.Height
    AlertsHistoric.BorderLabel = "Alerts Historic"

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

}

func uiLoop(){
    for {
        select {
        case s := <- monitoring.MonitoringDataChan:
            stats := []string{
                fmt.Sprintf(TotalHTTPRequests, s.TotalRequests),
                fmt.Sprintf(TotalSizeEmitted, s.TotalSize),
                getMostViewedSections(s.MapURLsection),
            }
            Statistics.Items = stats

        case s := <- alert.HighTrafficAlert_chan:
            LastAlert.Text = s
            if len(Alerts) > 4 {
                Alerts = Alerts[1:]
            }
            Alerts = append(Alerts, s)
            AlertsHistoric.Items = Alerts

        case s := <- monitoring.LineLog_chan:
            if len(LogLines) > 9 {
                LogLines = LogLines[1:]
            }
            LogLines = append(LogLines, s)
            Logs.Items = LogLines
        }

        // Every time something changes -> update the GUI
        termui.Render(termui.Body)
    }
}

func Gui() {
    err := termui.Init()
    if err != nil {
        panic(err)
    }

    setUI()
    go uiLoop()

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
    termui.Close()
    os.Exit(0)
}
