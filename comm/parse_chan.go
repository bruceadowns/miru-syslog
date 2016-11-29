package comm

import (
	"fmt"
	"log"
	"net"

	"github.com/bruceadowns/miru-syslog/miru"

	"github.com/jeromer/syslogparser/rfc3164"
	"github.com/jeromer/syslogparser/rfc5424"
)

// MakoJSON holds mako structured json
type MakoJSON struct {
	Timestamp          string `json:"@timestamp,omitempty"`
	Version            int    `json:"@version,omitempty"`
	Message            string `json:"message,omitempty"`
	LoggerName         string `json:"logger_name,omitempty"`
	ThreadName         string `json:"thread_name,omitempty"`
	Level              string `json:"level,omitempty"`
	LevelValue         int    `json:"level_value,omitempty"`
	ServiceName        string `json:"service_name,omitempty"`
	ServiceEnvironment string `json:"service_environment,omitempty"`
	ServicePipeline    string `json:"service_pipeline,omitempty"`
	ServiceVersion     string `json:"service_version,omitempty"`
}

// Packet holds the incoming traffic info
type Packet struct {
	Address  net.Addr
	Message  []byte
	LogEvent *miru.LogEvent
}

func (p *Packet) String() string {
	return fmt.Sprintf("Address: %s '%s'", p.Address, p.Message)
}

// IsValid returns T/F
func (p *Packet) IsValid() bool {
	if len(p.Address.String()) == 0 {
		log.Print("Address is empty")
		return false
	}

	if len(p.Message) == 0 {
		log.Print("Message is empty")
		return false
	}

	return true
}

// Mill determines message type and parses into a LogEvent
func (p *Packet) Mill() (res *miru.LogEvent) {
	log.Printf("%s", p)

	/*
		{"@timestamp":"2016-11-29T00:01:44.658+00:00","@version":1,"message":"172.16.3.0 - developer [29/Nov/2016:00:01:44 +0000] \"GET /releases HTTP/1.1\" 200 - \"https://cloud-jcx-api.ms-integ.svc.jivehosted.com/ui\" \"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36\" 221","logger_name":"http.request","thread_name":"dw-14717","level":"INFO","level_value":20000,"service_name":"cloud-jcx-api","service_environment":"ms-integ","service_pipeline":"main","service_version":"caaa4f863f114bd697dbd06569fb01f9c8667681"}

		Nov 29 09:10:30 LT-A10-122418 login[32465]: USER_PROCESS: 32465 ttys010
		Nov 29 17:10:18 ip-10-126-5-155 dhclient[2346]: bound to 10.126.5.155 -- renewal in 1649 seconds.

		2016-11-29T09:11:04.880190-08:00 soa-prime-data1 /usr/sbin/gmetad[2263]: data_thread() got no answer from any [my cluster] datasource
	*/

	for {
		{
			parser := rfc3164.NewParser(p.Message)
			if err := parser.Parse(); err == nil {
				res = &miru.LogEvent{}
				break
			}
		}

		{
			parser := rfc5424.NewParser(p.Message)
			if err := parser.Parse(); err == nil {
				res = &miru.LogEvent{}
				break
			}
		}

		{
			res = &miru.LogEvent{
				DataCenter: "bad-dc",
				Cluster:    "bad-cluster",
				Host:       "bad-host",
				Service:    "bad-service",
				Instance:   "bad-instance",
				Version:    "1.0",
				Level:      miru.LevelInfo,
				Message:    fmt.Sprintf("%s", p.Message),
			}
			break
		}
	}

	if res == nil {
		log.Printf("Message from %s not parsed: %s", p.Address, p.Message)
	} else {
		log.Print(res)
	}

	return
}

// ParseChan creates and returns a buffered channel used to capture line input
func ParseChan(size int, postChan chan *miru.LogEvent) (ch chan *Packet) {
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
