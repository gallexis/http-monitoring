package main

import (
    "fmt"
    "time"
)

var (
    trafficAlertMessage     = "[ %s : High traffic generated an alert - hits = %d, threshold = %d](fg-white,bg-red)"
    TrafficRecoveredMessage = "[ %s : High traffic recovered](fg-white,bg-green)"
)

/*
   Function supposed to be called each X seconds, sending alert/recovery messages if number of requests
   received during the last 2 calls is too high/low
*/
func (md *MonitoringData) CheckHttpThresholdCrossed() {
    md.Mux.Lock()
    defer md.Mux.Unlock()
    alertMessage := ""

    // get the number of requests received between the last check (delta = total - previousTotal)
    deltaRequests := md.HttpRequestsCount - md.PreviousTotalRequests

    if deltaRequests >= md.RequestsThreshold {
        md.IsTrafficAlert = true
        alertMessage = fmt.Sprintf(trafficAlertMessage, time.Now().String(), deltaRequests, md.RequestsThreshold)

    } else if md.IsTrafficAlert && deltaRequests < md.RequestsThreshold {
        md.IsTrafficAlert = false
        alertMessage = fmt.Sprintf(TrafficRecoveredMessage, time.Now().String())
    }

    md.PreviousTotalRequests = md.HttpRequestsCount
    md.LastAlertMessage = alertMessage
}
