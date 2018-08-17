package main

import (
    "log"
    "time"
    "github.com/hpcloud/tail"
)

// UI Channels
var DisplayLogLineChan = make(chan MonitoringData, 1)
var DisplayMonitoringDataChan = make(chan MonitoringData, 1)
var DisplayTrafficAlertChan = make(chan MonitoringData, 1)

//Read a log file, parse each line then send them to the LogLine channel
func FileFollower(commonLogFile *tail.Tail) {
    for line := range commonLogFile.Lines {
        LogLineChan <- line.Text
    }
}

func (md *MonitoringData) StartMonitoring(monitoringDataDelay, alertsDelay int) {
    // Tickers triggered each X seconds
    tickerMonitoringData := time.NewTicker(time.Second * time.Duration(monitoringDataDelay)).C
    tickerAlerts := time.NewTicker(time.Second * time.Duration(alertsDelay)).C
    DisplayMonitoringDataChan <- *md

    for {
        select {

        // Receive new line from the log file
        case line := <-LogLineChan:
            // Add this raw string line to the Ring
            md.LogLines.Value = line
            md.LogLines = md.LogLines.Next()

            // Then convert it to a LogLine struct
            logLine, err := ParseLine(line)
            if err != nil {
                log.Println(err.Error())
            }

            // Update the monitoring data struct with the new line received ...
            md.Update(logLine)

            // ... then send it to the UI
            DisplayLogLineChan <- *md

            // Send monitoring stats to UI on each tick
        case <-tickerMonitoringData:
            DisplayMonitoringDataChan <- *md

            // Check the number of HTTP requests on each X tick
        case <-tickerAlerts:
            md.CheckHttpThresholdCrossed()

            if md.LastAlertMessage != ""{
                md.Alerts.Value = md.LastAlertMessage
                md.Alerts = md.Alerts.Next()
                DisplayTrafficAlertChan <- *md
            }
        }
    }
}
