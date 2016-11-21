package common

import (
	"log"
	"os"
)

var activeMakoEnv makoEnv

type makoEnv struct {
	serviceID   string
	environment string
	pipeline    string
	version     string
	statsdHost  string
	statsdPort  int
}

func hideEnvInit() {
	activeMakoEnv.environment = os.Getenv("MAKO_ENVIRONMENT")
	if activeMakoEnv.environment == "" {
		log.Fatalf("MAKO_ENVIRONMENT not present in mako environment")
	}

	activeMakoEnv.pipeline = os.Getenv("MAKO_PIPELINE")
	if activeMakoEnv.pipeline == "" {
		log.Fatalf("MAKO_PIPELINE not present in mako environment")
	}

	activeMakoEnv.version = os.Getenv("MAKO_VERSION")
	if activeMakoEnv.version == "" {
		log.Fatalf("MAKO_VERSION not present in mako environment")
	}
}
