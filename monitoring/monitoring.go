package monitoring

import (
    "time"
    "sync"
    "strconv"
    "log"
)

/*
    Structure representing the metrics to be monitored
    Used between many goroutines, so don't forget the mutex (Mux)
*/
type monitoringData struct {
    TotalSize     uint64
    TotalRequests uint64
    MapURLsection map[string]int
    MapHTTPstatus map[string]int
    Mux           sync.Mutex
}

// Initialisation of the monitoringData struct.
var MonitoringDataStruct = monitoringData{
    TotalSize:     0,
    TotalRequests: 0,
    MapURLsection: make(map[string]int),
    MapHTTPstatus: make(map[string]int),
}

// Channel used to send the MonitoringDataStruct to the UI
var MonitoringData_chan = make(chan monitoringData, 1)

// We can add many functions with the same trigger time, i.e: Monitor( 10, f1, f2, f3...)
func Monitor(seconds uint32, fs ...func()) {
    duration := time.Duration(seconds) * time.Second
    for {
        time.Sleep(duration)
        for _, f := range fs {
            go f()
        }
    }
}

func UpdateMonitoringDataStruct(s LineLogStruct) {
    MonitoringDataStruct.Mux.Lock()
    defer MonitoringDataStruct.Mux.Unlock()

    section := GetURLSection(s)

    // URL section counter (i.e : /pages : 4 )
    MonitoringDataStruct.MapURLsection[section] += 1

    // HTTP status counter (i.e : 200, 404, 500...)
    MonitoringDataStruct.MapHTTPstatus[s.Status] += 1

    // Increase counter of HTTP requests
    MonitoringDataStruct.TotalRequests += 1

    size, err := strconv.ParseUint(s.Size, 10, 64)
    if err != nil {
        log.Println(err)
        size = 0
    }

    MonitoringDataStruct.TotalSize += size
}

func SendMonitoringDataToUI() {
    MonitoringDataStruct.Mux.Lock()
    defer MonitoringDataStruct.Mux.Unlock()

    MonitoringData_chan <- MonitoringDataStruct
}