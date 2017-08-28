package alert

import (
    "time"
    "fmt"
    "http-monitoring/monitoring"
)

const (
    requestsThreshold uint64 = 25
)

var (
    alertTraffic     = "[ %s : High traffic generated an alert - hits = %d](fg-white,bg-red)"
    recoveredTraffic = "[ %s : High traffic recovered](fg-white,bg-green)"
)

var (
    isAlertState          = false
    HighTrafficAlert_chan = make(chan string)
    oldTotalRequests   uint64 = 0
)

func AlertThresholdHTTP() {
    monitoring.MonitoringD.Mux.Lock()
    diff := monitoring.MonitoringD.TotalRequests - oldTotalRequests

    if diff > requestsThreshold {
        isAlertState = true
        HighTrafficAlert_chan <- fmt.Sprintf(alertTraffic, time.Now().String(), diff)

    } else if isAlertState && diff < requestsThreshold {
        isAlertState = false
        HighTrafficAlert_chan <- fmt.Sprintf(recoveredTraffic, time.Now().String())
    }

    oldTotalRequests = monitoring.MonitoringD.TotalRequests
    monitoring.MonitoringD.Mux.Unlock()
}
