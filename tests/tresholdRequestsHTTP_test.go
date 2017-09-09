package tests

import (
	"testing"
	"time"

	"github.com/DanMayor/http-monitoring/alerts"
	"github.com/DanMayor/http-monitoring/monitoring"
)

func Test_AlertThresholdHTTP(t *testing.T) {
	var d uint32 = 2

	monitoring.MonitoringDataStruct.Mux.Lock()
	// Set the number of TotalRequest at the threshold limit to trigger an alert
	monitoring.MonitoringDataStruct.TotalRequests = alerts.RequestsThreshold
	monitoring.MonitoringDataStruct.Mux.Unlock()

	go monitoring.Monitor(d, alerts.AlertThresholdHTTP)

	// IsAlertState must be True
	<-alerts.HighTrafficAlert_chan
	if !alerts.IsAlertState {
		t.Error("Exception raised, should have received a traffic alerts.")
		return
	}

	// IsAlertState must be False because there is no new http request
	<-alerts.HighTrafficAlert_chan
	if alerts.IsAlertState {
		t.Error("Exception raised, should have received a traffic recovery.")
		return
	}

	// IsAlertState must be False again because there is no new http request
	time.Sleep(time.Duration(d) * time.Second)
	if alerts.IsAlertState {
		t.Error("Exception raised, IsAlertState must be at false.")
		return
	}

	close(alerts.HighTrafficAlert_chan)
}
