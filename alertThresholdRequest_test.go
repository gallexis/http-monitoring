package main

import (
    "testing"
)

func Test_AlertThresholdHTTP(t *testing.T) {
    metrics := NewMetricStruct()

    // Set the number of Requests at the threshold limit to trigger an alert
    metrics.Requests = requestsThreshold
    alertMessage := CheckRequestsThreshold(metrics)
    if !trafficAlert || alertMessage == ""{
        t.Error("Exception raised, should have received a traffic alert.")
        return
    }

    // Set the number of requests under the threshold limit to trigger a recovery alert
    previousRequests = metrics.Requests
    metrics.Requests += requestsThreshold - 1
    alertMessage = CheckRequestsThreshold(metrics)
    if trafficAlert || alertMessage == ""{
        t.Error("Exception raised, should have received a traffic recover.")
        return
    }

    // Still under the threshold limit : alertMessage must be empty, and no alert must be triggered
    alertMessage = CheckRequestsThreshold(metrics)
    if trafficAlert || alertMessage != ""{
        t.Error("Exception raised, IsAlertState must be at false.")
        return
    }
}