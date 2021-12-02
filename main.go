package runlogger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

type severety string

const (
	default_severety   severety = "DEFAULT"   // The log entry has no assigned severity level.
	debug_severety     severety = "DEBUG"     // Debug or trace information.
	info_severety      severety = "INFO"      // Routine information, such as ongoing status or performance.
	notice_severety    severety = "NOTICE"    // Normal but significant events, such as start up, shut down, or a configuration change.
	warning_severety   severety = "WARNING"   // Warning events might cause problems.
	error_severety     severety = "ERROR"     // Error events are likely to cause problems.
	critical_severety  severety = "CRITICAL"  // Critical events cause more severe problems or outages.
	alert_severety     severety = "ALERT"     // A person must take an action immediately.
	emergency_severety severety = "EMERGENCY" // One or more systems are unusable.)
)

type Logger struct{}

// StructuredLogger is used to have structured logging in stackdriver (Google Cloud Platform)
func StructuredLogger() *Logger {
	return &Logger{}
}

// PlainLogger is used when you are not in a cloud run environment
func PlainLogger() *Logger {
	return nil
}

func (l *Logger) Debug(v ...interface{}) {
	l.writeLog(debug_severety, fmt.Sprint(v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.writeLog(info_severety, fmt.Sprint(v...))
}

func (l *Logger) Notice(v ...interface{}) {
	l.writeLog(notice_severety, fmt.Sprint(v...))
}

func (l *Logger) Warning(v ...interface{}) {
	l.writeLog(warning_severety, fmt.Sprint(v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.writeLog(error_severety, fmt.Sprint(v...))
}

func (l *Logger) Critical(v ...interface{}) {
	l.writeLog(critical_severety, fmt.Sprint(v...))
}

func (l *Logger) Alert(v ...interface{}) {
	l.writeLog(alert_severety, fmt.Sprint(v...))
}

func (l *Logger) Emergency(v ...interface{}) {
	l.writeLog(emergency_severety, fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.writeLog(debug_severety, fmt.Sprintf(format, v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.writeLog(info_severety, fmt.Sprintf(format, v...))
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.writeLog(notice_severety, fmt.Sprintf(format, v...))
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.writeLog(warning_severety, fmt.Sprintf(format, v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.writeLog(error_severety, fmt.Sprintf(format, v...))
}

func (l *Logger) Criticalf(format string, v ...interface{}) {
	l.writeLog(critical_severety, fmt.Sprintf(format, v...))
}

func (l *Logger) Alertf(format string, v ...interface{}) {
	l.writeLog(alert_severety, fmt.Sprintf(format, v...))
}

func (l *Logger) Emergencyf(format string, v ...interface{}) {
	l.writeLog(emergency_severety, fmt.Sprintf(format, v...))
}

func (l *Logger) Debugj(obj interface{}) {
	l.writeJson(debug_severety, l.objToJson(obj))
}

func (l *Logger) Infoj(obj interface{}) {
	l.writeJson(info_severety, l.objToJson(obj))
}

func (l *Logger) Noticej(obj interface{}) {
	l.writeJson(notice_severety, l.objToJson(obj))
}

func (l *Logger) Warningj(obj interface{}) {
	l.writeJson(warning_severety, l.objToJson(obj))
}

func (l *Logger) Errorj(obj interface{}) {
	l.writeJson(error_severety, l.objToJson(obj))
}

func (l *Logger) Criticalj(obj interface{}) {
	l.writeJson(critical_severety, l.objToJson(obj))
}

func (l *Logger) Alertj(obj interface{}) {
	l.writeJson(alert_severety, l.objToJson(obj))
}

func (l *Logger) Emergencyj(obj interface{}) {
	l.writeJson(emergency_severety, l.objToJson(obj))
}

func (l *Logger) writeLog(severety severety, message string) {
	pc, file, line, _ := runtime.Caller(2)
	if l == nil {
		fmt.Printf("%s in [%s:%d]: %s\n", severety, file, line, message)
		return
	}

	payload := &stackdriverLogStruct{
		TextPayload: message,
		Severity:    severety,
		Timestamp:   time.Now(),
		SourceLocation: &sourceLocation{
			File:     file,
			Function: runtime.FuncForPC(pc).Name(),
			Line:     strconv.Itoa(line),
		},
	}
	err := json.NewEncoder(os.Stdout).Encode(payload)
	if err != nil {
		panic("could not log TextPayload because of err: " + err.Error())
	}
}

func (l *Logger) writeJson(severety severety, jsonBytes []byte) {
	pc, file, line, _ := runtime.Caller(2)
	if l == nil {
		fmt.Printf("%s in [%s:%d]:\n%s\n", severety, file, line, string(jsonBytes))
		return
	}

	payload := &stackdriverLogStruct{
		JsonPayload: string(jsonBytes),
		Severity:    severety,
		Timestamp:   time.Now(),
		SourceLocation: &sourceLocation{
			File:     file,
			Function: runtime.FuncForPC(pc).Name(),
			Line:     strconv.Itoa(line),
		},
	}
	err := json.NewEncoder(os.Stdout).Encode(payload)
	if err != nil {
		panic("could not log JsonPayload because of err: " + err.Error())
	}
}

func (l *Logger) objToJson(obj interface{}) []byte {
	b, err := json.Marshal(obj)
	if err != nil {
		l.Criticalf("got error %v when marshaling %+v\n", err, obj)
		return []byte{}
	}
	return b
}

// stackdriverLogStruct source https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
type stackdriverLogStruct struct {
	TextPayload    string          `json:"message,omitempty"`
	JsonPayload    string          `json:"jsonPayload,omitempty"`
	Severity       severety        `json:"severity"`
	Timestamp      time.Time       `json:"timestamp`
	SourceLocation *sourceLocation `json:"sourceLocation,omitempty"`
}
type sourceLocation struct {
	File     string `json:"file"`
	Line     string `json:"line"`
	Function string `json:"function"`
}
