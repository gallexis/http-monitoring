package main

import (
    "fmt"
    "time"
    "regexp"
    "github.com/hpcloud/tail"
    "strings"
)

var (
    totalRequests uint32 = 0
    sectionMap           = make(map[string]int)
    httpStatus           = make(map[string]int)
)

type lineLog struct{
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
            fmt.Println("-----------------")
        }
    }
}

func browseMostViewedSections(){
    for k,v := range sectionMap{
        fmt.Println(k,v)
    }
}

func displayTotalViews(){
    fmt.Println("Total requests: ", totalRequests)
}

func getSection(url string) string {
    sections := strings.Split(url, "/")
    section := strings.Split(sections[1], "?")

    return "/"+section[0]
}

func lineParser(line string) (lineLog, error){
    str := lineLog{}

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

func main(){
    quickMonitoring := []func(){
        browseMostViewedSections,
        displayTotalViews,
    }

    go monitor(2, quickMonitoring...)

    t, err := tail.TailFile("l.log", tail.Config{
        Follow:   true,
        ReOpen:   true,
    })
    if err != nil {
        fmt.Println(err)
        return
    }

    for line := range t.Lines {
        s, _ := lineParser(line.Text)
        section := getSection(s.resource)

        httpStatus[s.status] += 1
        sectionMap[section]  += 1
        totalRequests        += 1
    }
}