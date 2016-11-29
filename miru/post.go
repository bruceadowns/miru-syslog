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
func PostOneEvent(addr string) error {
	if len(addr) == 0 {
		log.Print("Stumptown address is empty.")
		return nil
	}

	events := []LogEvent{{
		DataCenter: "bad-dc",
		Cluster:    "bad-cluster",
		Host:       "bad-host",
		Service:    "bad-service",
		Instance:   "bad-instance",
		Version:    "1.0",
		Level:      LevelInfo,
		Message:    "My info message",
	}}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(events); err != nil {
		return err
	}
	log.Print(buf)

	resp, err := http.Post(fmt.Sprintf("http://%s/miru/stumptown/intake", addr), "application/json", buf)
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
func PostManyEvents(addr string) error {
	if len(addr) == 0 {
		log.Print("Stumptown address is empty.")
		return nil
	}

	const SIZE = 10
	events := make([]LogEvent, SIZE)

	for i := 0; i < SIZE; i++ {
		event := LogEvent{
			DataCenter: "bad-dc",
			Cluster:    "bad-cluster",
			Host:       "bad-host",
			Service:    "bad-service",
			Instance:   "bad-instance",
			Version:    "1.0",
			Level:      LevelWarn,
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

	resp, err := http.Post(fmt.Sprintf("http://%s/miru/stumptown/intake", addr), "application/json", buf)
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
