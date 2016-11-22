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
		ServiceName:        ActiveEnv.serviceID,
		ServiceEnvironment: ActiveEnv.environment,
		ServicePipeline:    ActiveEnv.pipeline,
		ServiceVersion:     ActiveEnv.version,
		Message:            message,
		Level:              level,
	})
	if err != nil {
		log.Fatalf("Error occurred logging json %s: %s", level, message)
	}

	fmt.Fprintln(os.Stdout, string(j))
}

// LogError logs a debug message to stdout as json
func LogError(message string) {
	logPrint(message, "ERROR")
}

// LogInfo logs a debug message to stdout as json
func LogInfo(message string) {
	logPrint(message, "INFO")
}
