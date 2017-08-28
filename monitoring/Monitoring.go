package monitoring

import (
    "time"
    "github.com/hpcloud/tail"
    "log"
    "sync"
)

type MonitoringData struct{
    TotalSize     uint64
    TotalRequests uint64
    MapSection    map[string]int
    HttpStatus    map[string]int
    Mux           sync.Mutex
}

var MonitoringD MonitoringData = MonitoringData{
    TotalSize:0,
    TotalRequests:0,
    MapSection: make( map[string]int),
    HttpStatus:make(map[string]int) ,
}

var (
    Stats_chan = make(chan MonitoringData)
)

/*
    we can add many function which have the same trigger time, i.e: monitor( 10, f1, f2...)
 */
func Monitor(seconds uint32, fs ...func()){
    duration := time.Duration(seconds) * time.Second
    for {
        time.Sleep(duration)
        for _,f := range fs{
            go f()
        }
    }
}

func SendMonitoringToUI(){
    MonitoringD.Mux.Lock()
    Stats_chan <- MonitoringD
    MonitoringD.Mux.Unlock()
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
        LineLog_chan <- line.Text
        s, err := lineToLogStruct(line.Text)
        if err != nil{
            // log
            continue
        }

        UpdateData(s)
    }
}
