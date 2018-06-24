// Copyright (c) Jeevanandam M (https://github.com/jeevatkm)
// go-aah/log source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"aahframework.org/essentials.v0"
)

// Format flags used to define log message format for each log entry
const (
	FmtFlagLevel ess.FmtFlag = iota
	FmtFlagAppName
	FmtFlagInstanceName
	FmtFlagRequestID
	FmtFlagPrincipal
	FmtFlagTime
	FmtFlagUTCTime
	FmtFlagLongfile
	FmtFlagShortfile
	FmtFlagLine
	FmtFlagMessage
	FmtFlagFields
	FmtFlagCustom
	FmtFlagUnknown
)

const (
	textFmt = "text"
	jsonFmt = "json"
	space   = " "
)

type (
	// FlagPart is indiviual flag details
	//  For e.g.:
	//    part := FlagPart{
	//      Flag:   fmtFlagTime,
	//      Name:   "time",
	//      Format: "2006-01-02 15:04:05.000",
	//    }
	FlagPart struct {
		Flag   ess.FmtFlag
		Name   string
		Format string
	}
)

var (
	// DefaultPattern is default log entry pattern in aah/log. Only applicable to
	// text formatter.
	// For e.g:
	//    2006-01-02 15:04:05.000 INFO  This is my message
	DefaultPattern = "%time:2006-01-02 15:04:05.000 %level:-5 %message"

	// FmtFlags is the list of log format flags supported by aah log library
	// Usage of flag order is up to format composition.
	//    level     - outputs ERROR, WARN, INFO, DEBUG, TRACE
	//    appname   - outputs Application Name
	//    insname   - outputs Application Instance Name
	//    reqid     - outputs Request ID from HTTP header
	//    principal - outputs Logged-In subject primary principal value
	//    time      - outputs local time as per format supplied
	//    utctime   - outputs UTC time as per format supplied
	//    longfile  - outputs full file name: /a/b/c/d.go
	//    shortfile - outputs final file name element: d.go
	//    line      - outputs file line number: L23
	//    message   - outputs given message along supplied arguments if they present
	//    fields    - outputs field values into log entry
	//    custom    - outputs string as-is into log entry
	FmtFlags = map[string]ess.FmtFlag{
		"level":     FmtFlagLevel,
		"appname":   FmtFlagAppName,
		"insname":   FmtFlagInstanceName,
		"reqid":     FmtFlagRequestID,
		"principal": FmtFlagPrincipal,
		"time":      FmtFlagTime,
		"utctime":   FmtFlagUTCTime,
		"longfile":  FmtFlagLongfile,
		"shortfile": FmtFlagShortfile,
		"line":      FmtFlagLine,
		"message":   FmtFlagMessage,
		"fields":    FmtFlagFields,
		"custom":    FmtFlagCustom,
	}
)

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// textFormatter
//___________________________________

// textFormatter formats the `Entry` object details as per log `pattern`
// 	For e.g.:
// 		2016-07-02 22:26:01.530 INFO formatter_test.go L29 - Yes, I would love to see
func textFormatter(flags []ess.FmtFlagPart, entry *Entry) []byte {
	buf := new(bytes.Buffer)

	for _, part := range flags {
		switch part.Flag {
		case FmtFlagLevel:
			buf.WriteString(fmt.Sprintf(part.Format, entry.Level) + space)
		case FmtFlagAppName:
			if len(entry.AppName) > 0 {
				buf.WriteString(entry.AppName + space)
			}
		case FmtFlagInstanceName:
			if len(entry.InstanceName) > 0 {
				buf.WriteString(entry.InstanceName + space)
			}
		case FmtFlagRequestID:
			if len(entry.RequestID) > 0 {
				buf.WriteString(entry.RequestID + space)
			}
		case FmtFlagPrincipal:
			if len(entry.Principal) > 0 {
				buf.WriteString(entry.Principal + space)
			}
		case FmtFlagTime:
			buf.WriteString(entry.Time.Format(part.Format) + space)
		case FmtFlagUTCTime:
			buf.WriteString(entry.Time.UTC().Format(part.Format) + space)
		case FmtFlagLongfile, FmtFlagShortfile:
			if part.Flag == FmtFlagShortfile {
				entry.File = filepath.Base(entry.File)
			}
			buf.WriteString(fmt.Sprintf(part.Format, entry.File) + space)
		case FmtFlagLine:
			buf.WriteString("L" + fmt.Sprintf(part.Format, entry.Line) + space)
		case FmtFlagMessage:
			buf.WriteString(entry.Message + space)
		case FmtFlagCustom:
			buf.WriteString(part.Format + space)
		case FmtFlagFields:
			fs := make([]string, 0)
			for k, v := range entry.Fields {
				if !entry.isSkipField(k) {
					fs = append(fs, fmt.Sprintf("%v: %v", k, v))
				}
			}

			if len(fs) > 0 {
				buf.WriteString("fields[" + strings.Join(fs, ", ") + "] ")
			}
		}
	}

	buf.WriteByte('\n')
	return buf.Bytes()
}
