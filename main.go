package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/bruceadowns/miru-syslog/lib"
)

type miruEnv struct {
	tcpListenAddress          string
	stumptownAddress          string
	miruStumptownIntakeURL    string
	channelBufferSizeParse    int
	channelBufferSizeAccum    int
	channelBufferSizePost     int
	channelBufferAccumBatch   int
	channelBufferAccumDelayMs int
}

var (
	activeMiruEnv miruEnv

	parseChan chan lib.Packet
	accumChan chan lib.LogEvent
	postChan  chan lib.LogEvents
)

func handleTCPConnection(c net.Conn) {
	log.Print("New TCP connection")

	buf := bufio.NewReader(c)

	var err error
	for err == nil {
		var line []byte
		line, err = buf.ReadBytes('\n')
		if err == nil {
			p := &lib.Packet{Address: c.RemoteAddr().String(), Message: line}
			log.Printf("Read tcp buffer: %s", p)
			parseChan <- *p
		} else if err == io.EOF {
			if len(line) > 0 {
				log.Fatal("Unexpected buffer len with EOF")
			}

			log.Print("tcp buffer EOF")
			break
		} else {
			log.Print(err)
			break
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
	activeMiruEnv.stumptownAddress = lib.GetEnvStr("MIRU_STUMPTOWN_ADDR_PORT", "")
	activeMiruEnv.miruStumptownIntakeURL = lib.GetEnvStr("MIRU_STUMPTOWN_INTAKE_URL", "/miru/stumptown/intake")
	activeMiruEnv.channelBufferSizeParse = lib.GetEnvInt("CHANNEL_BUFFER_SIZE_PARSE", 1024)
	activeMiruEnv.channelBufferSizeAccum = lib.GetEnvInt("CHANNEL_BUFFER_SIZE_ACCUM", 1024)
	activeMiruEnv.channelBufferSizePost = lib.GetEnvInt("CHANNEL_BUFFER_SIZE_POST", 1024)
	activeMiruEnv.channelBufferAccumBatch = lib.GetEnvInt("CHANNEL_BUFFER_ACCUM_BATCH", 1000)
	activeMiruEnv.channelBufferAccumDelayMs = lib.GetEnvInt("CHANNEL_BUFFER_ACCUM_DELAY_MS", 100)

	postChan = lib.PostChan(
		activeMiruEnv.channelBufferSizePost,
		activeMiruEnv.stumptownAddress,
		activeMiruEnv.miruStumptownIntakeURL)
	accumChan = lib.AccumChan(
		activeMiruEnv.channelBufferSizeAccum,
		activeMiruEnv.channelBufferAccumBatch,
		time.Millisecond*time.Duration(activeMiruEnv.channelBufferAccumDelayMs),
		postChan)
	parseChan = lib.ParseChan(
		activeMiruEnv.channelBufferSizeParse,
		accumChan)
}

func main() {
	var wg sync.WaitGroup

	log.Print("Start tcp pump")
	tcpMessagePump(&wg)

	log.Print("Wait for message pump to finish")
	wg.Wait()

	log.Print("Done")
}
