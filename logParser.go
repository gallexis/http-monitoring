package main

import (
		"log"
	"regexp"
	"strings"
	"time"
    )

var LogLineChan = make(chan string)

// Structure representing a line of a log file
type LogLine struct {
	RemoteHostIP string
	Identd       string
	UserID       string
	Method, Resource, Protocol string
	Status string
	Size   string
	Date         time.Time
}

/*
   Get the section of a url : section for "http://my.site.com/pages/create' is "http://my.site.com/pages"
   Removes also the arguments in a URL, i.e :  /pages?k=v' -> "/pages"
*/
func (l LogLine) GetSection() string {
	sections := strings.Split(l.Resource, "/")
	section  := strings.Split(sections[1], "?")

	return "/" + section[0]
}

// Decode a log line and return a LogLine struct
func ParseLine(line string) (LogLine, error) {
	logLine := LogLine{}

	r := regexp.MustCompile(`^(?P<RemoteHostIP>[\d\.]+) (?P<Identd>.*) (?P<user_id>.*) \[(?P<date>.*)\] "(?P<method>.*) (?P<resource>.*) (?P<protocol>.*)" (?P<Status>\d+) (?P<Size>\d+)`)
	fields := r.FindStringSubmatch(line)

	logLine.RemoteHostIP = fields[1]
	logLine.Identd = fields[2]
	logLine.UserID = fields[3]
	logLine.Method = fields[5]
	logLine.Resource = fields[6]
	logLine.Protocol = fields[7]
	logLine.Status = fields[8]
	logLine.Size = fields[9]

    date, err := time.Parse("02/Jan/2006:15:04:05 -0700", fields[4])
    if err != nil {
        log.Println(err.Error())
    }
    logLine.Date = date

	return logLine, nil
}