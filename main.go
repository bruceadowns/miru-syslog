package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/bruceadowns/miru-syslog/comm"
	// "github.com/bruceadowns/miru-syslog/miru"
	// "github.com/jeromer/syslogparser"
	// "github.com/jeromer/syslogparser/rfc3164"
)

type miruEnv struct {
	tcpListenAddress       string
	udpListenAddress       string
	stumptownAddress       string
	miruStumptownIntakeURL string
	channelBufferSizeParse int
}

var (
	activeMiruEnv miruEnv
	parseChan     chan *comm.Packet
)

func udpMessagePump(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		if activeMiruEnv.udpListenAddress == "" {
			log.Printf("Not listening for for udp traffic")
			return
		}

		log.Printf("Listen for udp traffic on %s", activeMiruEnv.udpListenAddress)

		pc, err := net.ListenPacket("udp", activeMiruEnv.udpListenAddress)
		if err != nil {
			return
		}
		defer pc.Close()

		log.Print("Handle UDP connections")
		for {
			buffer := make([]byte, 1024)
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				log.Print(err)
				return
			}

			p := &comm.Packet{Address: addr, Message: buffer[:n]}
			log.Printf("Read udp buffer from %s", p)
			parseChan <- p
		}
	}()
}

func handleTCPConnection(c net.Conn) {
	log.Print("New TCP connection")

	buf := bufio.NewReader(c)
	var err error
	c.LocalAddr()

	for err == nil {
		var line []byte
		line, err = buf.ReadBytes('\n')
		if err == nil {
			p := &comm.Packet{Address: c.RemoteAddr(), Message: line}
			log.Printf("Send tcp buffer to parse channel: %s", p)
			parseChan <- p
		} else if err == io.EOF {
			p := &comm.Packet{Address: c.RemoteAddr(), Message: line}
			log.Printf("Send tcp buffer to parse channel: %s", p)
			parseChan <- p
			log.Printf("Send last tcp buffer to parse channel: %s", line)
		} else {
			log.Print(err)
			return
		}
	}
}

func tcpMessagePump(wg *sync.WaitGroup) {
	wg.Add(1)

	if activeMiruEnv.tcpListenAddress == "" {
		log.Printf("Not listening for for tcp traffic")
		return
	}

	go func() {
		defer wg.Done()

		log.Printf("Listen for tcp traffic on %s", activeMiruEnv.tcpListenAddress)

		l, err := net.Listen("tcp", activeMiruEnv.tcpListenAddress)
		if err != nil {
			log.Print(err)
			return
		}
		defer l.Close()

		for {
			log.Printf("Accept connections")
			c, err := l.Accept()
			if err != nil {
				log.Print(err)
				return
			}

			go handleTCPConnection(c)
		}
	}()
}

func init() {
	activeMiruEnv.tcpListenAddress = os.Getenv("MIRU_SYSLOG_TCP_ADDR_PORT")
	if activeMiruEnv.tcpListenAddress == "" {
		log.Print("MIRU_SYSLOG_TCP_ADDR_PORT not present in environment. Not listening for tcp traffic.")
	}
	log.Printf("MIRU_SYSLOG_TCP_ADDR_PORT set to %s.", activeMiruEnv.tcpListenAddress)

	activeMiruEnv.udpListenAddress = os.Getenv("MIRU_SYSLOG_UDP_ADDR_PORT")
	if activeMiruEnv.udpListenAddress == "" {
		log.Print("MIRU_SYSLOG_UDP_ADDR_PORT not present in environment. Not listening for udp traffic.")
	}
	log.Printf("MIRU_SYSLOG_UDP_ADDR_PORT set to %s.", activeMiruEnv.udpListenAddress)

	activeMiruEnv.stumptownAddress = os.Getenv("MIRU_STUMPTOWN_ADDR_PORT")
	if activeMiruEnv.stumptownAddress == "" {
		log.Print("MIRU_STUMPTOWN_ADDR_PORT not present in environment.")
	}
	log.Printf("MIRU_STUMPTOWN_ADDR_PORT set to %s.", activeMiruEnv.stumptownAddress)

	activeMiruEnv.miruStumptownIntakeURL = "/miru/stumptown/intake"
	sIntake := os.Getenv("MIRU_STUMPTOWN_INTAKE_URL")
	if sIntake == "" {
		log.Print("MIRU_STUMPTOWN_INTAKE_URL not present in environment.")
	} else {
		activeMiruEnv.miruStumptownIntakeURL = sIntake
	}
	log.Printf("MIRU_STUMPTOWN_INTAKE_URL set to %s.", activeMiruEnv.miruStumptownIntakeURL)

	activeMiruEnv.channelBufferSizeParse = 100
	sSize := os.Getenv("CHANNEL_BUFFER_SIZE_PARSE")
	if sSize == "" {
		log.Printf("CHANNEL_BUFFER_SIZE_PARSE not present in environment.")
	} else {
		iSize, err := strconv.Atoi(sSize)
		if err == nil {
			activeMiruEnv.channelBufferSizeParse = iSize
		} else {
			log.Printf("CHANNEL_BUFFER_SIZE_PARSE not numeric %s.", sSize)
		}
	}
	log.Printf("CHANNEL_BUFFER_SIZE_PARSE set to %d.", activeMiruEnv.channelBufferSizeParse)
}

func main() {
	parseChan = comm.ParseChan(activeMiruEnv.channelBufferSizeParse)

	var wg sync.WaitGroup

	log.Print("Start udp handler")
	udpMessagePump(&wg)

	log.Print("Start tcp pump")
	tcpMessagePump(&wg)

	log.Print("Wait for both message pumps to finish")
	wg.Wait()
}
