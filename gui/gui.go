package gui

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	// This is because of my install path, should be able to delete and save the file to override these
	// to your install paths ~Dan Mayor
	"github.com/DanMayor/http-monitoring/alerts"
	"github.com/DanMayor/http-monitoring/monitoring"
	"github.com/gizak/termui"
)

/*
   Arrays of strings used to display the contents of the previous alerts messages and the current lines parsed
   from the file log
*/
var (
	Alerts   = []string{}
	LogLines = []string{}
)

// UI Layouts
var (
	LastAlert      = termui.NewPar("")
	Statistics     = termui.NewList()
	Logs           = termui.NewList()
	AlertsHistoric = termui.NewList()
	Info           = termui.NewPar("Press 'Q' to quit.")
)

// String Formats
var (
	TotalHTTPRequests  = "Total HTTP requests : %d"
	TotalSizeEmitted   = "Total Size emitted in bytes : %d"
	MostViewedSections = "Most viewed sections : %s"
	HTTPStatus         = "HTTP Status : %s"
)

// Pair is a structure used for the ordering functions of a map[string]int
type Pair struct {
	Key   string
	Value int
}

// orderMapByValue will return a new []Pair slice sorted by Value
func orderMapByValue(m map[string]int) []Pair {
	/*
	   When ordering the map, we first put each of its key/value in a Pair structure.
	   Each of these Pair structures will be put in an array, then sorted by Value
	*/
	var array []Pair

	for k, v := range m {
		array = append(array, Pair{k, v})
	}

	sort.Slice(array, func(i, j int) bool {
		return array[i].Value > array[j].Value
	})

	return array
}

// getMostViewedSections returns a string displaying the most viewed section of the site, ordered from the highest to the lower
func getMostViewedSections(mapSection map[string]int) string {
	var buffer string
	orderedArray := orderMapByValue(mapSection)

	for i, kv := range orderedArray {
		if i > 2 {
			break
		}

		// This will append [key : value] to the buffer
		buffer = fmt.Sprintf("%s[%s : %d] ", buffer, kv.Key, kv.Value)
	}

	// This trims white space (the trailing space) and returns the buffer
	// using the MostViewedSections format.
	return fmt.Sprintf(MostViewedSections, strings.TrimRight(buffer, " "))
}

// Get a string displaying the HTTP status of the site, ordered from the highest to the lower
func getHTTPstatus(mapHTTPstatus map[string]int) string {
	var buffer string
	orderedArray := orderMapByValue(mapHTTPstatus)

	for i, kv := range orderedArray {
		if i > 4 {
			break
		}

		// This will append [key : value] to the buffer
		buffer = fmt.Sprintf("%s[%s : %d] ", buffer, kv.Key, kv.Value)
	}

	// This trims white space (the trailing space) and returns the buffer
	// using the HTTPStatus format.
	return fmt.Sprintf(HTTPStatus, strings.TrimRight(buffer, " "))
}

func displayUI() {
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

// Every time we receive something from a channel, we update the GUI
func uiLoop() {
	for {
		select {

		// Update the statistics
		case s := <-monitoring.MonitoringData_chan:
			stats := []string{
				fmt.Sprintf(TotalHTTPRequests, s.TotalRequests),
				fmt.Sprintf(TotalSizeEmitted, s.TotalSize),
				getHTTPstatus(s.MapHTTPstatus),
				getMostViewedSections(s.MapURLsection),
			}
			Statistics.Items = stats

			// Update the "LastAlert", and append it to "AlertsHistoric"
		case s := <-alerts.HighTrafficAlert_chan:
			LastAlert.Text = s
			if len(Alerts) > 4 {
				Alerts = Alerts[1:]
			}
			Alerts = append(Alerts, s)
			AlertsHistoric.Items = Alerts

			// Update the "LineLogs" with the latest line received then parsed from the log file
		case s := <-monitoring.LineLog_chan:
			if len(LogLines) > 9 {
				LogLines = LogLines[1:]
			}
			LogLines = append(LogLines, s)
			Logs.Items = LogLines
		}

		// Refresh the UI with the updated data
		termui.Render(termui.Body)
	}
}

func cleanExit() {
	termui.Close()

	close(alerts.HighTrafficAlert_chan)
	close(monitoring.LineLog_chan)
	close(monitoring.MonitoringData_chan)

	os.Exit(0)
}

func Gui() {
	err := termui.Init()
	if err != nil {
		// I get "The parameter is incorrect" error here, don't know the termui package :( ~Dan Mayor
		log.Panic(err)
	}

	displayUI()
	go uiLoop()

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

	termui.Loop()

	cleanExit()
}
