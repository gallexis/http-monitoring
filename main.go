package main

import (
    "http-monitoring/monitoring"
    "http-monitoring/gui"
    "http-monitoring/alerts"
    "os"
    "log"
)

func main() {

    // Setting log options
    f, err := os.OpenFile("errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        panic("Error opening file : " + err.Error())
    }
    defer f.Close()
    log.SetOutput(f)

    go gui.Gui()
    go monitoring.Monitor(10, monitoring.SendMonitoringDataToUI)
    go monitoring.Monitor(120, alerts.AlertThresholdHTTP)

    monitoring.FollowLogFile("l.log", monitoring.UpdateMonitoringDataStruct)
}
