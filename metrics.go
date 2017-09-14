package main

import (
    "fmt"
    "log"
    "sort"
    "strconv"
    "strings"
    "sync"
)

/*
   Structure representing the metrics to be monitored
   Used between many goroutines so don't forget the mutex (Mux)
*/
type Metrics struct {
    TotalSize    uint64
    Requests     uint64
    URL_sections map[string]int
    HTTP_status  map[string]int
    Mux          sync.Mutex
}

func NewMetricStruct() Metrics {
    return Metrics{
        URL_sections: make(map[string]int),
        HTTP_status:  make(map[string]int),
    }
}

// Update the metrics struct from a commonLog struct
func (m *Metrics) Update(commonLog CommonLog) {
    m.Mux.Lock()
    defer m.Mux.Unlock()

    section := commonLog.GetSection()

    // URL section counter (i.e : /pages : 4 )
    m.URL_sections[section] += 1

    // HTTP status counter (i.e : 200, 404, 500...)
    m.HTTP_status[commonLog.Status] += 1

    // Increase counter of HTTP requests
    m.Requests += 1

    // Parse the Size from a String to an uint64
    size, err := strconv.ParseUint(commonLog.Size, 10, 64)
    if err != nil {
        log.Panic(err)
    }
    m.TotalSize += size
}

func (m Metrics) Display() []string {
    // String Formats
    totalHTTPRequests := "Total HTTP requests : %d"
    totalSizeEmitted := "Total Size emitted in bytes : %d"
    mostViewedSections := "Most viewed sections : %s"
    HTTPstatus := "HTTP Status : %s"

    return []string{
        fmt.Sprintf(totalHTTPRequests, m.Requests),
        fmt.Sprintf(totalSizeEmitted, m.TotalSize),
        getHTTPstatus(m.HTTP_status, HTTPstatus),
        getMostViewedSections(m.URL_sections, mostViewedSections),
    }
}





/*
   This part is used for ordering the values of the map[string]int in descending order,
    and then display the most viewed section and the HTTP status
*/

type Pair struct {
    Key   string
    Value int
}

// orderMapByValue will return a new []Pair slice sorted by descending order of Value
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

// getMostViewedSections returns a string displaying the most viewed sections of the site,
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
    // using the mostViewedSections format.
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
    // using the HTTPstatus format.
    return fmt.Sprintf(HTTPStatus, strings.TrimRight(buffer, " "))
}
