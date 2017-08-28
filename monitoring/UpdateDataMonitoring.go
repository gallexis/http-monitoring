package monitoring

import (
    "strconv"
)

func UpdateData(s LineStruct){
    MonitoringD.Mux.Lock()
    section := GetURLSection(s)

    // URL section counter (i.e : /pages : 4 )
    MonitoringD.MapSection[section]  += 1

    // HTTP status counter (i.e : 200, 404, 500...)
    MonitoringD.HttpStatus[s.Status] += 1

    // Increase number of requests (set at 0 every 2min)
    MonitoringD.TotalRequests   += 1

    size, err := strconv.ParseUint(s.Size, 10, 64)
    if err != nil {
        // log error
    } else{
        MonitoringD.TotalSize += size
    }
    MonitoringD.Mux.Unlock()
}
