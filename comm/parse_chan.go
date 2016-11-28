package comm

import (
	"fmt"
	"log"
	"net"
)

// Packet holds the incoming traffic info
type Packet struct {
	Address net.Addr
	Message []byte
}

func (p *Packet) String() string {
	return fmt.Sprintf("%s: '%s'", p.Address, p.Message)
}

func parse(b []byte) error {
	log.Printf("Parse '%s'", b)
	return nil
}

// ParseChan creates and returns a buffered channel used to capture line input
func ParseChan(size int) (ch chan *Packet) {
	ch = make(chan *Packet, size)

	go func() {
		for {
			select {
			case m := <-ch:
				if len(m.Address.String()) == 0 {
					log.Print("Address is empty")
					continue
				}

				if len(m.Message) == 0 {
					log.Print("Message is empty")
					continue
				}

				if err := parse(m.Message); err != nil {
					log.Printf("Error parsing message '%s' [%s].", m, err)
				}
			}
		}
	}()

	return
}
