package webshell

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	logdate_fmt   = "20060102"
	timestamp_fmt = "2006-01-02T15:04:05Z"
)

var (
	ipScrub        = regexp.MustCompile("^([^:]):.*$")
	access_logfile = "logs/access"
)

type DefaultRequestLogger struct {
	AccessLog       string
	ErrorLog        string
	TimestampFormat string
}

type DTSRequestLogger DefaultRequestLogger

type RequestLogger interface {
	GetAccessLog() string
	GetTimestampFormat() string
	LogfileName(string) string
}

type RequestLogItem struct {
	ClientIP  string
	Timestamp string
	Path      string
	Status    int
	Logfile   string
}

func NewDefaultRequestLogger(accesslog, errorlog string) *DefaultRequestLogger {
	logger := new(DefaultRequestLogger)
	if accesslog == "" {
		logger.AccessLog = access_logfile + ".log"
	} else {
		logger.AccessLog = accesslog
	}
	if errorlog == "" {
		logger.ErrorLog = logger.AccessLog
	} else {
		logger.ErrorLog = errorlog
	}
	logger.TimestampFormat = timestamp_fmt
	return logger
}

func (l *DefaultRequestLogger) GetTimestampFormat() string {
	return l.TimestampFormat
}

func (l *DefaultRequestLogger) GetAccessLog() string {
	return l.AccessLog
}

func NewDTSRequestLogger(accesslog, errorlog string) *DTSRequestLogger {
	logger := new(DTSRequestLogger)
	if accesslog == "" {
		logger.AccessLog = access_logfile
	} else {
		logger.AccessLog = accesslog
	}
	if errorlog == "" {
		logger.ErrorLog = logger.AccessLog
	} else {
		logger.ErrorLog = errorlog
	}
	logger.TimestampFormat = timestamp_fmt
	return logger
}

func (l *DTSRequestLogger) GetTimestampFormat() string {
	return l.TimestampFormat
}

func (l *DTSRequestLogger) GetAccessLog() string {
	return l.AccessLog
}

func GetTimestamp(l RequestLogger) string {
	return time.Now().Format(l.GetTimestampFormat())
}

func GetClientIp(r *http.Request) string {
	if r.Header.Get("X-Real-Ip") == "" {
		return r.RemoteAddr
	} else {
		return r.Header.Get("X-Real-Ip")
	}
	return r.RemoteAddr
}

func NewRequestLogItem(l RequestLogger, r *http.Request) *RequestLogItem {
	req := new(RequestLogItem)
	req.ClientIP = GetClientIp(r)
	req.Timestamp = GetTimestamp(l)
	req.Path = r.URL.String()

	return req
}

func LogRequest(l RequestLogger, r *http.Request) {
	req := NewRequestLogItem(l, r)
	req.LogRequest(l)
}

func (req *RequestLogItem) LogRequest(l RequestLogger) {
	req.Logfile = l.LogfileName(l.GetAccessLog())
	log_line := fmt.Sprintf("%s %s %s %d\n", req.ClientIP,
		req.Timestamp, req.Path, req.Status)
	fmt.Printf(log_line)
	if err := writeLogEntry(access_logfile, log_line); err != nil {
		fmt.Printf("[!] error writing to %s: %s\n", req.Logfile,
			err.Error())
	}
}

func writeLogEntry(logfile, line string) (err error) {
	file, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(logfile)
	}

	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(line)
	return
}

func (l *DTSRequestLogger) LogfileName(logfile string) string {
	tmpname := fmt.Sprintf("%s-%s.log", logfile,
		time.Now().Format(logdate_fmt))
	name, err := filepath.Abs(tmpname)
	if err != nil {
		name = filepath.Base(tmpname)
	}
	return name
}

func (l *DefaultRequestLogger) LogfileName(logfile string) string {
	return logfile
}
