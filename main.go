package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"sync"

	"github.com/bruceadowns/miru-syslog/lib"
)

type miruEnv struct {
	tcpListenAddress       string
	udpListenAddress       string
	stumptownAddress       string
	miruStumptownIntakeURL string
	channelBufferSizeParse int
	channelBufferSizePost  int
	udpReceiveBufferSize   int
}

var (
	activeMiruEnv miruEnv
	parseChan     chan *lib.Packet
	postChan      chan *lib.LogEvent
)

func udpMessagePump(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		if len(activeMiruEnv.udpListenAddress) == 0 {
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
			buffer := make([]byte, activeMiruEnv.udpReceiveBufferSize)
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				log.Print(err)
				return
			}

			log.Printf("Read udp buffer from %s", addr)
			for _, line := range bytes.Split(buffer[:n], []byte{'\n'}) {
				if len(line) > 0 {
					parseChan <- &lib.Packet{Address: addr, Message: line}
				}
			}
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
		if err == nil || err == io.EOF {
			if len(line) > 0 {
				p := &lib.Packet{Address: c.RemoteAddr(), Message: line}
				log.Printf("Read tcp buffer from: %s", p.Address)
				parseChan <- p
			}
		} else {
			log.Print(err)
			return
		}
	}
}

func tcpMessagePump(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		if len(activeMiruEnv.tcpListenAddress) == 0 {
			log.Printf("Not listening for for tcp traffic")
			return
		}

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
	activeMiruEnv.tcpListenAddress = lib.GetEnvStr("MIRU_SYSLOG_TCP_ADDR_PORT", "")
	activeMiruEnv.udpListenAddress = lib.GetEnvStr("MIRU_SYSLOG_UDP_ADDR_PORT", "")
	activeMiruEnv.stumptownAddress = lib.GetEnvStr("MIRU_STUMPTOWN_ADDR_PORT", "")
	activeMiruEnv.miruStumptownIntakeURL = lib.GetEnvStr("MIRU_STUMPTOWN_INTAKE_URL", "/miru/stumptown/intake")
	activeMiruEnv.channelBufferSizeParse = lib.GetEnvInt("CHANNEL_BUFFER_SIZE_PARSE", 1024)
	activeMiruEnv.channelBufferSizePost = lib.GetEnvInt("CHANNEL_BUFFER_SIZE_POST", 1024)
	activeMiruEnv.udpReceiveBufferSize = lib.GetEnvInt("UDP_RECEIVE_BUFFER_SIZE", 2*1024*1024)

	postChan = lib.PostChan(
		activeMiruEnv.channelBufferSizePost,
		activeMiruEnv.stumptownAddress,
		activeMiruEnv.miruStumptownIntakeURL)
	parseChan = lib.ParseChan(activeMiruEnv.channelBufferSizeParse, postChan)
}

func main() {
	var wg sync.WaitGroup

	log.Print("Start udp handler")
	udpMessagePump(&wg)

	log.Print("Start tcp pump")
	tcpMessagePump(&wg)

	log.Print("Wait for both message pumps to finish")
	wg.Wait()
}
