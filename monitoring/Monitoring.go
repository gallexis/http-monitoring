package monitoring

import (
    "time"
    "github.com/hpcloud/tail"
    "log"
    "sync"
)

type monitoringData struct{
    TotalSize     uint64
    TotalRequests uint64
    MapURLsection map[string]int
    MapHTTPstatus map[string]int
    Mux           sync.Mutex
}

var MonitoringDataStruct = monitoringData{
    TotalSize:     0,
    TotalRequests: 0,
    MapURLsection: make(map[string]int),
    MapHTTPstatus: make(map[string]int) ,
}

var (
    MonitoringDataChan = make(chan monitoringData,1)
)

/*
    we can add many functions which have the same trigger time, i.e: Monitor( 10, f1, f2, f3...)
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

func SendMonitoringDataToUI(){
    MonitoringDataStruct.Mux.Lock()
    MonitoringDataChan <- MonitoringDataStruct
    MonitoringDataStruct.Mux.Unlock()
}

func FollowLogFile(file string){

    t, err := tail.TailFile(file, tail.Config{
        Follow:   true,
        ReOpen:   true,
        Logger:   tail.DiscardingLogger,
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
        UpdateMonitoringDataStruct(s)
    }
}
