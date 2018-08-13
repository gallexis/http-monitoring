package main

import (
	"testing"
    "time"
)

func Test_GetSection(t *testing.T) {

    s := LogLine{Resource: "/pages/foo/bar?k=v"}
    if s.GetSection() != "/pages"{
        t.Error("Error in GetSection().")
        return
    }

}


// Test each field of a commonLog struct
func Test_ParseCommonLog(t *testing.T) {
    line := "127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326"
    commonLog, err := ParseLine(line)
    if err != nil{
        t.Error("ParseLine() shouldn't have thrown an error:", err)
        return
    }

    date, err := time.Parse("02/Jan/2006:15:04:05 -0700", "10/Oct/2000:13:55:36 -0700")
    if err != nil {
        t.Error("time.Parse threw an error:", err)
        return
    }

    if  commonLog.RemoteHostIP != "127.0.0.1" ||
        commonLog.Identd != "user-identifier" ||
        commonLog.UserID != "frank" ||
        commonLog.Date.String() != date.String() ||
        commonLog.Method != "GET" ||
        commonLog.Resource != "/apache_pb.gif" ||
        commonLog.Protocol != "HTTP/1.0" ||
        commonLog.Status != "200"||
        commonLog.Size != "2326" {

        t.Error("At least one field is bad decoded")
        t.Error(commonLog)
        return
    }

}




