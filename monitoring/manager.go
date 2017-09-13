package monitoring

import (
    "http-monitoring/alert"
    "http-monitoring/data"
    "http-monitoring/parser"
    "log"
    "time"
)

// Tickers triggered each X seconds, caught by the Start() function
var metricsToUI_ticker = time.NewTicker(time.Second * 5).C
var alertRequestsHTTP_ticker = time.NewTicker(time.Second * 10).C

// Channels
var CommonLogToUI_chan = make(chan string, 1)
var MetricStructToUI_chan = make(chan data.MetricStruct, 1)
var TrafficAlertToUI_chan = make(chan string, 1)

func Start(metricStruct *data.MetricStruct) {
    for {

        select {

            // Receive new line from the log file
        case line := <-parser.CommonLog_chan:

            // Send this line to the UI
            CommonLogToUI_chan <- line

            // Convert the line to a CommonLogStruct
            logStruct, err := parser.ParseLine(line)
            if err != nil {
                log.Panic(err)
            }

            // Update the metricStruct with the new data from logStruct
            metricStruct.Update(logStruct)

            // Send monitoring stats to UI each X seconds
        case <-metricsToUI_ticker:
            MetricStructToUI_chan <- *metricStruct

            // Check for monitoring alert each X seconds
        case <-alertRequestsHTTP_ticker:
            alertMessage := alert.CheckRequestsThreshold(*metricStruct)
            if alertMessage != "" {
                TrafficAlertToUI_chan <- alertMessage
            }

        }
    }

}
