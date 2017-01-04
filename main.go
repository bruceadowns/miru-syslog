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
	tcpListenAddress       string
	stumptownAddress       string
	miruStumptownIntakeURL string

	channelBufferSizeParse     int
	channelBufferSizeMiruAccum int
	channelBufferSizeMiruPost  int

	channelBufferSizeS3Accum int
	channelBufferSizeS3Post  int

	channelBufferMiruAccumBatch       int
	channelBufferMiruAccumDelayMs     int
	channelBufferMiruDelayOnSuccessMs int
	channelBufferMiruDelayOnErrorMs   int

	channelBufferS3AccumBatchBytes  int
	channelBufferS3AccumDelayMs     int
	channelBufferS3DelayOnSuccessMs int
	channelBufferS3DelayOnErrorMs   int

	awsInfo lib.AWSInfo
}

var (
	activeMiruEnv miruEnv
	sb            lib.SwitchBoard
)

func handleTCPConnection(c net.Conn) {
	log.Printf("New TCP connection: %s:%s", c.LocalAddr(), c.RemoteAddr())

	buf := bufio.NewReader(c)

	var err error
	for err == nil {
		var line []byte
		line, err = buf.ReadBytes('\n')
		if err == nil {
			p := &lib.Packet{Address: c.RemoteAddr().String(), Message: line}
			log.Printf("Read tcp buffer: %s", p)
			sb.ParseChan <- p
			sb.S3AccumChan <- p
		} else if err == io.EOF {
			if len(line) > 0 {
				log.Printf("Unexpected (ignore) buffer on EOF %s", line)
			}

			log.Print("tcp buffer EOF")
			break
		} else {
			log.Printf("Error reading bytes from channel: %s", err)
			break
		}
	}
}

func tcpMessagePump(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		if activeMiruEnv.tcpListenAddress == "" {
			log.Printf("Not listening for for tcp traffic")
			return
		}

		log.Printf("Listen for tcp traffic on %s", activeMiruEnv.tcpListenAddress)

		l, err := net.Listen("tcp", activeMiruEnv.tcpListenAddress)
		if err != nil {
			log.Printf("Error listening on tcp %s: %s", activeMiruEnv.tcpListenAddress, err)
			return
		}
		defer l.Close()

		for {
			log.Printf("Accept connections")
			c, err := l.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %s", err)
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
	activeMiruEnv.channelBufferSizeMiruAccum = lib.GetEnvInt("CHANNEL_BUFFER_SIZE_MIRU_ACCUM", 1024)
	activeMiruEnv.channelBufferSizeMiruPost = lib.GetEnvInt("CHANNEL_BUFFER_SIZE_MIRU_POST", 1024)

	activeMiruEnv.channelBufferSizeS3Accum = lib.GetEnvInt("CHANNEL_BUFFER_SIZE_S3_ACCUM", 1024)
	activeMiruEnv.channelBufferSizeS3Post = lib.GetEnvInt("CHANNEL_BUFFER_SIZE_S3_POST", 1024)

	activeMiruEnv.channelBufferMiruAccumBatch = lib.GetEnvInt("CHANNEL_BUFFER_MIRU_ACCUM_BATCH", 1024)
	activeMiruEnv.channelBufferMiruAccumDelayMs = lib.GetEnvInt("CHANNEL_BUFFER_MIRU_ACCUM_DELAY_MS", 1000)
	activeMiruEnv.channelBufferMiruDelayOnSuccessMs = lib.GetEnvInt("CHANNEL_BUFFER_MIRU_DELAY_ON_SUCCESS_MS", 500)
	activeMiruEnv.channelBufferMiruDelayOnErrorMs = lib.GetEnvInt("CHANNEL_BUFFER_MIRU_DELAY_ON_ERROR_MS", 5000)

	activeMiruEnv.channelBufferS3AccumBatchBytes = lib.GetEnvInt("CHANNEL_BUFFER_S3_ACCUM_BATCH_BYTES", 10*1024*1024)
	activeMiruEnv.channelBufferS3AccumDelayMs = lib.GetEnvInt("CHANNEL_BUFFER_S3_ACCUM_DELAY_MS", 8*60*60*1000*100)
	activeMiruEnv.channelBufferS3DelayOnSuccessMs = lib.GetEnvInt("CHANNEL_BUFFER_S3_DELAY_ON_SUCCESS_MS", 1000)
	activeMiruEnv.channelBufferS3DelayOnErrorMs = lib.GetEnvInt("CHANNEL_BUFFER_S3_DELAY_ON_ERROR_MS", 5000)

	activeMiruEnv.awsInfo = lib.AWSInfo{
		AwsRegion:          lib.GetEnvStr("AWS_REGION", "us-west-2"),
		S3Bucket:           lib.GetEnvStr("AWS_S3_BUCKET_NAME", "miru-syslog"),
		AwsAccessKeyID:     lib.GetEnvStr("AWS_ACCESS_KEY_ID", ""),
		AwsSecretAccessKey: lib.GetEnvStr("AWS_SECRET_ACCESS_KEY", "")}
}

func initChannels() {
	log.Print("Initialize channels")

	sb.MiruPostChan = lib.MiruPostChan(
		activeMiruEnv.channelBufferSizeMiruPost,
		activeMiruEnv.stumptownAddress,
		activeMiruEnv.miruStumptownIntakeURL,
		time.Millisecond*time.Duration(activeMiruEnv.channelBufferMiruDelayOnSuccessMs),
		time.Millisecond*time.Duration(activeMiruEnv.channelBufferMiruDelayOnErrorMs))

	sb.MiruAccumChan = lib.MiruAccumChan(
		activeMiruEnv.channelBufferSizeMiruAccum,
		activeMiruEnv.channelBufferMiruAccumBatch,
		time.Millisecond*time.Duration(activeMiruEnv.channelBufferMiruAccumDelayMs),
		sb.MiruPostChan)

	sb.ParseChan = lib.ParseChan(
		activeMiruEnv.channelBufferSizeParse,
		sb.MiruAccumChan)

	sb.S3PostChan = lib.S3PostChan(
		activeMiruEnv.channelBufferSizeS3Post,
		activeMiruEnv.awsInfo,
		time.Millisecond*time.Duration(activeMiruEnv.channelBufferS3DelayOnSuccessMs),
		time.Millisecond*time.Duration(activeMiruEnv.channelBufferS3DelayOnErrorMs))

	sb.S3AccumChan = lib.S3AccumChan(
		activeMiruEnv.channelBufferSizeS3Accum,
		activeMiruEnv.channelBufferS3AccumBatchBytes,
		time.Millisecond*time.Duration(activeMiruEnv.channelBufferS3AccumDelayMs),
		sb.S3PostChan)
}

func initS3() {
	log.Print("Initialize AWS S3")

	if err := lib.InitS3(activeMiruEnv.awsInfo); err != nil {
		log.Printf("Error initializing S3 (bucket): %s", err)
	}
}

func main() {
	var wg sync.WaitGroup

	go initChannels()
	go initS3()

	log.Print("Start tcp pump")
	tcpMessagePump(&wg)

	log.Print("Wait for message pump to finish")
	wg.Wait()

	log.Print("Done")
}
