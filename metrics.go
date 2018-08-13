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
   Used between many goroutines so don't forget the mutex
*/
type Metrics struct {
    TotalSize    uint64
    HttpRequests uint64
    HttpStatuses map[string]int
    UrlSections  map[string]int
    Mux          sync.Mutex
}

func NewMetrics() Metrics {
    return Metrics{
        UrlSections:  make(map[string]int),
        HttpStatuses: make(map[string]int),
    }
}

// Update the metric structs from a LogLine struct
func (m *Metrics) Update(logLine LogLine) {
    m.Mux.Lock()
    defer m.Mux.Unlock()

    section := logLine.GetSection()

    m.UrlSections[section] += 1
    m.HttpStatuses[logLine.Status] += 1
    m.HttpRequests += 1

    // Parse a size in String to uint64
    size, err := strconv.ParseUint(logLine.Size, 10, 64)
    if err != nil {
        log.Println(err.Error())
    }
    m.TotalSize += size
}

func (m Metrics) Display() []string {
    return []string{
        fmt.Sprintf("Total HTTP requests : %d", m.HttpRequests),
        fmt.Sprintf("Total Size emitted in bytes : %d", m.TotalSize),
        getHTTPstatus(m.HttpStatuses, "HTTP Status : %s", 5),
        getMostViewedSections(m.UrlSections, "Most viewed sections : %s", 3),
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

// sortMap will return a new []Pair slice sorted by descending order of Value
func sortMap(m map[string]int) []Pair {
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
func getMostViewedSections(URLsections map[string]int, mostViewedSections string, limit int) string {
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
    return fmt.Sprintf(mostViewedSections, strings.TrimRight(buffer, " "))
}

// getHTTPstatus returns a string displaying the HTTP status of the site,
// ordered from the highest to the lower
func getHTTPstatus(HTTPstatus map[string]int, HTTPStatus string, limit int) string {
    var buffer string
    orderedArray := sortMap(HTTPstatus)

    for i, kv := range orderedArray {
        if i >= limit {
            break
        }

        // This will append [key : value] to the buffer
        buffer = fmt.Sprintf("%s[%s : %d] ", buffer, kv.Key, kv.Value)
    }

    // This trims the last white space and returns the buffer
    // using the HTTPstatus format.
    return fmt.Sprintf(HTTPStatus, strings.TrimRight(buffer, " "))
}
