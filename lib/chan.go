package lib

import (
	"log"
	"time"
)

// ParseChan creates and returns a buffered channel used to capture line input
func ParseChan(size int, accumChan chan LogEvent) (ch chan Packet) {
	ch = make(chan Packet, size)

	go func() {
		for {
			select {
			case m := <-ch:
				if !m.IsValid() {
					continue
				}

				if logEvent, err := m.Mill(); err != nil {
					log.Printf("Error parsing message: [%s]", m)
				} else {
					log.Printf("Send log event to accumulator: [%s]", logEvent)
					accumChan <- *logEvent
				}
			}
		}
	}()

	return
}

// AccumChan creates and returns a buffered channel used to accumulate events
func AccumChan(size int, batchSize int, delay time.Duration, postChan chan LogEvents) (ch chan LogEvent) {
	ch = make(chan LogEvent, size)

	go func() {
		var logEvents LogEvents
		for {
			select {
			case <-time.After(delay):
				if len(logEvents) > 0 {
					log.Printf("Post log events: %d (delay: %d)", len(logEvents), delay)
					postChan <- logEvents
					logEvents = make(LogEvents, 0)
				}
			case logEvent := <-ch:
				log.Printf("Accumulate log event: %s", logEvent)
				logEvents = append(logEvents, logEvent)

				if len(logEvents) >= batchSize {
					log.Printf("Post log events: %d", len(logEvents))
					postChan <- logEvents
					logEvents = make(LogEvents, 0)
				}
			}
		}
	}()

	return
}

// PostChan creates and returns a buffered channel used to post events to stumptown
func PostChan(size int, addr, url string) (ch chan LogEvents) {
	ch = make(chan LogEvents, size)

	go func() {
		for {
			select {
			case logEvents := <-ch:
				log.Printf("Send %d events to stumptown", len(logEvents))
				if len(logEvents) > 0 {
					if err := logEvents.Post(addr, url); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}()

	return
}
