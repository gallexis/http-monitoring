package main

import (
    "testing"
)

func Test_AlertThresholdHTTP(t *testing.T) {
    md := NewMonitoringData(25)

    // Set the number of HttpRequestsCount at the threshold limit to trigger an alert
    md.HttpRequestsCount = md.RequestsThreshold
    md.CheckHttpThresholdCrossed()
    if !md.IsTrafficAlert || md.LastAlertMessage == "" {
        t.Error("Exception raised, should have received a traffic alert.")
        return
    }

    // Set the number of requests under the threshold limit to trigger a recovery alert
    md.PreviousTotalRequests = md.HttpRequestsCount
    md.HttpRequestsCount += md.RequestsThreshold - 1
    md.CheckHttpThresholdCrossed()
    if md.IsTrafficAlert || md.LastAlertMessage == "" {
        t.Error("Exception raised, should have received a traffic recover.")
        return
    }

    // Still under the threshold limit : alertMessage must be empty, and no alert must be triggered
    md.CheckHttpThresholdCrossed()
    if md.IsTrafficAlert || md.LastAlertMessage != "" {
        t.Error("Exception raised, IsAlertState must be at false.")
        return
    }
}
