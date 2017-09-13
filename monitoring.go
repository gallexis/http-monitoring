package main

import (
    "log"
    "time"
)

// Channels
var CommonLogToUI_chan = make(chan string, 1)
var MetricsToUI_chan = make(chan Metrics, 1)
var TrafficAlertToUI_chan = make(chan string, 1)

func Monitoring(metrics *Metrics) {

    // Tickers triggered each X seconds
    metricsToUI_ticker       := time.NewTicker(time.Second * 5).C
    alertRequestsHTTP_ticker := time.NewTicker(time.Second * 10).C

    for {
        select {

            // Receive new line from the log file
        case line := <-CommonLog_chan:

            // Send this line to the UI
            CommonLogToUI_chan <- line

            // Convert the line to a commonLog struct
            parsedCommonLog, err := ParseCommonLog(line)
            if err != nil {
                log.Panic(err)
            }

            // Update the metrics with the new data from a commonLog struct
            metrics.Update(parsedCommonLog)

            // Send monitoring stats to UI on each tick
        case <-metricsToUI_ticker:
            MetricsToUI_chan <- *metrics

            // Check the number of HTTP requests on each X tick
        case <-alertRequestsHTTP_ticker:
            alertMessage := CheckRequestsThreshold(*metrics)
            if alertMessage != "" {
                TrafficAlertToUI_chan <- alertMessage
            }

        }
    }

}
