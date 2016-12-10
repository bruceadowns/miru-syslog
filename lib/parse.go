package lib

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/bruceadowns/syslogparser"
	"github.com/bruceadowns/syslogparser/journaljson"
	"github.com/bruceadowns/syslogparser/journalmako"
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

var remoteTypeCache = make(map[net.Addr]string)

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

type noopParser struct {
	buff []byte
	host net.Addr
}

func newNoopParser(buff []byte, host net.Addr) *noopParser {
	return &noopParser{buff: buff, host: host}
}

// Parse ...
func (p *noopParser) Parse() error {
	return nil
}

// Dump ...
func (p *noopParser) Dump() syslogparser.LogParts {
	return syslogparser.LogParts{
		"message":  string(p.buff),
		"hostname": p.host.String(),
	}
}

func populate(p syslogparser.LogParser) (res *LogEvent) {
	if p == nil {
		return
	}

	logParts := p.Dump()

	app := logParts["service_name"]
	if len(app) == 0 {
		app = logParts["app_name"]
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
	}

	level := logParts["level"]
	if len(level) == 0 {
		switch logParts["severity"] {
		case "0", "1", "2", "3":
			level = "ERROR"
		case "4":
			level = "WARN"
		case "7":
			level = "DEBUG"
		//case "5", "6":
		default:
			level = "INFO"
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

	stackTrace := []string{}
	fullStackTrace := logParts["stack_trace"]
	if len(fullStackTrace) > 0 {
		stackTrace = strings.Split(fullStackTrace, "\n")
	}

	res = &LogEvent{
		DataCenter:       logParts["service_environment"],
		Cluster:          logParts["service_pipeline"],
		Host:             logParts["hostname"],
		Service:          app,
		Instance:         pid,
		Version:          version,
		Level:            level,
		ThreadName:       logParts["thread_name"],
		LoggerName:       logParts["logger_name"],
		Message:          message,
		Timestamp:        timestamp,
		ThrownStackTrace: stackTrace,
	}

	return
}

func determineParser(p *Packet) (res syslogparser.LogParser) {
	res = journalmako.NewParser(p.Message)
	if err := res.Parse(); err == nil {
		remoteTypeCache[p.Address] = "journalmako"
		log.Printf("%s - %s", p.Address.String(), "journalmako")
		return
	}

	res = journaljson.NewParser(p.Message)
	if err := res.Parse(); err == nil {
		remoteTypeCache[p.Address] = "journaljson"
		log.Printf("%s - %s", p.Address.String(), "journaljson")
		return
	}

	res = mako.NewParser(p.Message, p.Address)
	if err := res.Parse(); err == nil {
		remoteTypeCache[p.Address] = "mako"
		log.Printf("%s - %s", p.Address.String(), "mako")
		return
	}

	res = syslogmako.NewParser(p.Message)
	if err := res.Parse(); err == nil {
		remoteTypeCache[p.Address] = "syslogmako"
		log.Printf("%s - %s", p.Address.String(), "syslogmako")
		return
	}

	res = rfc5424raw.NewParser(p.Message)
	if err := res.Parse(); err == nil {
		remoteTypeCache[p.Address] = "rfc5424raw"
		log.Printf("%s - %s", p.Address.String(), "rfc5424raw")
		return
	}

	res = rfc3164raw.NewParser(p.Message)
	if err := res.Parse(); err == nil {
		remoteTypeCache[p.Address] = "rfc3164raw"
		log.Printf("%s - %s", p.Address.String(), "rfc3164raw")
		return
	}

	res = rfc3164.NewParser(p.Message)
	if err := res.Parse(); err == nil {
		remoteTypeCache[p.Address] = "rfc3164"
		log.Printf("%s - %s", p.Address.String(), "rfc3164")
		return
	}

	res = rfc5424.NewParser(p.Message)
	if err := res.Parse(); err == nil {
		remoteTypeCache[p.Address] = "rfc5424"
		log.Printf("%s - %s", p.Address.String(), "rfc5424")
		return
	}

	res = newNoopParser(p.Message, p.Address)
	log.Printf("%s - %s", p.Address.String(), "noop")
	return
}

// Mill determines message type and parses into a LogEvent
func (p *Packet) Mill() (res *LogEvent) {
	log.Printf("%s", p)

	var parser syslogparser.LogParser

	switch remoteTypeCache[p.Address] {

	case "mako":
		parser = mako.NewParser(p.Message, p.Address)
		if err := parser.Parse(); err != nil {
			log.Print(err)
			parser = determineParser(p)
		}

	case "syslogmako":
		parser = syslogmako.NewParser(p.Message)
		if err := parser.Parse(); err != nil {
			log.Print(err)
			parser = determineParser(p)
		}

	case "rfc5424raw":
		parser = rfc5424raw.NewParser(p.Message)
		if err := parser.Parse(); err != nil {
			log.Print(err)
			parser = determineParser(p)
		}

	case "rfc3164raw":
		parser = rfc3164raw.NewParser(p.Message)
		if err := parser.Parse(); err != nil {
			log.Print(err)
			parser = determineParser(p)
		}

	case "rfc5424":
		parser = rfc5424.NewParser(p.Message)
		if err := parser.Parse(); err != nil {
			log.Print(err)
			parser = determineParser(p)
		}

	case "rfc3164":
		parser = rfc3164raw.NewParser(p.Message)
		if err := parser.Parse(); err != nil {
			log.Print(err)
			parser = determineParser(p)
		}

	case "noop":
		parser = newNoopParser(p.Message, p.Address)
		if err := parser.Parse(); err != nil {
			log.Print(err)
			parser = determineParser(p)
		}

	default:
		parser = determineParser(p)
	}

	res = populate(parser)
	if res == nil {
		log.Printf("Message from %s not parsed: %s", p.Address, p.Message)
	} else {
		log.Print(res)
	}

	return
}
