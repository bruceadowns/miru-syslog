package lib

import (
	"fmt"
	"log"
	"net"

	"github.com/bruceadowns/syslogparser"
	"github.com/bruceadowns/syslogparser/mako"
	"github.com/bruceadowns/syslogparser/rfc3164"
	"github.com/bruceadowns/syslogparser/rfc3164raw"
	"github.com/bruceadowns/syslogparser/rfc5424"
	"github.com/bruceadowns/syslogparser/rfc5424raw"
	"github.com/bruceadowns/syslogparser/syslogmako"
)

// Packet holds the incoming traffic info
type Packet struct {
	Address  net.Addr
	Message  []byte
	LogEvent *LogEvent
}

var remoteTypeCache = make(map[string]string)

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

func populate(p syslogparser.LogParser) (res *LogEvent) {
	if p == nil {
		return
	}

	logParts := p.Dump()

	app := logParts["app_name"]
	if len(app) == 0 {
		app = logParts["service_name"]
	}

	pid := logParts["service_version"]
	if len(pid) == 0 {
		pid = logParts["proc_id"]
		if len(pid) == 0 {
			pid = logParts["tag"]
		}
	}

	version := logParts["@version"]
	if len(version) == 0 {
		version = logParts["version"]
		if len(version) == 0 {
			version = logParts["v"]
		}
	}

	message := logParts["message"]
	if len(message) == 0 {
		message = logParts["content"]
	}

	timestamp := logParts["@timestamp"]
	if len(timestamp) == 0 {
		timestamp = logParts["timestamp"]
	}

	res = &LogEvent{
		DataCenter: logParts["service_environment"],
		Cluster:    logParts["service_pipeline"],
		Host:       logParts["hostname"],
		Service:    app,
		Instance:   pid,
		Version:    version,
		Level:      logParts["level"],
		ThreadName: logParts["thread_name"],
		LoggerName: logParts["logger_name"],
		Message:    message,
		Timestamp:  timestamp,
		//ThrownStackTrace
	}

	return
}

// Mill determines message type and parses into a LogEvent
func (p *Packet) Mill() (res *LogEvent) {
	log.Printf("%s", p)

	var parser syslogparser.LogParser

	for {
		// check cache for known type

		{
			parser = mako.NewParser(p.Message, p.Address.String())
			if err := parser.Parse(); err == nil {
				remoteTypeCache[p.Address.String()] = "mako"
				log.Printf("mako")
				break
			}
		}

		{
			parser = syslogmako.NewParser(p.Message)
			if err := parser.Parse(); err == nil {
				remoteTypeCache[p.Address.String()] = "syslogmako"
				log.Printf("syslogmako")
				break
			}
		}

		{
			parser = rfc5424raw.NewParser(p.Message)
			if err := parser.Parse(); err == nil {
				remoteTypeCache[p.Address.String()] = "rfc5424raw"
				log.Printf("rfc5424raw")
				break
			}
		}

		{
			parser = rfc3164raw.NewParser(p.Message)
			if err := parser.Parse(); err == nil {
				remoteTypeCache[p.Address.String()] = "rfc3164raw"
				log.Printf("rfc3164raw")
				break
			}
		}

		{
			parser = rfc3164.NewParser(p.Message)
			if err := parser.Parse(); err == nil {
				remoteTypeCache[p.Address.String()] = "rfc3164"
				log.Printf("rfc3164")
				break
			}
		}

		{
			parser = rfc5424.NewParser(p.Message)
			if err := parser.Parse(); err == nil {
				remoteTypeCache[p.Address.String()] = "rfc5424"
				log.Printf("rfc5424")
				break
			}
		}

		log.Printf("none")
		parser = nil
		break
	}

	res = populate(parser)
	if res == nil {
		log.Printf("Message from %s not parsed: %s", p.Address, p.Message)
	} else {
		log.Print(res)
	}

	return
}
