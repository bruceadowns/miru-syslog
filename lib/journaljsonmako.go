package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// JournalJSONMako struct
type JournalJSONMako struct {
	journalJSON JournalJSON
	makoJSON    MakoJSON
}

// Extract ...
func (p JournalJSONMako) Extract(hn string, bb *bytes.Buffer) (res map[string]string, err error) {
	replacer := strings.NewReplacer(
		"\"level\":10,", "\"level\":\"TRACE\",",
		"\"level\":20,", "\"level\":\"DEBUG\",",
		"\"level\":30,", "\"level\":\"INFO\",",
		"\"level\":40,", "\"level\":\"WARN\",",
		"\"level\":50,", "\"level\":\"ERROR\",",
		"\"level\":60,", "\"level\":\"ERROR\",",
		"\"@timestamp\"", "\"timestamp\"",
		"\"@version\"", "\"version\"")

	jj := replacer.Replace(bb.String())
	jj = reVersionStrung.ReplaceAllString(jj, "\"version\":0")
	if err = json.NewDecoder(bytes.NewBufferString(jj)).Decode(&p.journalJSON); err != nil {
		return
	}

	if len(p.journalJSON.HostName) == 0 {
		return nil, fmt.Errorf("Host name not found")
	}

	level := "INFO"
	switch p.journalJSON.Priority {
	case "0", "1", "2", "3":
		level = "ERROR"
	case "4":
		level = "WARN"
	case "7":
		level = "DEBUG"
	//case "5", "6":
	default:
		level = "INFO"
	}

	timestamp := p.journalJSON.SourceRealtimeTimestamp
	if len(timestamp) == 0 {
		timestamp = p.journalJSON.RealtimeTimestamp
	}
	if len(timestamp) == 16 {
		timestamp = timestamp[:13]
	}

	hostname := p.journalJSON.HostName
	loggerName := p.journalJSON.Transport
	message := p.journalJSON.Message
	serviceName := p.journalJSON.Exe
	serviceVersion := p.journalJSON.PID
	serviceEnvironment := ""
	servicePipeline := ""
	stackTrace := ""
	threadName := ""
	version := ""

	mj := replacer.Replace(p.journalJSON.Message)
	mj = reVersionStrung.ReplaceAllString(mj, "\"version\":0")
	if err := json.NewDecoder(bytes.NewBufferString(mj)).Decode(&p.makoJSON); err == nil {
		loggerName = p.makoJSON.LoggerName
		level = p.makoJSON.Level
		message = p.makoJSON.Message
		serviceName = p.makoJSON.ServiceName
		serviceVersion = p.makoJSON.ServiceVersion
		serviceEnvironment = p.makoJSON.ServiceEnvironment
		servicePipeline = p.makoJSON.ServicePipeline
		stackTrace = p.makoJSON.StackTrace
		threadName = p.makoJSON.ThreadName
		version = strconv.Itoa(p.makoJSON.Version)
	}

	return map[string]string{
		"hostname":            hostname,
		"logger_name":         loggerName,
		"level":               level,
		"message":             message,
		"service_environment": serviceEnvironment,
		"service_name":        serviceName,
		"service_pipeline":    servicePipeline,
		"service_version":     serviceVersion,
		"stack_trace":         stackTrace,
		"thread_name":         threadName,
		"timestamp":           timestamp,
		"version":             version,
	}, nil
}

// Name ...
func (p JournalJSONMako) Name() string {
	return "journaljsonmako"
}
