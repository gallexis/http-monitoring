package main

import (
	"log"
	"os"

	// This is because of my install path, should be able to delete and save the file to override these
	// to your install paths ~Dan Mayor
	"github.com/DanMayor/http-monitoring/alerts"
	"github.com/DanMayor/http-monitoring/gui"
	"github.com/DanMayor/http-monitoring/monitoring"
)

func main() {
	// Setting log options
	f, err := os.OpenFile("errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Error opening file : " + err.Error())
	}
	defer f.Close()
	log.SetOutput(f)

	go gui.Gui()
	go monitoring.Monitor(10, monitoring.SendMonitoringDataToUI)
	go monitoring.Monitor(120, alerts.AlertThresholdHTTP)

	monitoring.FollowLogFile("l.log", monitoring.UpdateMonitoringDataStruct)
}
