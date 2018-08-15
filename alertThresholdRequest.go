package main

import (
    "fmt"
    "time"
)

// Number max of HTTP requests before raising an alert, between 2 CheckRequestsThreshold() calls

var (
    trafficAlertMessage     = "[ %s : High traffic generated an alert - hits = %d, threshold = %d](fg-white,bg-red)"
    TrafficRecoveredMessage = "[ %s : High traffic recovered](fg-white,bg-green)"
)

// Used to calculate the number of requests received between 2 CheckRequestsThreshold() calls
// var previousTotalRequests uint64 = 0

/*
   Function supposed to be called each X seconds, sending alert/recovery messages if number of requests
   received during its last 2 calls is too high/low
*/
func (m *Metrics) CheckRequestsThreshold() {
    m.Mux.Lock()
    defer m.Mux.Unlock()
    alertMessage := ""

    // get the number of requests received between the last check (delta = total - previousTotal)
    deltaRequests := m.HttpRequests - m.PreviousTotalRequests

    if deltaRequests >= m.RequestsThreshold {
        m.IsTrafficAlert = true
        alertMessage = fmt.Sprintf(trafficAlertMessage, time.Now().String(), deltaRequests, m.RequestsThreshold)

    } else if m.IsTrafficAlert && deltaRequests < m.RequestsThreshold {
        m.IsTrafficAlert = false
        alertMessage = fmt.Sprintf(TrafficRecoveredMessage, time.Now().String())
    }

    m.PreviousTotalRequests = m.HttpRequests
    m.LastAlertMessage = alertMessage
}
