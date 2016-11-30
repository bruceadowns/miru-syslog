package lib

import (
	"log"
)

// ParseChan creates and returns a buffered channel used to capture line input
func ParseChan(size int, postChan chan *LogEvent) (ch chan *Packet) {
	ch = make(chan *Packet, size)

	go func() {
		for {
			select {
			case m := <-ch:
				if !m.IsValid() {
					continue
				}

				if logEvent := m.Mill(); logEvent == nil {
					log.Printf("Error parsing message: [%s]", m)
				} else {
					log.Printf("Posting log event: [%s]", logEvent)
					postChan <- logEvent
				}
			}
		}
	}()

	return
}

// PostChan creates and returns a buffered channel used to post events to stumptown
func PostChan(size int, addr, url string) (ch chan *LogEvent) {
	ch = make(chan *LogEvent, size)

	go func() {
		for {
			select {
			case logEvent := <-ch:
				log.Printf("Send to stumptown: %s", logEvent)
				logEvent.Post(addr, url)
			}
		}
	}()

	return
}
