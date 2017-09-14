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
    messageTrafficAlert     = "[ %s : High traffic generated an alert - hits = %d, threshold = %d](fg-white,bg-red)"
    messageTrafficRecovered = "[ %s : High traffic recovered](fg-white,bg-green)"
)

var trafficAlert = false

// Used to calculate the number of requests received between 2 CheckRequestsThreshold() calls
var previousRequests uint64 = 0

/*
   Function supposed to be called each X seconds, sending alert/recovery messages if number of requests
   received during its last 2 calls is too high/low
*/
func CheckRequestsThreshold(metrics Metrics) string {
    metrics.Mux.Lock()
    defer metrics.Mux.Unlock()
    alertMessage := ""

    // get the number of the requests received between the last check (delta = total - previous)
    deltaRequests := metrics.Requests - previousRequests

    if deltaRequests >= requestsThreshold {
        trafficAlert = true
        alertMessage = fmt.Sprintf(messageTrafficAlert, time.Now().String(), deltaRequests, requestsThreshold)

    } else if trafficAlert && deltaRequests < requestsThreshold {
        trafficAlert = false
        alertMessage = fmt.Sprintf(messageTrafficRecovered, time.Now().String())
    }

    previousRequests = metrics.Requests
    return alertMessage
}
