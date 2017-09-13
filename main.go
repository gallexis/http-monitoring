package main

import (
	"http-monitoring/UI"
	"http-monitoring/data"
	"http-monitoring/monitoring"
	"http-monitoring/parser"
	"log"
	"os"
)

func main() {
	// Set log options
	f, err := os.OpenFile("errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Error opening file : " + err.Error())
	}
	defer f.Close()
	log.SetOutput(f)

	metricsStruct := data.NewMetricStruct()

	go monitoring.Start(&metricsStruct)
	go parser.Follow("l.log")

	UI.Start()
}
