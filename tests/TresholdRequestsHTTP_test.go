package tests

import (
    "testing"
    "http-monitoring/alert"
    "http-monitoring/monitoring"
    "time"
)

func Test_AlertThresholdHTTP(t *testing.T) {
    var d uint32 = 2
    monitoring.MonitoringDataStruct.Mux.Lock()
    monitoring.MonitoringDataStruct.TotalRequests = alert.RequestsThreshold
    monitoring.MonitoringDataStruct.Mux.Unlock()

    // Generate Alert
    go monitoring.Monitor(d, alert.AlertThresholdHTTP)

    <- alert.HighTrafficAlert_chan
    if !alert.IsAlertState{
        t.Error("Exception raised, should have received a traffic alert.")
        return
    }

    <- alert.HighTrafficAlert_chan
    if alert.IsAlertState{
        t.Error("Exception raised, should have received a traffic recover.")
        return
    }

    time.Sleep(time.Duration(d) * time.Second)
    if alert.IsAlertState{
        t.Error("Exception raised, IsAlertState must be at false.")
        return
    }

    close(alert.HighTrafficAlert_chan)
}
