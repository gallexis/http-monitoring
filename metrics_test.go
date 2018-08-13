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
    cl.UserID = "frank"
    cl.Date = date
    cl.Method = "GET"
    cl.Resource = "/apache_pb.gif"
    cl.Protocol = "HTTP/1.0"
    cl.Status = "200"
    cl.Size = "2326"

    return cl
}

func Test_Update(t *testing.T) {
    m := NewMetrics()

    // we update 'm' 2 times (with the same LogLine)
    m.Update(commonLogMock())
    m.Update(commonLogMock())

    if v, ok := m.UrlSections["/apache_pb.gif"]; !ok || v != 2{
        t.Error("UrlSections has not been updated correctly")
        return
    }

    if v, ok := m.HttpStatuses["200"]; !ok || v != 2{
        t.Error("UrlSections has not been updated correctly")
        return
    }

    if m.HttpRequests != 2{
        t.Error("HttpRequests should be 2")
        return
    }

    if m.TotalSize != 2326 * 2{
        t.Error("TotalSize should be 4652 (2326 * 2)")
        return
    }
}

// Test if sortMap sorts the values of the map in descending order
func Test_orderMapByValue(t *testing.T){
    m := map[string]int{
        "r" : -1,
        "b" : 1,
        "a" : 0,
    }
    array := sortMap(m)

    if  array[0] != (Pair{"b", 1})  ||
        array[1] != (Pair{"a",  0}) ||
        array[2] != (Pair{"r",  -1}) {

        t.Error("Array of Pair not sorted correctly in descending order")
        return
    }
}