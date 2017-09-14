package main

import (
    "log"
    "time"
)

// Channels
var CommonLogToUIChan = make(chan string, 1)
var MetricsToUIChan = make(chan Metrics, 1)
var TrafficAlertToUIChan = make(chan string, 1)


func Monitoring(metrics Metrics) {

    // Tickers triggered each X seconds
    tickerMetricsChan := time.NewTicker(time.Second * 10).C
    tickerAlertRequestsChan := time.NewTicker(time.Second * 120).C

    for {
        select {

            // Receive new line from the log file
        case line := <-CommonLogChan:

            // Send this line to the UI
            CommonLogToUIChan <- line

            // Convert the line to a commonLog struct
            commonLog, err := ParseLine(line)
            if err != nil {
                log.Panic(err)
            }

            // Update the metrics with the new data from a commonLog struct
            metrics.Update(commonLog)

            // Send monitoring stats to UI on each tick
        case <-tickerMetricsChan:
            MetricsToUIChan <- metrics

            // Check the number of HTTP requests on each X tick
        case <-tickerAlertRequestsChan:
            alertMessage := CheckRequestsThreshold(metrics)
            if alertMessage != "" {
                TrafficAlertToUIChan <- alertMessage
            }

        }
    }

}
