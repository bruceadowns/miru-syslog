package lib

import (
	"log"
)

// PostChan creates and returns a buffered channel used to post events to stumptown
func PostChan(size int, addr, url string) (ch chan *MiruLogEvent) {
	ch = make(chan *MiruLogEvent, size)

	go func() {
		for {
			select {
			case logEvent := <-ch:
				log.Printf("Send to stumptown: %s", logEvent)
				Post(addr, url, logEvent)
			}
		}
	}()

	return
}
