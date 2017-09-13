package parser

import (
	"github.com/hpcloud/tail"
	"log"
	"regexp"
	"strings"
	"time"
)

var CommonLog_chan = make(chan string)

// Structure representing a line of a Common Log file
type CommonLogStruct struct {
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
func (s CommonLogStruct) GetSection() string {
	sections := strings.Split(s.Resource, "/")
	section := strings.Split(sections[1], "?")

	return "/" + section[0]
}

// Decode a CommonLog line and return a CommonLogStruct
func ParseLine(line string) (CommonLogStruct, error) {
	logStruct := CommonLogStruct{}

	r := regexp.MustCompile(`^(?P<RemoteHostIP>[\d\.]+) (?P<Identd>.*) (?P<user_id>.*) \[(?P<date>.*)\] "(?P<method>.*) (?P<resource>.*) (?P<protocol>.*)" (?P<Status>\d+) (?P<Size>\d+)`)
	fields := r.FindStringSubmatch(line)

	logStruct.RemoteHostIP = fields[1]
	logStruct.Identd = fields[2]
	logStruct.User_id = fields[3]

	date, err := time.Parse("02/Jan/2006:15:04:05 -0700", fields[4])
	if err != nil {
		log.Println(err)
		return logStruct, err
	}
	logStruct.Date = date

	logStruct.Method = fields[5]
	logStruct.Resource = fields[6]
	logStruct.Protocol = fields[7]
	logStruct.Status = fields[8]
	logStruct.Size = fields[9]

	return logStruct, nil
}

/*
   Read a Common Log file, parse each line, convert them as a CommonLogStruct, then send it to the channel
*/
func Follow(file string) {
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
