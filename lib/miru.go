package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// LevelInfo stringizes info log level
	LevelInfo = "INFO"

	// LevelWarn stringizes warn log level
	LevelWarn = "WARN"
)

// LogEvent holds the stumptown event
type LogEvent struct {
	DataCenter       string   `json:"datacenter,omitempty"`
	Cluster          string   `json:"cluster,omitempty"`
	Host             string   `json:"host,omitempty"`
	Service          string   `json:"service,omitempty"`
	Instance         string   `json:"instance,omitempty"`
	Version          string   `json:"version,omitempty"`
	Level            string   `json:"level,omitempty"`
	ThreadName       string   `json:"threadName,omitempty"`
	LoggerName       string   `json:"loggerName,omitempty"`
	Message          string   `json:"message,omitempty"`
	Timestamp        string   `json:"timestamp,omitempty"`
	ThrownStackTrace []string `json:"thrownStackTrace,omitempty"`
}

func (l *LogEvent) String() string {
	return fmt.Sprintf("datacenter: %s - cluster: %s - message: %s", l.DataCenter, l.Cluster, l.Message)
}

// Post sends a single log event to stumptown
func (l *LogEvent) Post(a, u string) error {
	if len(a) == 0 {
		log.Print("Stumptown address is empty.")
		return nil
	}
	if len(u) == 0 {
		log.Print("Stumptown intake url is empty.")
		return nil
	}

	events := []*LogEvent{l}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(events); err != nil {
		return err
	}
	log.Print(buf)

	url := fmt.Sprintf("http://%s%s", a, u)
	log.Print(url)
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Response: %d '%s'\n", resp.StatusCode, body)
	return nil
}
