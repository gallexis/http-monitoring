package main

import (
    "log"
    "time"
    "regexp"
    "github.com/hpcloud/tail"
    "strings"
    "strconv"
    "bytes"
)

var (
    totalSize  uint64 = 0
    totalRequests_2min  uint64 = 0
    totalRequests       uint64 = 0
    isAlertState               = false
    mapSection                 = make(map[string]int)
    httpStatus                 = make(map[string]int)

    Stats_chan                 = make(chan []string)
    HighTrafficAlert_chan      = make(chan string)
    lineLog_chan               = make(chan string)
)

const (
    requestsTreshold uint64 = 25
)

type lineStruct struct{
    remote_host_ip string
    identd         string
    user_id        string
    date           time.Time

    method, resource, protocol string

    status         string
    size           string
}
/*
    we can add many function which have the same trigger time, i.e: monitor( 10, f1, f2...)
 */
func monitor(seconds uint32, fs ...func()){
    duration := time.Duration(seconds) * time.Second
    for {
        time.Sleep(duration)
        for _,f := range fs{
            go f()
        }
    }
}

func DisplayMostViewedSections() string{
    var buffer bytes.Buffer

    for k,v := range mapSection {
        buffer.WriteString(k)
        buffer.WriteString(" : ")
        buffer.WriteString(strconv.Itoa(v))
        buffer.WriteString(" - ")
    }

    return buffer.String()
}

func DisplayTotalHTTPRequests() string{
    return "Total HTTP requests : " + strconv.FormatUint(totalRequests, 10)
}

func DisplayTotalSizeEmitted() string{
    return "Total Size emitted in bytes : " + strconv.FormatUint(totalSize, 10)
}

func manageStatistics(){
    stats := []string{
        DisplayTotalHTTPRequests(),
        DisplayTotalSizeEmitted(),
        DisplayMostViewedSections(),
    }

    Stats_chan <- stats
}

func manageAlerts(){
    if totalRequests_2min > requestsTreshold {
        isAlertState = true
        HighTrafficAlert_chan <- "[" + time.Now().String() + " : High traffic generated an alert - hits = " +
            strconv.FormatUint(totalRequests_2min, 10) + "](fg-white,bg-red)"

    } else if isAlertState && totalRequests_2min < requestsTreshold {
        isAlertState = false
        HighTrafficAlert_chan <- "[" + time.Now().String() + " : High traffic recovered.](fg-white,bg-green)"
    }

    totalRequests_2min = 0
}

func getSection(url string) string {
    sections := strings.Split(url, "/")
    section  := strings.Split(sections[1], "?")

    return "/"+section[0]
}

func lineToLogStruct(line string) (lineStruct, error){
    str := lineStruct{}

    r := regexp.MustCompile(`^(?P<remote_host_ip>[\d\.]+) (?P<identd>.*) (?P<user_id>.*) \[(?P<date>.*)\] "(?P<method>.*) (?P<resource>.*) (?P<protocol>.*)" (?P<status>\d+) (?P<size>\d+)`)
    fields := r.FindStringSubmatch(line)

    str.remote_host_ip = fields[1]
    str.identd = fields[2]
    str.user_id = fields[3]

    t, err := time.Parse("02/Jan/2006:15:04:05 -0700", fields[4])
    if err != nil{
        return str, err
    }
    str.date = t
    str.method = fields[5]
    str.resource = fields[6]
    str.protocol = fields[7]
    str.status = fields[8]
    str.size = fields[9]

    return str,nil
}

func FollowLogFile(file string){

    t, err := tail.TailFile(file, tail.Config{
        Follow:   true,
        ReOpen:   true,
        Logger: tail.DiscardingLogger,
    })
    if err != nil {
        log.Panic(err)
    }

    for line := range t.Lines {
        lineLog_chan <- line.Text
        s, _ := lineToLogStruct(line.Text)
        section := getSection(s.resource)

        // HTTP status counter (i.e : 200, 404, 500...)
        httpStatus[s.status] += 1

        // URL section counter (i.e : /pages : 4 )
        mapSection[section]  += 1

        // Increase number of requests (set at 0 every 2min)
        totalRequests_2min   += 1

        // Increase number of requests (set at 0 every 2min)
        totalRequests   += 1

        size, err := strconv.ParseUint(s.size, 10, 64)
        if err != nil {
            panic(err)
        }
        totalSize += size
    }
}

func startMonitoring(){

    go monitor(2, manageStatistics)
    go monitor(10, manageAlerts)

}

func main(){

    go Gui()
    startMonitoring()
    FollowLogFile("l.log")

}