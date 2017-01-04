package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	return fmt.Sprintf(
		"datacenter: %s - cluster: %s - service: %s - message: %s",
		l.DataCenter, l.Cluster, l.Service, Trunc(l.Message))
}

// LogEvents ...
type LogEvents []LogEvent

// helper post function to facilitate retry
func doPost(u string, bb *bytes.Buffer) error {
	resp, err := http.DefaultClient.Post(u, "application/json", bb)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Invalid status code [%d] posting to %s", resp.StatusCode, u)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Unexpected error reading accepted post", err)
	}

	log.Printf("Response: %d '%s'", resp.StatusCode, body)
	return nil
}

// Post sends log events to stumptown
func (l *LogEvents) Post(a, u string, delaySuccess, delayError time.Duration) error {
	if len(a) == 0 {
		log.Print("Stumptown address is empty.")
		return nil
	}
	if len(u) == 0 {
		log.Print("Stumptown intake url is empty.")
		return nil
	}

	log.Printf("Send %d log events to stumptown.", len(*l))
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(l); err != nil {
		return err
	}
	log.Printf("Post buffer: %s", Trunc(buf.String()))

	url := fmt.Sprintf("http://%s%s", a, u)
	log.Print(url)

	for {
		err := doPost(url, buf)
		if err == nil {
			if delaySuccess > 0 {
				log.Printf("Miru delay on success %dms", delaySuccess)
				time.Sleep(delaySuccess)
			}

			break
		}

		log.Printf("Error posting to %s: %s", url, err)

		if delayError > 0 {
			log.Printf("Miru delay on error %dms", delayError)
			time.Sleep(delayError)
		}
	}

	return nil
}
