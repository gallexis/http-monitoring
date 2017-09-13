package main

import (
	"github.com/hpcloud/tail"
	"log"
	"regexp"
	"strings"
	"time"
)

var CommonLog_chan = make(chan string)

// Structure representing a line of a Common Log file
type CommonLog struct {
	RemoteHostIP string
	Identd       string
	User_id      string
	Date         time.Time

	Method, Resource, Protocol string

	Status string
	Size   string
}

/*
   Get the section of a url : section for "http://my.site.com/pages/create' is "http://my.site.com/pages"
   Removes also the arguments in a URL, i.e :  /pages?k=v' -> "/pages"
*/
func (cl CommonLog) GetSection() string {
	sections := strings.Split(cl.Resource, "/")
	section  := strings.Split(sections[1], "?")

	return "/" + section[0]
}

// Decode a CommonLog line and return a CommonLog struct
func ParseLine(line string) (CommonLog, error) {
	commonLog := CommonLog{}

	r := regexp.MustCompile(`^(?P<RemoteHostIP>[\d\.]+) (?P<Identd>.*) (?P<user_id>.*) \[(?P<date>.*)\] "(?P<method>.*) (?P<resource>.*) (?P<protocol>.*)" (?P<Status>\d+) (?P<Size>\d+)`)
	fields := r.FindStringSubmatch(line)

	commonLog.RemoteHostIP = fields[1]
	commonLog.Identd = fields[2]
	commonLog.User_id = fields[3]

	date, err := time.Parse("02/Jan/2006:15:04:05 -0700", fields[4])
	if err != nil {
		log.Println(err)
		return commonLog, err
	}
	commonLog.Date = date

	commonLog.Method = fields[5]
	commonLog.Resource = fields[6]
	commonLog.Protocol = fields[7]
	commonLog.Status = fields[8]
	commonLog.Size = fields[9]

	return commonLog, nil
}

/*
   Read a Common Log file, parse each line, convert them as a CommonLog struct, then send it to the channel
*/
func FollowLog(file string) {
	t, err := tail.TailFile(file, tail.Config{
		Follow: true,
		ReOpen: true,
		Logger: tail.DiscardingLogger,
	})
	if err != nil {
		log.Panic(err)
	}

	for line := range t.Lines {
		CommonLog_chan <- line.Text
	}
}
