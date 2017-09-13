package main

import (
    "testing"
)

func Test_AlertThresholdHTTP(t *testing.T) {
    metrics := NewMetricStruct()

    metrics.TotalRequests = RequestsThreshold
    alertMessage := CheckRequestsThreshold(metrics)
    if !IsAlertState || alertMessage == ""{
        t.Error("Exception raised, should have received a traffic alert.")
        return
    }

    previousTotalRequests = metrics.TotalRequests
    metrics.TotalRequests += RequestsThreshold - 1
    alertMessage = CheckRequestsThreshold(metrics)
    if IsAlertState || alertMessage == ""{
        t.Error("Exception raised, should have received a traffic recover.")
        return
    }

    alertMessage = CheckRequestsThreshold(metrics)
    if IsAlertState || alertMessage != ""{
        t.Error("Exception raised, IsAlertState must be at false.")
        return
    }
}