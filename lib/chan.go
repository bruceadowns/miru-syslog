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
		for {
			select {
			case p := <-ch:
				log.Printf("Parse message: [%s]", p)

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
		var logEvents LogEvents

		for {
			select {
			case <-time.After(delay):
				if len(logEvents) > 0 {
					log.Printf("Post miru log events: %d (delay: %d)", len(logEvents), delay)
					postChan <- logEvents
					logEvents = make(LogEvents, 0)
				}

			case e := <-ch:
				log.Printf("Accumulate log event for miru: %s", e)
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
func MiruPostChan(size int, addr, url string) (ch chan LogEvents) {
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

// S3AccumChan ...
func S3AccumChan(size, batchBytes int, delay time.Duration, s3PostChan chan bytes.Buffer) (ch chan *Packet) {
	ch = make(chan *Packet, size)

	go func() {
		bb := bytes.Buffer{}
		gw := gzip.NewWriter(&bb)

		for {
			select {
			case <-time.After(delay):
				if bb.Len() > 0 {
					gw.Close()

					log.Printf("Post %d bytes to S3 (delay: %d)", bb.Len(), delay)
					s3PostChan <- bb

					bb = bytes.Buffer{}
					gw = gzip.NewWriter(&bb)
				}

			case p := <-ch:
				log.Printf("Accumulate (gzip) packet for S3: %s", p)
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
func S3PostChan(size int, a AWSInfo) (ch chan bytes.Buffer) {
	ch = make(chan bytes.Buffer, size)

	go func() {
		for {
			select {
			case bb := <-ch:
				log.Printf("Send %d bytes to S3", bb.Len())
				PostS3(bb, a)
			}
		}
	}()

	return
}
