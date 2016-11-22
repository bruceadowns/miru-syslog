package mako

import (
	"log"
	"os"
)

// MakoEnv tracks mako related environment variables
type MakoEnv struct {
	serviceID   string
	environment string
	pipeline    string
	version     string
	statsdHost  string
	statsdPort  int
}

// ActiveMakoEnv holds the active environment variables
var ActiveMakoEnv MakoEnv

func hideEnvInit() {
	ActiveMakoEnv.environment = os.Getenv("MAKO_ENVIRONMENT")
	if ActiveMakoEnv.environment == "" {
		log.Fatalf("MAKO_ENVIRONMENT not present in mako environment")
	}

	ActiveMakoEnv.pipeline = os.Getenv("MAKO_PIPELINE")
	if ActiveMakoEnv.pipeline == "" {
		log.Fatalf("MAKO_PIPELINE not present in mako environment")
	}

	ActiveMakoEnv.version = os.Getenv("MAKO_VERSION")
	if ActiveMakoEnv.version == "" {
		log.Fatalf("MAKO_VERSION not present in mako environment")
	}
}
