package monitoring

import (
    "time"
    "strings"
    "regexp"
)

var (
    LineLog_chan               = make(chan string)
)

type LineStruct struct{
    Remote_host_ip string
    Identd         string
    User_id        string
    Date           time.Time

    Method, Resource, Protocol string

    Status string
    Size   string
}

func GetURLSection(s LineStruct) string {
    sections := strings.Split(s.Resource, "/")
    section  := strings.Split(sections[1], "?")

    return "/"+section[0]
}

func lineToLogStruct(line string) (LineStruct, error){
    str := LineStruct{}

    r := regexp.MustCompile(`^(?P<Remote_host_ip>[\d\.]+) (?P<Identd>.*) (?P<user_id>.*) \[(?P<date>.*)\] "(?P<method>.*) (?P<resource>.*) (?P<protocol>.*)" (?P<Status>\d+) (?P<Size>\d+)`)
    fields := r.FindStringSubmatch(line)

    str.Remote_host_ip = fields[1]
    str.Identd = fields[2]
    str.User_id = fields[3]

    t, err := time.Parse("02/Jan/2006:15:04:05 -0700", fields[4])
    if err != nil{
        return str, err
    }
    str.Date = t
    str.Method = fields[5]
    str.Resource = fields[6]
    str.Protocol = fields[7]
    str.Status = fields[8]
    str.Size = fields[9]

    return str,nil
}
