package main

import (
	"os"
	"log"
	"github.com/hpcloud/tail"
)

var commonLogPath = "access.log"
var commonLogFile *tail.Tail

var errorsLogPath = "errors.log"
var errorsLogFile *os.File
var err error

func init() {
	// Init output error file
	errorsLogFile, err = os.OpenFile(errorsLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err.Error())
	}

	// Init input log reader
	commonLogFile, err = tail.TailFile(commonLogPath, tail.Config{
		Follow: true,
		ReOpen: true,
		Logger: tail.DiscardingLogger,
	})
	if err != nil {
		log.Panic(err.Error())
	}
}

func main() {
	defer errorsLogFile.Close()
	log.SetOutput(errorsLogFile)
	m := NewMonitoringData()

	go m.StartMonitoring()
	go StartLogFileFollower(commonLogFile)

	RunUI()
}
