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

var alertState = false

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

    // get the number of the new requests between the last call
    lastRequests := metrics.Requests - previousTotalRequests

    if lastRequests >= requestsThreshold {
        alertState = true
        alertMessage = fmt.Sprintf(messageTrafficAlert, time.Now().String(), lastRequests, requestsThreshold)

    } else if alertState && lastRequests < requestsThreshold {
        alertState = false
        alertMessage = fmt.Sprintf(messageTrafficRecovered, time.Now().String())
    }

    previousTotalRequests = metrics.Requests
    return alertMessage
}
