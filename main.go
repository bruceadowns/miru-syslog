package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"sync"

	"github.com/bruceadowns/miru-syslog/comm"
	"github.com/bruceadowns/miru-syslog/common"
	"github.com/bruceadowns/miru-syslog/miru"
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
	parseChan     chan *comm.Packet
	postChan      chan *miru.LogEvent
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
					parseChan <- &comm.Packet{Address: addr, Message: line}
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
				p := &comm.Packet{Address: c.RemoteAddr(), Message: line}
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
	activeMiruEnv.tcpListenAddress = common.GetEnvStr("MIRU_SYSLOG_TCP_ADDR_PORT", "")
	activeMiruEnv.udpListenAddress = common.GetEnvStr("MIRU_SYSLOG_UDP_ADDR_PORT", "")
	activeMiruEnv.stumptownAddress = common.GetEnvStr("MIRU_STUMPTOWN_ADDR_PORT", "")
	activeMiruEnv.miruStumptownIntakeURL = common.GetEnvStr("MIRU_STUMPTOWN_INTAKE_URL", "/miru/stumptown/intake")
	activeMiruEnv.channelBufferSizeParse = common.GetEnvInt("CHANNEL_BUFFER_SIZE_PARSE", 1024)
	activeMiruEnv.channelBufferSizePost = common.GetEnvInt("CHANNEL_BUFFER_SIZE_POST", 1024)
	activeMiruEnv.udpReceiveBufferSize = common.GetEnvInt("UDP_RECEIVE_BUFFER_SIZE", 2*1024*1024)

	postChan = comm.PostChan(
		activeMiruEnv.channelBufferSizePost,
		activeMiruEnv.stumptownAddress,
		activeMiruEnv.miruStumptownIntakeURL)
	parseChan = comm.ParseChan(activeMiruEnv.channelBufferSizeParse, postChan)
}

func main() {
	log.Print("hello world")

	var wg sync.WaitGroup

	log.Print("Start udp handler")
	udpMessagePump(&wg)

	log.Print("Start tcp pump")
	tcpMessagePump(&wg)

	log.Print("Wait for both message pumps to finish")
	wg.Wait()
}
