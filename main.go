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
    totalRequests_2min  uint32 = 0
    isAlertState               = false
    mapSection                 = make(map[string]int)
    httpStatus                 = make(map[string]int)
)

const (
    requestsTrshld_2min uint32 = 5
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

func monitor(seconds uint32, fs ...func()){
    duration := time.Duration(seconds) * time.Second
    for {
        time.Sleep(duration)
        for _,f := range fs{
            go f()
        }
    }
}

func browseMostViewedSections(){
    var buffer bytes.Buffer

    for k,v := range mapSection {
        buffer.WriteString("\n")
        buffer.WriteString(k)
        buffer.WriteString(" : ")
        buffer.WriteString(strconv.Itoa(v))
    }
    log.Println(buffer.String())
}

func displayTotalHTTPRequests(){
    log.Println("Total HTTP requests :", totalRequests_2min)
}

func displayTotalSizeEmitted(){
    log.Println("Total Size emitted in bytes :", totalSize)
}

func manageAlerts(){
    if totalRequests_2min > requestsTrshld_2min {
        isAlertState = true
        log.Println("High traffic generated an alert - hits =", totalRequests_2min,", triggered at", time.Now())

    } else if isAlertState && totalRequests_2min < requestsTrshld_2min{
        isAlertState = false
        log.Println("High traffic recovered at", time.Now())
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
    })
    if err != nil {
        log.Panic(err)
    }

    for line := range t.Lines {
        s, _ := lineToLogStruct(line.Text)
        section := getSection(s.resource)

        // HTTP status counter (i.e : 200, 404, 500...)
        httpStatus[s.status] += 1

        // URL section counter (i.e : /pages : 4 )
        mapSection[section]  += 1

        // Increase number of requests (set at 0 every 2min)
        totalRequests_2min   += 1

        size, err := strconv.ParseUint(s.size, 10, 64)
        if err != nil {
            panic(err)
        }
        totalSize += size
    }
}

func startMonitoring(){
    quickMonitoring := []func(){
        browseMostViewedSections,
        displayTotalSizeEmitted,
        displayTotalHTTPRequests,
    }

    go monitor(2, quickMonitoring...)
    go monitor(10, manageAlerts)
}

func main(){

    startMonitoring()
    FollowLogFile("l.log")

}