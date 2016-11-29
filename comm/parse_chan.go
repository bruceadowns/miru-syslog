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

// IsValid returns T/F
func (p *Packet) IsValid() bool {
	if p.Address.String() == "" {
		log.Print("Address is empty")
		return false
	}

	if len(p.Message) == 0 {
		log.Print("Message is empty")
		return false
	}

	return true
}

// Parse determines message type and parses into a LogEvent
func (p *Packet) Parse() error {
	log.Printf("Parse '%s'", p)
	return nil
}

// ParseChan creates and returns a buffered channel used to capture line input
func ParseChan(size int) (ch chan *Packet) {
	ch = make(chan *Packet, size)

	go func() {
		for {
			select {
			case m := <-ch:
				if !m.IsValid() {
					continue
				}

				if err := m.Parse(); err != nil {
					log.Printf("Error parsing message '%s' [%s].", m, err)
				}
			}
		}
	}()

	return
}
