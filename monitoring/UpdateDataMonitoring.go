package monitoring

import (
    "strconv"
)

func UpdateMonitoringDataStruct(s LineStruct){
    MonitoringDataStruct.Mux.Lock()
    section := GetURLSection(s)

    // URL section counter (i.e : /pages : 4 )
    MonitoringDataStruct.MapURLsection[section]  += 1

    // HTTP status counter (i.e : 200, 404, 500...)
    MonitoringDataStruct.MapHTTPstatus[s.Status] += 1

    // Increase number of requests (set at 0 every 2min)
    MonitoringDataStruct.TotalRequests   += 1

    size, err := strconv.ParseUint(s.Size, 10, 64)
    if err != nil {
        // log error
    } else{
        MonitoringDataStruct.TotalSize += size
    }
    MonitoringDataStruct.Mux.Unlock()
}
