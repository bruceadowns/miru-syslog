package lib

import (
	"bytes"
	"compress/gzip"
	"log"
	"time"
)

// SwitchBoard ...
type SwitchBoard struct {
	ParseChan     chan *Packet
	MiruAccumChan chan *LogEvent
	MiruPostChan  chan LogEvents

	S3AccumChan chan *Packet
	S3PostChan  chan bytes.Buffer
}

// ParseChan creates and returns a buffered channel used to capture line input
func ParseChan(size int, accumChan chan *LogEvent) (ch chan *Packet) {
	ch = make(chan *Packet, size)

	go func() {
		t := 0
		for {
			select {
			case p := <-ch:
				t++
				log.Printf("Parse message [%d]: %s", t, p)

				if e, err := p.Mill(); err != nil {
					log.Printf("Error parsing message: [%s] [%s]", p, err)
				} else {
					log.Printf("Send log event to miru accumulator: [%s]", e)
					accumChan <- e
				}
			}
		}
	}()

	return
}

// MiruAccumChan creates and returns a buffered channel used to accumulate events
func MiruAccumChan(size int, batchSize int, delay time.Duration, postChan chan LogEvents) (ch chan *LogEvent) {
	ch = make(chan *LogEvent, size)

	go func() {
		t := 0
		var logEvents LogEvents

		for {
			select {
			case <-time.After(delay):
				if len(logEvents) > 0 {
					log.Printf("Post miru log events: %d (delay: %dms)", len(logEvents), delay)
					postChan <- logEvents
					logEvents = make(LogEvents, 0)
				}

			case e := <-ch:
				t++
				log.Printf("Accumulate log event [%d] for miru: %s", t, e)
				logEvents = append(logEvents, *e)

				if len(logEvents) >= batchSize {
					log.Printf("Post log events to miru: %d", len(logEvents))
					postChan <- logEvents
					logEvents = make(LogEvents, 0)
				}
			}
		}
	}()

	return
}

// MiruPostChan creates and returns a buffered channel used to post events to stumptown
func MiruPostChan(size int, addr, url string, delaySuccess, delayError time.Duration) (ch chan LogEvents) {
	ch = make(chan LogEvents, size)

	go func() {
		t, tl := 0, 0
		for {
			select {
			case logEvents := <-ch:
				t++
				tl += len(logEvents)
				log.Printf("Send %d events to stumptown [%d - %d]", len(logEvents), t, tl)
				if len(logEvents) > 0 {
					if err := logEvents.Post(addr, url, delaySuccess, delayError); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}()

	return
}

// S3AccumChan ...
func S3AccumChan(size, batchBytes int, delay time.Duration, s3PostChan chan bytes.Buffer) (ch chan *Packet) {
	ch = make(chan *Packet, size)

	go func() {
		t, tl := 0, 0
		bb := bytes.Buffer{}
		gw := gzip.NewWriter(&bb)

		for {
			select {
			case <-time.After(delay):
				if bb.Len() > 0 {
					gw.Close()

					log.Printf("Post %d bytes to S3 (delay: %dms)", bb.Len(), delay)
					s3PostChan <- bb

					bb = bytes.Buffer{}
					gw = gzip.NewWriter(&bb)
				}

			case p := <-ch:
				t++
				tl += len(p.Message)
				log.Printf("Accumulate (gzip) packet [%d - %d] for S3: %s", t, tl, p)
				gw.Write(p.Message)

				if bb.Len() >= batchBytes {
					gw.Close()

					log.Printf("Post %d bytes to S3", bb.Len())
					s3PostChan <- bb

					bb = bytes.Buffer{}
					gw = gzip.NewWriter(&bb)
				}
			}
		}
	}()

	return
}

// S3PostChan ...
func S3PostChan(size int, a AWSInfo, delaySuccess, delayError time.Duration) (ch chan bytes.Buffer) {
	ch = make(chan bytes.Buffer, size)

	go func() {
		t, tl := 0, 0
		for {
			select {
			case bb := <-ch:
				t++
				tl += bb.Len()
				log.Printf("Send %d bytes to S3. [%d - %d]", bb.Len(), t, tl)
				PostS3(bb, a, delaySuccess, delayError)
			}
		}
	}()

	return
}
