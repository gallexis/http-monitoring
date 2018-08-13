package main

import (
    "fmt"
    "time"
)

// Number max of HTTP requests before raising an alert, between 2 CheckRequestsThreshold() calls
const (
    requestsThreshold uint64 = 25
)

var (
    trafficAlertMessage     = "[ %s : High traffic generated an alert - hits = %d, threshold = %d](fg-white,bg-red)"
    TrafficRecoveredMessage = "[ %s : High traffic recovered](fg-white,bg-green)"
)

var isTrafficAlert = false

// Used to calculate the number of requests received between 2 CheckRequestsThreshold() calls
var previousTotalRequests uint64 = 0

/*
   Function supposed to be called each X seconds, sending alert/recovery messages if number of requests
   received during its last 2 calls is too high/low
*/
func CheckRequestsThreshold(metrics Metrics) string {
    metrics.Mux.Lock()
    defer metrics.Mux.Unlock()
    alertMessage := ""

    // get the number of requests received between the last check (delta = total - previousTotal)
    deltaRequests := metrics.HttpRequests - previousTotalRequests

    if deltaRequests >= requestsThreshold {
        isTrafficAlert = true
        alertMessage = fmt.Sprintf(trafficAlertMessage, time.Now().String(), deltaRequests, requestsThreshold)

    } else if isTrafficAlert && deltaRequests < requestsThreshold {
        isTrafficAlert = false
        alertMessage = fmt.Sprintf(TrafficRecoveredMessage, time.Now().String())
    }

    previousTotalRequests = metrics.HttpRequests
    return alertMessage
}
