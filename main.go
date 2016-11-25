package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"

	"github.com/bruceadowns/miru-syslog/miru"
	// "github.com/jeromer/syslogparser"
	// "github.com/jeromer/syslogparser/rfc3164"
)

type miruEnv struct {
	listenAddress    string
	stumptownAddress string
}

var activeMiruEnv miruEnv

func udpMessagePump() error {
	log.Printf("Listen for udp traffic on %s", activeMiruEnv.listenAddress)

	pc, err := net.ListenPacket("udp", activeMiruEnv.listenAddress)
	if err != nil {
		return err
	}
	defer pc.Close()

	log.Print("Handle UDP connections")
	for {
		buffer := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			log.Print(err)
			return err
		}

		log.Printf("Read udp buffer from %s: %s", addr, buffer[:n])

		err = miru.PostManyEvents(activeMiruEnv.stumptownAddress)
		if err != nil {
			log.Print(err)
		}
	}
}

func handleTCPConnection(c net.Conn) {
	log.Print("New TCP connection")

	buf := bufio.NewReader(c)
	var err error

	for err == nil {
		var line []byte
		line, err = buf.ReadBytes('\n')
		if err == nil {
			log.Printf("Send tcp buffer to parse channel: %s", line)
		} else if err == io.EOF {
			log.Printf("Send last tcp buffer to parse channel: %s", line)
		} else {
			log.Print(err)
			return
		}
	}
}

func tcpMessagePump() error {
	log.Printf("Listen for tcp traffic on %s", activeMiruEnv.listenAddress)

	l, err := net.Listen("tcp", activeMiruEnv.listenAddress)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		log.Printf("Accept connections")
		c, err := l.Accept()
		if err != nil {
			return err
		}

		go handleTCPConnection(c)
	}
}

func init() {
	activeMiruEnv.listenAddress = os.Getenv("MIRU_SYSLOG_HOST_PORT")
	if activeMiruEnv.listenAddress == "" {
		log.Print("MIRU_SYSLOG_HOST_PORT not present in environment. Default to :514.")
		activeMiruEnv.listenAddress = ":514"
	}

	activeMiruEnv.stumptownAddress = os.Getenv("MIRU_STUMPTOWN_HOST_PORT")
	if activeMiruEnv.stumptownAddress == "" {
		log.Print("MIRU_STUMPTOWN_HOST_PORT not present in environment.")
	}
}

func main() {
	log.Print("Start udp handler")
	go udpMessagePump()

	log.Print("Start tcp pump")
	err := tcpMessagePump()
	if err != nil {
		log.Fatal(err)
	}
}
