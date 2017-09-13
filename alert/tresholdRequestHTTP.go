package alert

import (
    "fmt"
    "http-monitoring/data"
    "time"
)

// Number max of HTTP requests before raising an alert, between 2 CheckRequestsThreshold() calls
const (
    RequestsThreshold uint64 = 25
)

var (
    messageTrafficAlert     = "[ %s : High traffic generated an alert - hits = %d, threshold = %d](fg-white,bg-red)"
    messageTrafficRecovered = "[ %s : High traffic recovered](fg-white,bg-green)"
)

var IsAlertState = false

// Used to calculate the number of requests received between 2 CheckRequestsThreshold() calls
var previousTotalRequests uint64 = 0

/*
   Function supposed to be called each X seconds, sending alert/recovery messages if number of requests
   received during its last 2 calls is too high/low
   (send nothing if IsAlertState == false and lastNumberRequests < RequestsThreshold)
*/
func CheckRequestsThreshold(ms data.MetricStruct) string {
    ms.Mux.Lock()
    defer ms.Mux.Unlock()
    alertMessage := ""

    lastNumberRequests := ms.TotalRequests - previousTotalRequests

    if lastNumberRequests >= RequestsThreshold {
        IsAlertState = true
        alertMessage = fmt.Sprintf(messageTrafficAlert, time.Now().String(), lastNumberRequests, RequestsThreshold)

    } else if IsAlertState && lastNumberRequests < RequestsThreshold {
        IsAlertState = false
        alertMessage = fmt.Sprintf(messageTrafficRecovered, time.Now().String())
    }

    previousTotalRequests = ms.TotalRequests
    return alertMessage
}
