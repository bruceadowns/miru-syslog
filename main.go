package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/bruceadowns/miru-syslog/mako"
	"github.com/bruceadowns/miru-syslog/miru"
)

const listenAddr = ":514"

var timeStart time.Time

func adminRootHandler(w http.ResponseWriter, r *http.Request) {
	common.LogInfo(fmt.Sprintf("%v", r))
	w.Write([]byte("Hello Admin Root"))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	common.LogInfo(fmt.Sprintf("%v", r))
	common.DataDogClient.Count("ping", 1, nil, 0)

	w.Write([]byte("pong"))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	common.LogInfo(fmt.Sprintf("%v", r))
	w.Write([]byte("Hello Health Check"))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	common.LogInfo(fmt.Sprintf("%v", r))
	w.Write([]byte("Hello Metrics"))
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	common.LogInfo(fmt.Sprintf("%v", r))
	common.DataDogClient.Count("serviceHandler", 1, nil, 0)
	common.DataDogClient.SimpleEvent("serviceHandler", "Service Handler was called")
	w.Write([]byte("Hello Root"))
}

func udpMessagePump() error {
	log.Printf("Listen for udp traffic on %s", listenAddr)

	pc, err := net.ListenPacket("udp", listenAddr)
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
	log.Printf("Listen for tcp traffic on %s", listenAddr)

	l, err := net.Listen("tcp", listenAddr)
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
