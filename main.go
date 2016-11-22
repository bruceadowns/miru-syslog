package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/bruceadowns/miru-syslog/mako"
	"github.com/bruceadowns/miru-syslog/miru"
)

type miruEnv struct {
	miruAddress string
}

var activeMiruEnv miruEnv

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

func udpMessagePump() error {
	mako.LogInfo(fmt.Sprintf("Listen for udp traffic on %s", activeMiruEnv.miruAddress))

	pc, err := net.ListenPacket("udp", activeMiruEnv.miruAddress)
	if err != nil {
		return err
	}
	defer pc.Close()

	mako.LogInfo("Handle UDP connections")
	for {
		buffer := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			mako.LogError(err.Error())
			return err
		}

		mako.LogInfo(fmt.Sprintf("Read udp buffer from %s: %s", addr, buffer[:n]))

		err = miru.PostManyEvents()
		if err != nil {
			mako.LogInfo(err.Error())
		}
	}
}

func handleTCPConnection(c net.Conn) {
	mako.LogInfo("New TCP connection")

	buffer := make([]byte, 1024)
	n, err := c.Read(buffer)
	if err != nil {
		mako.LogError(err.Error())
		return
	}

	mako.LogInfo(fmt.Sprintf("Read tcp buffer: %s", buffer[:n]))

	err = miru.PostOneEvent()
	if err != nil {
		mako.LogInfo(err.Error())
	}
}

func tcpMessagePump() error {
	mako.LogInfo(fmt.Sprintf("Listen for tcp traffic on %s", activeMiruEnv.miruAddress))

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
		mako.LogInfo("MIRU_STUMPTOWN_HOST_PORT not present in environment. Default to :514.")
		activeMiruEnv.miruAddress = ":514"
	}
}

func main() {
	mako.LogInfo("Start admin handler")
	go func() {
		admin := mux.NewRouter()
		admin.HandleFunc("/", adminRootHandler)
		admin.HandleFunc("/ping", pingHandler)
		admin.HandleFunc("/healthcheck", healthCheckHandler)
		admin.HandleFunc("/metrics", metricsHandler)
		http.ListenAndServe(":8081", admin)
	}()

	var err error

	mako.LogInfo("Start udp handler")
	go udpMessagePump()

	mako.LogInfo("Start tcp pump")
	err = tcpMessagePump()
	if err != nil {
		mako.LogError(err.Error())
	}
}
