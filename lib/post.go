package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Post sends a single log event to stumptown
func Post(a, u string, logEvent *LogEvent) error {
	if len(a) == 0 {
		log.Print("Stumptown address is empty.")
		return nil
	}
	if len(u) == 0 {
		log.Print("Stumptown intake url is empty.")
		return nil
	}

	events := []*LogEvent{logEvent}

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
