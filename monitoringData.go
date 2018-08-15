package main

import (
    "fmt"
    "log"
    "sort"
    "strconv"
    "sync"
    "container/ring"
)

/*
   Structure representing the metrics to be monitored
   Used between many goroutines so don't forget the mutex
*/
type MonitoringData struct {
    HttpRequestsSize      uint64
    HttpRequestsCount     uint64
    HttpStatuses          map[string]uint64
    UrlSections           map[string]uint64
    Alerts                *ring.Ring
    LogLines              *ring.Ring
    IsTrafficAlert        bool
    LastAlertMessage      string
    PreviousTotalRequests uint64
    RequestsThreshold     uint64
    Mux                   sync.Mutex
}

func NewMonitoringData() MonitoringData {
    return MonitoringData{
        UrlSections:       make(map[string]uint64),
        HttpStatuses:      make(map[string]uint64),
        Alerts:            ring.New(5),
        LogLines:          ring.New(10),
        IsTrafficAlert:    false,
        RequestsThreshold: 25,
    }
}


/*
    The ring is just a circular array.
    ringToStringArray allows to put all the elements of the ring into an array of string
 */
func (md *MonitoringData) ringToStringArray(r *ring.Ring) []string {
    var str []string

    r.Do(func(p interface{}) {
        if p != nil {
            str = append(str, p.(string))
        }
    })
    return str
}

// Update the metric structs from a LogLine struct
func (md *MonitoringData) Update(logLine LogLine) {
    md.Mux.Lock()
    defer md.Mux.Unlock()

    section := logLine.GetSection()

    md.UrlSections[section] += 1
    md.HttpStatuses[logLine.Status] += 1
    md.HttpRequestsCount += 1

    // Parse a size in String to uint64
    size, err := strconv.ParseUint(logLine.Size, 10, 64)
    if err != nil {
        log.Println(err.Error())
    }
    md.HttpRequestsSize += size
}

func (md MonitoringData) Display() []string {
    return []string{
        fmt.Sprintf("Total HTTP requests : %d", md.HttpRequestsCount),
        fmt.Sprintf("Total Size emitted in bytes : %d", md.HttpRequestsSize),
        getHttpStatus(md.HttpStatuses, "HTTP Status : %s", 5),
        getMostViewedSections(md.UrlSections, "Most viewed sections : %s", 3),
    }
}

/*
   This part is used for ordering the values of the map[string]uint64 in descending order,
    and then display the most viewed section and the HTTP status
*/
type Pair struct {
    Key   string
    Value uint64
}

// sortMap will return a new []Pair slice sorted by descending order of Value
func sortMap(m map[string]uint64) []Pair {
    /*
       When sorting the map, we first put each of its key/value in a Pair structure.
       Each of these Pair structures will be put in an array, then sorted by Value
    */
    var array []Pair

    for k, v := range m {
        array = append(array, Pair{k, v})
    }

    sort.Slice(array, func(i, j int) bool {
        return array[i].Value > array[j].Value
    })

    return array
}

// getMostViewedSections returns a string displaying the most viewed sections of the site,
// ordered from the highest to the lower
func getMostViewedSections(URLsections map[string]uint64, mostViewedSections string, limit int) string {
    var buffer string
    orderedArray := sortMap(URLsections)

    for i, kv := range orderedArray {
        if i >= limit {
            break
        }

        // This will append [key : value] to the buffer
        buffer = fmt.Sprintf("%s[%s : %d] ", buffer, kv.Key, kv.Value)
    }

    // This trims the last white space and returns the buffer
    // using the mostViewedSections format.
    return fmt.Sprintf(mostViewedSections, buffer)
}

// getHttpStatus returns a string displaying the HTTP status of the site,
// ordered from the highest to the lower
func getHttpStatus(HttpStatusMap map[string]uint64, HttpStatus string, limit int) string {
    var buffer string
    orderedArray := sortMap(HttpStatusMap)

    for i, kv := range orderedArray {
        if i >= limit {
            break
        }

        // This will append [key : value] to the buffer
        buffer = fmt.Sprintf("%s[%s : %d] ", buffer, kv.Key, kv.Value)
    }

    // This trims the last white space and returns the buffer
    // using the HttpStatusMap format.
    return fmt.Sprintf(HttpStatus, buffer)
}
