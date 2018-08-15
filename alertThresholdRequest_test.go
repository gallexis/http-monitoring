package main

import (
    "testing"
)

var requestsThreshold = uint64(25)

func Test_AlertThresholdHTTP(t *testing.T) {
    md := NewMonitoringData()

    // Set the number of HttpRequestsCount at the threshold limit to trigger an alert
    md.HttpRequestsCount = requestsThreshold
    md.CheckHttpThresholdCrossed()
    if !md.IsTrafficAlert || md.LastAlertMessage == ""{
        t.Error("Exception raised, should have received a traffic alert.")
        return
    }

    // Set the number of requests under the threshold limit to trigger a recovery alert
    md.PreviousTotalRequests = md.HttpRequestsCount
    md.HttpRequestsCount += requestsThreshold - 1
    md.CheckHttpThresholdCrossed()
    if md.IsTrafficAlert || md.LastAlertMessage == ""{
        t.Error("Exception raised, should have received a traffic recover.")
        return
    }

    // Still under the threshold limit : alertMessage must be empty, and no alert must be triggered
    md.CheckHttpThresholdCrossed()
    if md.IsTrafficAlert || md.LastAlertMessage != ""{
        t.Error("Exception raised, IsAlertState must be at false.")
        return
    }
}