package main

import (
    "testing"
)

func Test_AlertThresholdHTTP(t *testing.T) {
    metrics := NewMetrics()

    // Set the number of HttpRequests at the threshold limit to trigger an alert
    metrics.HttpRequests = requestsThreshold
    alertMessage := CheckRequestsThreshold(metrics)
    if !isTrafficAlert || alertMessage == ""{
        t.Error("Exception raised, should have received a traffic alert.")
        return
    }

    // Set the number of requests under the threshold limit to trigger a recovery alert
    previousTotalRequests = metrics.HttpRequests
    metrics.HttpRequests += requestsThreshold - 1
    alertMessage = CheckRequestsThreshold(metrics)
    if isTrafficAlert || alertMessage == ""{
        t.Error("Exception raised, should have received a traffic recover.")
        return
    }

    // Still under the threshold limit : alertMessage must be empty, and no alert must be triggered
    alertMessage = CheckRequestsThreshold(metrics)
    if isTrafficAlert || alertMessage != ""{
        t.Error("Exception raised, IsAlertState must be at false.")
        return
    }
}