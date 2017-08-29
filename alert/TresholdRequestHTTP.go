package alert

import (
    "time"
    "fmt"
    "http-monitoring/monitoring"
)

const (
    RequestsThreshold uint64 = 25
)

var (
    alertTraffic     = "[ %s : High traffic generated an alert - hits = %d](fg-white,bg-red)"
    recoveredTraffic = "[ %s : High traffic recovered](fg-white,bg-green)"
)

var (
    IsAlertState                 = false
    HighTrafficAlert_chan        = make(chan string, 1)
    oldTotalRequests      uint64 = 0
)

func AlertThresholdHTTP() {
    monitoring.MonitoringDataStruct.Mux.Lock()
    diff := monitoring.MonitoringDataStruct.TotalRequests - oldTotalRequests

    if diff >= RequestsThreshold {
        IsAlertState = true
        HighTrafficAlert_chan <- fmt.Sprintf(alertTraffic, time.Now().String(), diff)

    } else if IsAlertState && diff < RequestsThreshold {
        IsAlertState = false
        HighTrafficAlert_chan <- fmt.Sprintf(recoveredTraffic, time.Now().String())
    }

    oldTotalRequests = monitoring.MonitoringDataStruct.TotalRequests
    monitoring.MonitoringDataStruct.Mux.Unlock()
}
