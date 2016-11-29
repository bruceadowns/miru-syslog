package comm

import (
	"log"

	"github.com/bruceadowns/miru-syslog/miru"
)

// PostChan creates and returns a buffered channel used to post events to stumptown
func PostChan(size int, addr, url string) (ch chan *miru.LogEvent) {
	ch = make(chan *miru.LogEvent, size)

	go func() {
		for {
			select {
			case logEvent := <-ch:
				log.Printf("Send to stumptown: %s", logEvent)
				miru.Post(addr, url, logEvent)
			}
		}
	}()

	return
}
