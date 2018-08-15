package main

import (
    "log"
    "time"
    "github.com/hpcloud/tail"
)

// UI Channels
var DisplayLogLineChan = make(chan Metrics, 1)
var DisplayMetricsChan = make(chan Metrics, 1)
var DisplayTrafficAlertChan = make(chan Metrics, 1)


//Read a log file, parse each line then send them to the LogLine channel
func StartLogFileFollower(commonLogFile *tail.Tail) {
    for line := range commonLogFile.Lines {
        LogLineChan <- line.Text
    }
}


func (metrics *Metrics) StartMonitoring() {

    // Tickers triggered each X seconds
    tickerMetrics := time.NewTicker(time.Second * 2).C
    tickerAlertRequests := time.NewTicker(time.Second * 4).C

    for {
        select {

            // Receive new line from the log file
        case line := <-LogLineChan:
            // Send this raw string line to the UI
            metrics.LogLines.Value = line
            metrics.LogLines = metrics.LogLines.Next()

            // Convert the line to a LogLine struct
            logLine, err := ParseLine(line)
            if err != nil {
                log.Println(err.Error())
            }

            // Update the metrics with the new data from a logLine struct
            metrics.Update(logLine)

            DisplayLogLineChan <- *metrics

            // Send monitoring stats to UI on each tick
        case <-tickerMetrics:
            DisplayMetricsChan <- *metrics

            // Check the number of HTTP requests on each X tick
        case <-tickerAlertRequests:
            metrics.CheckRequestsThreshold()

            if metrics.LastAlertMessage != "" {
                metrics.Alerts.Value = metrics.LastAlertMessage
                metrics.Alerts = metrics.Alerts.Next()
                DisplayTrafficAlertChan <- *metrics
            }

        }
    }
}