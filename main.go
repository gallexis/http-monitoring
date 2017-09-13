package main

import (
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

	metricsStruct := NewMetricStruct()

	go Monitoring(metricsStruct)
	go FollowLog("l.log")

    RunUI()
}
