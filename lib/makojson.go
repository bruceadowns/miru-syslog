package lib

import (
	"bytes"
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MakoJSON holds mako structured json
type MakoJSON struct {
	LoggerName         string `json:"logger_name,omitempty"`
	Message            string `json:"message,omitempty"`
	Level              string `json:"level,omitempty"`
	ServiceEnvironment string `json:"service_environment,omitempty"`
	ServiceName        string `json:"service_name,omitempty"`
	ServicePipeline    string `json:"service_pipeline,omitempty"`
	ServiceVersion     string `json:"service_version,omitempty"`
	ThreadName         string `json:"thread_name,omitempty"`
	Timestamp          string `json:"timestamp,omitempty"`
	Version            int    `json:"version,omitempty"`
	StackTrace         string `json:"stack_trace,omitempty"`
}

// global const in order to compile once
var reVersionStrung = regexp.MustCompile("\"version\":\"[0-9.]+\"")

// Extract ...
func (p MakoJSON) Extract(hostname string, bb *bytes.Buffer) (res map[string]string, err error) {
	replacer := strings.NewReplacer(
		"\"level\":10,", "\"level\":\"TRACE\",",
		"\"level\":20,", "\"level\":\"DEBUG\",",
		"\"level\":30,", "\"level\":\"INFO\",",
		"\"level\":40,", "\"level\":\"WARN\",",
		"\"level\":50,", "\"level\":\"ERROR\",",
		"\"level\":60,", "\"level\":\"ERROR\",",
		"\"@timestamp\"", "\"timestamp\"",
		"\"@version\"", "\"version\"")

	out := replacer.Replace(bb.String())
	out = reVersionStrung.ReplaceAllString(out, "\"version\":0")

	err = json.NewDecoder(bytes.NewBufferString(out)).Decode(&p)
	if err != nil {
		return nil, err
	}

	timestamp := "0"
	if ts, err := time.Parse(time.RFC3339, p.Timestamp); err == nil {
		timestamp = strconv.FormatInt(ts.UnixNano()/1000000, 10)
	} else {
		log.Printf("Error parsing timestamp: %s", err)
	}

	return map[string]string{
		"hostname":            hostname,
		"logger_name":         p.LoggerName,
		"level":               p.Level,
		"message":             p.Message,
		"service_environment": p.ServiceEnvironment,
		"service_name":        p.ServiceName,
		"service_pipeline":    p.ServicePipeline,
		"service_version":     p.ServiceVersion,
		"stack_trace":         p.StackTrace,
		"thread_name":         p.ThreadName,
		"timestamp":           timestamp,
		"version":             strconv.Itoa(p.Version),
	}, nil
}

// Name ...
func (p MakoJSON) Name() string {
	return "makojson"
}
