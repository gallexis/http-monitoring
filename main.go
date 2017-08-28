package main

import (
    "http-monitoring/monitoring"
    "http-monitoring/gui"
    "http-monitoring/alert"
)

func main(){

    go monitoring.Monitor(2,  monitoring.SendMonitoringToUI)
    go monitoring.Monitor(10, alert.AlertThresholdHTTP)

    go gui.Gui()
    monitoring.FollowLogFile("l.log")

}