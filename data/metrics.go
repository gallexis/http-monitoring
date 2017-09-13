package data

import (
    "fmt"
    "http-monitoring/parser"
    "log"
    "sort"
    "strconv"
    "strings"
    "sync"
)

/*
   Structure representing the metrics to be monitored
   Used between many goroutines, so don't forget the mutex (Mux)
*/
type MetricStruct struct {
    TotalSize     uint64
    TotalRequests uint64
    URL_sections  map[string]int
    HTTP_status   map[string]int
    Mux           sync.Mutex
}

func NewMetricStruct() MetricStruct {
    return MetricStruct{
        URL_sections: make(map[string]int),
        HTTP_status:  make(map[string]int),
    }
}

// Update the metricStruct with a commonLog struct
func (d *MetricStruct) Update(logStr parser.CommonLogStruct) {
    d.Mux.Lock()
    defer d.Mux.Unlock()

    section := logStr.GetSection()

    // URL section counter (i.e : /pages : 4 )
    d.URL_sections[section] += 1

    // HTTP status counter (i.e : 200, 404, 500...)
    d.HTTP_status[logStr.Status] += 1

    // Increase counter of HTTP requests
    d.TotalRequests += 1

    // Parse the Size from a String to an Uint64
    size, err := strconv.ParseUint(logStr.Size, 10, 64)
    if err != nil {
        log.Panic(err)
    }
    d.TotalSize += size
}

func (d MetricStruct) Display() []string {
    // String Formats
    totalHTTPRequests := "Total HTTP requests : %d"
    totalSizeEmitted := "Total Size emitted in bytes : %d"
    mostViewedSections := "Most viewed sections : %s"
    HTTPstatus := "HTTP Status : %s"

    return []string{
        fmt.Sprintf(totalHTTPRequests, d.TotalRequests),
        fmt.Sprintf(totalSizeEmitted, d.TotalSize),
        getHTTPstatus(d.HTTP_status, HTTPstatus),
        getMostViewedSections(d.URL_sections, mostViewedSections),
    }
}

/*
   This part is used for ordering the values of the map[string]int in descending order
*/

type Pair struct {
    Key   string
    Value int
}

// orderMapByValue will return a new []Pair slice sorted by Value
func orderMapByValue(m map[string]int) []Pair {
    /*
       When ordering the map, we first put each of its key/value in a Pair structure.
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

// getMostViewedSections a string displaying the most viewed sections of the site,
// ordered from the highest to the lower
func getMostViewedSections(URLsections map[string]int, mostViewedSections string) string {
    var buffer string
    orderedArray := orderMapByValue(URLsections)

    for i, kv := range orderedArray {
        if i > 2 {
            break
        }

        // This will append [key : value] to the buffer
        buffer = fmt.Sprintf("%s[%s : %d] ", buffer, kv.Key, kv.Value)
    }

    // This trims the last white space and returns the buffer
    // using the MostViewedSections format.
    return fmt.Sprintf(mostViewedSections, strings.TrimRight(buffer, " "))
}

// getHTTPstatus returns a string displaying the HTTP status of the site,
// ordered from the highest to the lower
func getHTTPstatus(HTTPstatus map[string]int, HTTPStatus string) string {
    var buffer string
    orderedArray := orderMapByValue(HTTPstatus)

    for i, kv := range orderedArray {
        if i > 4 {
            break
        }

        // This will append [key : value] to the buffer
        buffer = fmt.Sprintf("%s[%s : %d] ", buffer, kv.Key, kv.Value)
    }

    // This trims the last white space and returns the buffer
    // using the HTTPStatus format.
    return fmt.Sprintf(HTTPStatus, strings.TrimRight(buffer, " "))
}
