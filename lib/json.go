package lib

import "fmt"

const (
	// LevelInfo stringizes info log level
	LevelInfo = "INFO"

	// LevelWarn stringizes warn log level
	LevelWarn = "WARN"
)

// MiruLogEvent holds the stumptown event
type MiruLogEvent struct {
	DataCenter       string   `json:"datacenter,omitempty"`
	Cluster          string   `json:"cluster,omitempty"`
	Host             string   `json:"host,omitempty"`
	Service          string   `json:"service,omitempty"`
	Instance         string   `json:"instance,omitempty"`
	Version          string   `json:"version,omitempty"`
	Level            string   `json:"level,omitempty"`
	ThreadName       string   `json:"threadName,omitempty"`
	LoggerName       string   `json:"loggerName,omitempty"`
	Message          string   `json:"message,omitempty"`
	Timestamp        string   `json:"timestamp,omitempty"`
	ThrownStackTrace []string `json:"thrownStackTrace,omitempty"`
}

func (l *MiruLogEvent) String() string {
	return fmt.Sprintf("datacenter: %s - cluster: %s - message: %s", l.DataCenter, l.Cluster, l.Message)
}
