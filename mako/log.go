package mako

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type jsonLogEntry struct {
	Timestamp          string `json:"timestamp"`
	ServiceName        string `json:"service_name"`
	ServiceEnvironment string `json:"service_environment"`
	ServicePipeline    string `json:"service_pipeline"`
	ServiceVersion     string `json:"service_version"`
	Message            string `json:"message"`
	Level              string `json:"level"`
}

func logPrint(message, level string) {
	j, err := json.Marshal(&jsonLogEntry{
		Timestamp:          time.Now().Format(time.RFC3339),
		ServiceName:        ActiveMakoEnv.serviceID,
		ServiceEnvironment: ActiveMakoEnv.environment,
		ServicePipeline:    ActiveMakoEnv.pipeline,
		ServiceVersion:     ActiveMakoEnv.version,
		Message:            message,
		Level:              level,
	})
	if err != nil {
		log.Fatalf("Error occurred logging json %s: %s", level, message)
	}

	fmt.Fprintln(os.Stdout, string(j))
}

// LogDebug logs a debug message to stdout as json
func LogDebug(message string) {
	logPrint(message, "DEBUG")
}

// LogInfo logs a debug message to stdout as json
func LogInfo(message string) {
	logPrint(message, "INFO")
}
