package main

import (
    "log"
    "time"
    "github.com/hpcloud/tail"
)

// UI Channels
var DisplayLogLineChan = make(chan string, 1)
var DisplayMetricsChan = make(chan Metrics, 1)
var DisplayTrafficAlertChan = make(chan string, 1)


//Read a log file, parse each line then send them to the LogLine channel
func StartLogFileFollower(commonLogFile *tail.Tail) {
    for line := range commonLogFile.Lines {
        println(line.Text)
        LogLineChan <- line.Text
    }
}


func StartMonitoring(metrics Metrics) {

    // Tickers triggered each X seconds
    tickerMetrics := time.NewTicker(time.Second * 2).C
    tickerAlertRequests := time.NewTicker(time.Second * 4).C

    for {
        select {

            // Receive new line from the log file
        case line := <-LogLineChan:

            // Send this raw string line to the UI
            DisplayLogLineChan <- line

            // Convert the line to a LogLine struct
            logLine, err := ParseLine(line)
            if err != nil {
                log.Println(err.Error())
            }

            // Update the metrics with the new data from a logLine struct
            metrics.Update(logLine)

            // Send monitoring stats to UI on each tick
        case <-tickerMetrics:
            DisplayMetricsChan <- metrics

            // Check the number of HTTP requests on each X tick
        case <-tickerAlertRequests:
            alertMessage := CheckRequestsThreshold(metrics)
            if alertMessage != "" {
                DisplayTrafficAlertChan <- alertMessage
            }

        }
    }

}
