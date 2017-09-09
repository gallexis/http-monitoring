package alerts

import (
	"fmt"
	"time"

	"github.com/DanMayor/http-monitoring/monitoring"
)

// Number max of HTTP requests before raising an alert, between 2 AlertThresholdHTTP() calls
const (
	RequestsThreshold uint64 = 25
)

var (
	alertTraffic     = "[ %s : High traffic generated an alert - hits = %d, threshold = %d](fg-white,bg-red)"
	recoveredTraffic = "[ %s : High traffic recovered](fg-white,bg-green)"
)

var IsAlertState = false

// Channel used to send the high traffic alert/recovery to the GUI
var HighTrafficAlert_chan = make(chan string, 1)

// Used to calculate the number of requests received between 2 AlertThresholdHTTP() calls
var previousTotalRequests uint64 = 0

/*
   Function supposed to be called each X seconds, sending alert/recovery messages if number of requests
   received during its last 2 calls is too high/low
   (send nothing if IsAlertState == false and lastNumberRequests < RequestsThreshold)
*/
func AlertThresholdHTTP() {
	monitoring.MonitoringDataStruct.Mux.Lock()
	defer monitoring.MonitoringDataStruct.Mux.Unlock()

	lastNumberRequests := monitoring.MonitoringDataStruct.TotalRequests - previousTotalRequests

	if lastNumberRequests >= RequestsThreshold {
		IsAlertState = true
		HighTrafficAlert_chan <- fmt.Sprintf(alertTraffic, time.Now().String(), lastNumberRequests, RequestsThreshold)

	} else if IsAlertState && lastNumberRequests < RequestsThreshold {
		IsAlertState = false
		HighTrafficAlert_chan <- fmt.Sprintf(recoveredTraffic, time.Now().String())
	}

	previousTotalRequests = monitoring.MonitoringDataStruct.TotalRequests
}
