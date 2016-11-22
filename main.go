package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/bruceadowns/miru-syslog/mako"
	"github.com/bruceadowns/miru-syslog/miru"
)

type miruEnv struct {
	miruAddress string
}

var activeMiruEnv miruEnv

var timeStart time.Time

func adminRootHandler(w http.ResponseWriter, r *http.Request) {
	mako.LogInfo(fmt.Sprintf("%v", r))
	w.Write([]byte("Hello Admin Root"))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	mako.LogInfo(fmt.Sprintf("%v", r))
	mako.DataDogClient.Count("ping", 1, nil, 0)

	w.Write([]byte("pong"))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	mako.LogInfo(fmt.Sprintf("%v", r))
	w.Write([]byte("Hello Health Check"))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	mako.LogInfo(fmt.Sprintf("%v", r))
	w.Write([]byte("Hello Metrics"))
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	mako.LogInfo(fmt.Sprintf("%v", r))
	mako.DataDogClient.Count("serviceHandler", 1, nil, 0)
	mako.DataDogClient.SimpleEvent("serviceHandler", "Service Handler was called")
	w.Write([]byte("Hello Root"))
}

func udpMessagePump() error {
	log.Printf("Listen for udp traffic on %s", activeMiruEnv.miruAddress)

	pc, err := net.ListenPacket("udp", activeMiruEnv.miruAddress)
	if err != nil {
		return err
	}
	defer pc.Close()

	log.Printf("Handle UDP connections")
	for {
		buffer := make([]byte, 1024)
		pc.ReadFrom(buffer)
		log.Printf("Read udp buffer: %s", buffer)

		err := miru.PostManyEvents()
		if err != nil {
			log.Print(err)
		}
	}
}

func handleTCPConnection(c net.Conn) {
	log.Printf("New TCP connection")

	buffer := make([]byte, 1024)
	c.Read(buffer)
	log.Printf("Read tcp buffer: %s", buffer)

	err := miru.PostOneEvent()
	if err != nil {
		log.Print(err)
	}
}

func tcpMessagePump() error {
	log.Printf("Listen for tcp traffic on %s", activeMiruEnv.miruAddress)

	l, err := net.Listen("tcp", activeMiruEnv.miruAddress)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		go handleTCPConnection(c)
	}
}

func init() {
	activeMiruEnv.miruAddress = os.Getenv("MIRU_STUMPTOWN_HOST_PORT")
	if activeMiruEnv.miruAddress == "" {
		log.Print("MIRU_STUMPTOWN_HOST_PORT not present in environment. Default to :514.")
		activeMiruEnv.miruAddress = ":514"
	}
}

func main() {
	timeStart = time.Now()

	log.Print("Start admin handler")
	go func() {
		admin := mux.NewRouter()
		admin.HandleFunc("/", adminRootHandler)
		admin.HandleFunc("/ping", pingHandler)
		admin.HandleFunc("/healthcheck", healthCheckHandler)
		admin.HandleFunc("/metrics", metricsHandler)
		http.ListenAndServe(":8081", admin)
	}()

	var err error

	log.Print("Start udp handler")
	go udpMessagePump()

	log.Print("Start tcp pump")
	err = tcpMessagePump()
	if err != nil {
		log.Fatal(err)
	}
}
