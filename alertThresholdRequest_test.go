package main

import (
    "testing"
)

func Test_AlertThresholdHTTP(t *testing.T) {
    metrics := NewMetricStruct()

    // Set the number of TotalRequest at the threshold limit to trigger an alert
    metrics.Requests = requestsThreshold
    alertMessage := CheckRequestsThreshold(metrics)
    if !alertState || alertMessage == ""{
        t.Error("Exception raised, should have received a traffic alert.")
        return
    }

    // Set the number of TotalRequest under the threshold limit to trigger a recovery alert
    previousTotalRequests = metrics.Requests
    metrics.Requests += requestsThreshold - 1
    alertMessage = CheckRequestsThreshold(metrics)
    if alertState || alertMessage == ""{
        t.Error("Exception raised, should have received a traffic recover.")
        return
    }

    // Still under the threshold limit : alertMessage must be empty, and no alert to be triggered
    alertMessage = CheckRequestsThreshold(metrics)
    if alertState || alertMessage != ""{
        t.Error("Exception raised, IsAlertState must be at false.")
        return
    }
}