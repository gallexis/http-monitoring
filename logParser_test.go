package main

import (
    "testing"
)

func Test_GetSection(t *testing.T) {

    s := LogLine{Resource: "/pages/foo/bar?k=v"}
    if s.GetSection() != "/pages" {
        t.Error("Error in GetSection().")
        return
    }

}

// Test each field of a commonLog struct
func Test_ParseCommonLog(t *testing.T) {
    line := "127.0.0.1 user-identifier frank [wrong date] \"GET /apache_pb.gif HTTP/1.0\" 200 2326"
    commonLog, err := ParseLine(line)
    if err != nil {
        t.Error("ParseLine() shouldn't throw an error: ", err)
        return
    }

    if commonLog.RemoteHostIP != "127.0.0.1" ||
        commonLog.Identd != "user-identifier" ||
        commonLog.UserID != "frank" ||
        commonLog.Date.String() != "0001-01-01 00:00:00 +0000 UTC" ||
        commonLog.Method != "GET" ||
        commonLog.Resource != "/apache_pb.gif" ||
        commonLog.Protocol != "HTTP/1.0" ||
        commonLog.Status != "200" ||
        commonLog.Size != "2326" {
        t.Error("At least one field is bad decoded")
        t.Error(commonLog)
        return
    }

}
