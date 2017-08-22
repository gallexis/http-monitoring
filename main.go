package main

import (
    "fmt"
    //"github.com/hpcloud/tail"
    "time"
    "regexp"
    "github.com/hpcloud/tail"
)

type lineLog struct{
    remote_host_ip string
    identd         string
    user_id        string
    date           time.Time

    method, ressource, protocol        string

    status         string
    size           string
}

func parser(line string) (lineLog, error){
    str := lineLog{}

    r := regexp.MustCompile(`^(?P<remote_host_ip>[\d\.]+) (?P<identd>.*) (?P<user_id>.*) \[(?P<date>.*)\] "(?P<method>.*) (?P<ressource>.*) (?P<protocol>.*)" (?P<status>\d+) (?P<size>\d+)`)
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
    str.ressource = fields[6]
    str.protocol = fields[7]
    str.status = fields[8]
    str.size = fields[9]

    fmt.Println(str)

    return str,nil
}

func main(){

    //s := `127.0.0.1 useridentifier frank [10/Mar/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
    //parser(s)

    t, err := tail.TailFile("l.log", tail.Config{
        Follow:   true,
        ReOpen:   true,
    })
    if err != nil {
        fmt.Println(err)
    }

    for line := range t.Lines {
        parser(line.Text)
    }

}
