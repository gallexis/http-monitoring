package monitoring

import (
    "time"
    "strings"
    "regexp"
    "log"
    "github.com/hpcloud/tail"
)

// Channel used to send the LineLogStruct
var LineLog_chan = make(chan string, 1)

// Structure representing a line of a Common Log file
type LineLogStruct struct {
    Remote_host_ip string
    Identd         string
    User_id        string
    Date           time.Time

    Method, Resource, Protocol string

    Status string
    Size   string
}

/*
    Get the section of a url : section for "http://my.site.com/pages/create' is "http://my.site.com/pages"
    Removes also the arguments in a URL, i.e :  /pages?k=v' -> "/pages"
*/
func GetURLSection(s LineLogStruct) string {
    sections := strings.Split(s.Resource, "/")
    section := strings.Split(sections[1], "?")

    return "/" + section[0]
}

// Decode a commonLog line and return a LineLogStruct from this line
func lineToLogStruct(line string) (LineLogStruct, error) {
    logStruct := LineLogStruct{}

    r := regexp.MustCompile(`^(?P<Remote_host_ip>[\d\.]+) (?P<Identd>.*) (?P<user_id>.*) \[(?P<date>.*)\] "(?P<method>.*) (?P<resource>.*) (?P<protocol>.*)" (?P<Status>\d+) (?P<Size>\d+)`)
    fields := r.FindStringSubmatch(line)

    logStruct.Remote_host_ip = fields[1]
    logStruct.Identd = fields[2]
    logStruct.User_id = fields[3]

    t, err := time.Parse("02/Jan/2006:15:04:05 -0700", fields[4])
    if err != nil {
        log.Println(err)
        return logStruct, err
    }
    logStruct.Date = t
    logStruct.Method = fields[5]
    logStruct.Resource = fields[6]
    logStruct.Protocol = fields[7]
    logStruct.Status = fields[8]
    logStruct.Size = fields[9]

    return logStruct, nil
}

/*
    Read a Common Log file, parse each line, convert them as a LineLogStruct, then send it via the functions in fs
    i.e : we can send the structure to UpdateMonitoringDataStruct(), but we can send it elsewhere too
 */
func FollowLogFile(file string, fs ...func(LineLogStruct)) {
    t, err := tail.TailFile(file, tail.Config{
        Follow: true,
        ReOpen: true,
        Logger: tail.DiscardingLogger,
    })
    if err != nil {
        log.Panic(err)
    }

    for line := range t.Lines {
        // Send the raw text line to be displayed in the UI
        LineLog_chan <- line.Text

        // Convert the line to a LineLogStruct
        logStruct, err := lineToLogStruct(line.Text)
        if err != nil {
            log.Println(err)
            continue
        }

        for _, f := range fs {
            go f(logStruct)
        }
    }
}
