package runlogger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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

const maxSize = 102400

var errorMessageType = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"

var stdout = bufio.NewWriter(os.Stdout)

type Logger struct{}

var prefixPath string

// StructuredLogger is used to have structured logging in stackdriver (Google Cloud Platform)
func StructuredLogger() *Logger {
	setPrefixPath()
	return &Logger{}
}

// PlainLogger is used when you are not in a cloud run environment
func PlainLogger() *Logger {
	setPrefixPath()
	return nil
}

func setPrefixPath() {
	_, fileName, _, _ := runtime.Caller(2)
	prefixPath = filepath.Dir(fileName) + "/"
}

func (l *Logger) Debug(v ...interface{}) {
	l.writeLog(debug_severety, strings.TrimSpace(fmt.Sprintln(v...)), nil)
}

func (l *Logger) Info(v ...interface{}) {
	l.writeLog(info_severety, strings.TrimSpace(fmt.Sprintln(v...)), nil)
}

func (l *Logger) Notice(v ...interface{}) {
	l.writeLog(notice_severety, strings.TrimSpace(fmt.Sprintln(v...)), nil)
}

func (l *Logger) Warning(v ...interface{}) {
	l.writeLog(warning_severety, strings.TrimSpace(fmt.Sprintln(v...)), nil)
}

func (l *Logger) Error(v ...interface{}) {
	l.writeLog(error_severety, strings.TrimSpace(fmt.Sprintln(v...)), nil)
}

func (l *Logger) Critical(v ...interface{}) {
	l.writeLog(critical_severety, strings.TrimSpace(fmt.Sprintln(v...)), nil)
}

func (l *Logger) Alert(v ...interface{}) {
	l.writeLog(alert_severety, strings.TrimSpace(fmt.Sprintln(v...)), nil)
}

func (l *Logger) Emergency(v ...interface{}) {
	l.writeLog(emergency_severety, strings.TrimSpace(fmt.Sprintln(v...)), nil)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.writeLog(debug_severety, fmt.Sprintf(format, v...), nil)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.writeLog(info_severety, fmt.Sprintf(format, v...), nil)
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.writeLog(notice_severety, fmt.Sprintf(format, v...), nil)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.writeLog(warning_severety, fmt.Sprintf(format, v...), nil)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.writeLog(error_severety, fmt.Sprintf(format, v...), nil)
}

func (l *Logger) Criticalf(format string, v ...interface{}) {
	l.writeLog(critical_severety, fmt.Sprintf(format, v...), nil)
}

func (l *Logger) Alertf(format string, v ...interface{}) {
	l.writeLog(alert_severety, fmt.Sprintf(format, v...), nil)
}

func (l *Logger) Emergencyf(format string, v ...interface{}) {
	l.writeLog(emergency_severety, fmt.Sprintf(format, v...), nil)
}

func (l *Logger) Debugj(message string, obj interface{}) {
	l.writeLog(debug_severety, message, obj)
}

func (l *Logger) Infoj(message string, obj interface{}) {
	l.writeLog(info_severety, message, obj)
}

func (l *Logger) Noticej(message string, obj interface{}) {
	l.writeLog(notice_severety, message, obj)
}

func (l *Logger) Warningj(message string, obj interface{}) {
	l.writeLog(warning_severety, message, obj)
}

func (l *Logger) Errorj(message string, obj interface{}) {
	l.writeLog(error_severety, message, obj)
}

func (l *Logger) Criticalj(message string, obj interface{}) {
	l.writeLog(critical_severety, message, obj)
}

func (l *Logger) Alertj(message string, obj interface{}) {
	l.writeLog(alert_severety, message, obj)
}

func (l *Logger) Emergencyj(message string, obj interface{}) {
	l.writeLog(emergency_severety, message, obj)
}

func (l *Logger) writeLog(severety severety, message string, obj interface{}) {
	output := os.Stderr

	var isError bool
	switch severety {
	case error_severety, critical_severety, alert_severety, emergency_severety:
		isError = true
	}

	pc, file, line, _ := runtime.Caller(2)

	if l == nil {
		if obj != nil {
			j, _ := json.Marshal(obj)
			fmt.Fprintf(
				output,
				"%s in [%s:%d]: %s\n%s\n",
				severety,
				relative(file),
				line,
				message,
				j,
			)
		} else {
			fmt.Fprintf(
				output,
				"%s in [%s:%d]: %s\n",
				severety,
				relative(file),
				line,
				message,
			)
		}
		return
	}

	var (
		messageType    *string
		serviceContext *ServiceContext
	)
	if isError {
		output = os.Stderr
		messageType = &errorMessageType
	}
	if os.Getenv("K_SERVICE") != "" {
		serviceContext = &ServiceContext{
			Service: os.Getenv("K_SERVICE"),
		}
	}

	payload := &stackdriverLogStruct{
		JsonPayload: obj,
		Message:     message,
		Severity:    severety,
		Timestamp:   time.Now(),
		Type:        messageType,
		SourceLocation: &sourceLocation{
			File:     relative(file),
			Function: runtime.FuncForPC(pc).Name(),
			Line:     strconv.Itoa(line),
		},
		ServiceContext: serviceContext,
	}
	j, err := json.Marshal(payload)
	if err != nil {
		panic("could not log because of err: " + err.Error())
	}

	if len(j) >= maxSize {
		l.Errorf("log entry exeed max size of %d bytes: %.100000s", maxSize, j)
	} else {
		fmt.Fprintf(output, "%s\n", j)
	}
}

func relative(path string) string {
	if filepath.HasPrefix(path, prefixPath) {
		return path[len(prefixPath):]
	}
	return path
}

// stackdriverLogStruct source https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
type stackdriverLogStruct struct {
	Message        string          `json:"message"`
	JsonPayload    interface{}     `json:"jsonPayload,omitempty"`
	Severity       severety        `json:"severity"`
	Timestamp      time.Time       `json:"timestamp"`
	SourceLocation *sourceLocation `json:"logging.googleapis.com/sourceLocation"`
	Type           *string         `json:"@type,omitempty"`
	ServiceContext *ServiceContext `json:"serviceContext,omitempty"`
}
type ServiceContext struct {
	Service string `json:"service"`
}
type sourceLocation struct {
	File     string `json:"file"`
	Line     string `json:"line"`
	Function string `json:"function"`
}
