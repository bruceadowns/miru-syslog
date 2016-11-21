package common

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/DataDog/datadog-go/statsd"
)

// DataDogClient is the client for logging to datadoghq
var DataDogClient *statsd.Client

func hideDdInit() {
	activeMakoEnv.statsdHost = os.Getenv("MAKO_STATSD_HOST")
	if activeMakoEnv.statsdHost == "" {
		log.Fatalf("MAKO_STATSD_HOST not present in mako environment")
	}

	sPort := os.Getenv("MAKO_STATSD_PORT")
	if sPort == "" {
		log.Fatalf("MAKO_STATSD_PORT not present in mako environment")
	}
	iPort, err := strconv.Atoi(sPort)
	if err != nil {
		log.Fatalf("Error occurred converting MAKO_STATSD_PORT: %s", sPort)
	}
	activeMakoEnv.statsdPort = iPort

	//telemetry to statsd_host:statsd_port
	DataDogClient, err = statsd.New(
		fmt.Sprintf("%s:%d",
			activeMakoEnv.statsdHost, activeMakoEnv.statsdPort))
	if err != nil {
		log.Fatal(err)
	}

	activeMakoEnv.serviceID = os.Getenv("MAKO_SERVICE_ID")
	if activeMakoEnv.serviceID == "" {
		log.Fatalf("MAKO_SERVICE_ID not present in mako environment")
	}

	DataDogClient.Namespace = activeMakoEnv.serviceID
}
