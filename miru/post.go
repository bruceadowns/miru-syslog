package miru

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// PostOneEvent demonstrates sending a single log event to stumptown
func PostOneEvent() error {
	events := []miruLogEvent{miruLogEvent{
		DataCenter: "bad-dc",
		Cluster:    "bad-cluster",
		Host:       "bad-host",
		Service:    "bad-service",
		Instance:   "bad-instance",
		Version:    "1.0",
		Level:      levelInfo,
		Message:    "My info message",
	}}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(events); err != nil {
		return err
	}
	log.Print(buf)

	resp, err := http.Post("http://10.126.5.155:10004/miru/stumptown/intake", "application/json", buf)
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

// PostManyEvents demonstrates sending more than one log event to stumptown
func PostManyEvents() error {
	const SIZE = 10
	events := make([]miruLogEvent, SIZE)

	for i := 0; i < SIZE; i++ {
		event := miruLogEvent{
			DataCenter: "bad-dc",
			Cluster:    "bad-cluster",
			Host:       "bad-host",
			Service:    "bad-service",
			Instance:   "bad-instance",
			Version:    "1.0",
			Level:      levelWarn,
			Message:    fmt.Sprintf("My %d message", i),
		}

		//events = append(events, event)
		events[i] = event
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(events); err != nil {
		return err
	}
	log.Print(buf)

	resp, err := http.Post("http://10.126.5.155:10004/miru/stumptown/intake", "application/json", buf)
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
