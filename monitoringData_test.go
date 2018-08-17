package main

import (
    "testing"
    "time"
    "log"
)

// commonLogMock generates a fake commonLog struct
func commonLogMock() LogLine {
    cl := LogLine{}

    date, err := time.Parse("02/Jan/2006:15:04:05 -0700", "10/Oct/2000:13:55:36 -0700")
    if err != nil {
        log.Panic(err.Error())
    }

    cl.RemoteHostIP = "127.0.0.1"
    cl.Identd = "user-identifier"
    cl.UserID = "alexis"
    cl.Date = date
    cl.Method = "GET"
    cl.Resource = "/apache_pb.gif"
    cl.Protocol = "HTTP/1.0"
    cl.Status = "200"
    cl.Size = "2326"

    return cl
}

func Test_Update(t *testing.T) {
    m := NewMonitoringData(25)

    // we update 'm' 2 times (with the same log line)
    m.Update(commonLogMock())
    m.Update(commonLogMock())

    if v, ok := m.UrlSections["/apache_pb.gif"]; !ok || v != 2 {
        t.Error("UrlSections has not been updated correctly")
        return
    }

    if v, ok := m.HttpStatuses["200"]; !ok || v != 2 {
        t.Error("UrlSections has not been updated correctly")
        return
    }

    if m.HttpRequestsCount != 2 {
        t.Error("HttpRequestsCount should be 2")
        return
    }

    if m.HttpRequestsSize != 2326*2 {
        t.Error("HttpRequestsSize should be 4652 (2326 * 2)")
        return
    }
}

// Test if sortMap sorts the values of the map in descending order
func Test_orderMapByValue(t *testing.T) {
    m := map[string]uint64{
        "r": 0,
        "b": 2,
        "a": 1,
    }
    array := sortMap(m)

    if array[0] != (Pair{"b", 2}) ||
        array[1] != (Pair{"a", 1}) ||
        array[2] != (Pair{"r", 0}) {

        t.Error("Array of Pair not sorted correctly in descending order")
        return
    }
}

// Test if sortMap sorts the values of the map in descending order
func Test_DisplayMonitoringData(t *testing.T) {
    m := NewMonitoringData(25)

    // we update 'm' 2 times (with the same log line)
    m.Update(commonLogMock())
    m.Update(commonLogMock())

    display := m.Display()

    if display[0] != "Total HTTP requests : 2" {
        t.Error("Wrong Total HTTP requests count", m.HttpRequestsCount, display[0])
        return
    }

    if display[1] != "Total Size emitted in bytes : 4652" {
        t.Error("Wrong Total Size emitted in bytes", m.HttpRequestsSize, display[1])
        return
    }

    if display[2] != "HTTP Status : [200 : 2] " {
        t.Error("Wrong HTTP statuses", m.HttpStatuses, display[2])
        return
    }

    if display[3] != "Most viewed sections : [/apache_pb.gif : 2] " {
        t.Error("Wrong Total of most viewed sections", m.UrlSections, display[3])
        return
    }
}
