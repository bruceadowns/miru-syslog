package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Parser ...
type Parser interface {
	Name() string
	Extract(hostname string, bb *bytes.Buffer) (map[string]string, error)
}

// Packet holds the incoming traffic info
type Packet struct {
	Address  string
	Message  []byte
	LogEvent *LogEvent
}

// PreTag pretag json
type PreTag struct {
	Name, Type string
}

var hostToTypeCache = make(map[string]string)

func init() {
	if c, err := ioutil.ReadFile("pretag.json"); err == nil {
		var t []PreTag
		err = json.Unmarshal(c, &t)
		if err != nil {
			log.Printf("Error unmarshalling pretag.json: %s", err)
		}

		for _, v := range t {
			hostToTypeCache[v.Name] = v.Type
			log.Printf("pretag %s:%s", v.Name, v.Type)
		}
	} else {
		log.Printf("Error reading pretag.json: %s", err)
	}
}

func (p *Packet) String() string {
	return fmt.Sprintf("Address: %s '%s'", p.Address,
		Trunc(string(p.Message)))
}

func (p *Packet) determineParser() (fields map[string]string, err error) {
	var parser Parser

	for {
		parser = JournalJSONMako{}
		if fields, err = parser.Extract(p.Address, bytes.NewBuffer(p.Message)); err == nil {
			hostToTypeCache[p.Address] = parser.Name()
			break
		}

		parser = MakoJSON{}
		if fields, err = parser.Extract(p.Address, bytes.NewBuffer(p.Message)); err == nil {
			hostToTypeCache[p.Address] = parser.Name()
			break
		}

		parser = Base{}
		fields, err = parser.Extract(p.Address, bytes.NewBuffer(p.Message))
		break
	}

	if err == nil {
		log.Printf("determined %s - %s", p.Address, parser.Name())
	}

	return
}

// Mill determines message type and parses into a LogEvent
func (p *Packet) Mill() (res *LogEvent, err error) {
	if p.Address == "" {
		return nil, fmt.Errorf("Address is empty")
	}
	if len(p.Message) == 0 {
		return nil, fmt.Errorf("Address is empty")
	}

	log.Printf("Mill packet: %s", p)

	var fields map[string]string

	switch hostToTypeCache[p.Address] {
	case "journaljsonmako":
		parser := JournalJSONMako{}
		if fields, err = parser.Extract(p.Address, bytes.NewBuffer(p.Message)); err != nil {
			log.Printf("Error parsing %s [%s]", parser.Name(), err)

			parser := Base{}
			if fields, err = parser.Extract(p.Address, bytes.NewBuffer(p.Message)); err != nil {
				return nil, err
			}
		}

	case "makojson":
		parser := MakoJSON{}
		if fields, err = parser.Extract(p.Address, bytes.NewBuffer(p.Message)); err != nil {
			log.Printf("Error parsing %s [%s]", parser.Name(), err)

			parser := Base{}
			if fields, err = parser.Extract(p.Address, bytes.NewBuffer(p.Message)); err != nil {
				return nil, err
			}
		}

	case "base":
		parser := Base{}
		if fields, err = parser.Extract(p.Address, bytes.NewBuffer(p.Message)); err != nil {
			return nil, err
		}

	default:
		if fields, err = p.determineParser(); err != nil {
			return nil, err
		}
	}

	stackTrace := []string{}
	fullStackTrace := fields["stack_trace"]
	if len(fullStackTrace) > 0 {
		stackTrace = strings.Split(fullStackTrace, "\n")
	}

	return &LogEvent{
		DataCenter:       fields["service_environment"],
		Cluster:          fields["service_pipeline"],
		Host:             fields["hostname"],
		Service:          fields["service_name"],
		Instance:         fields["service_version"],
		Version:          fields["version"],
		Level:            fields["level"],
		ThreadName:       fields["thread_name"],
		LoggerName:       fields["logger_name"],
		Message:          fields["message"],
		Timestamp:        fields["timestamp"],
		ThrownStackTrace: stackTrace,
	}, nil
}
