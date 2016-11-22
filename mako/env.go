package mako

import (
	"log"
	"os"
)

// Env tracks mako related environment variables
type Env struct {
	serviceID   string
	environment string
	pipeline    string
	version     string
	statsdHost  string
	statsdPort  int
}

// ActiveEnv holds the active environment variables
var ActiveEnv Env

func hideEnvInit() {
	ActiveEnv.environment = os.Getenv("MAKO_ENVIRONMENT")
	if ActiveEnv.environment == "" {
		log.Fatalf("MAKO_ENVIRONMENT not present in mako environment")
	}

	ActiveEnv.pipeline = os.Getenv("MAKO_PIPELINE")
	if ActiveEnv.pipeline == "" {
		log.Fatalf("MAKO_PIPELINE not present in mako environment")
	}

	ActiveEnv.version = os.Getenv("MAKO_VERSION")
	if ActiveEnv.version == "" {
		log.Fatalf("MAKO_VERSION not present in mako environment")
	}
}
