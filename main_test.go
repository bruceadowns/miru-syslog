package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"testing"
)

func TestTcpClient(t *testing.T) {
	dockerAddr := os.Getenv("MIRU_SYSLOG_TCP_ADDR_PORT")
	if len(dockerAddr) == 0 {
		t.Skip("MIRU_SYSLOG_TCP_ADDR_PORT not found.")
	}

	fmt.Printf("Connect to tcp server at %s\n", dockerAddr)

	conn, err := net.Dial("tcp", dockerAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	message := "Hello from client"
	fmt.Printf("Simple write %s\n", message)
	conn.Write([]byte(message))
}

func TestUdpClient(t *testing.T) {
	dockerAddr := os.Getenv("MIRU_SYSLOG_UDP_ADDR_PORT")
	if len(dockerAddr) == 0 {
		t.Skip("MIRU_SYSLOG_UDP_ADDR_PORT not found.")
	}

	fmt.Printf("Connect to udp server at %s\n", dockerAddr)

	conn, err := net.Dial("udp", dockerAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	message := "Hello from client"
	fmt.Printf("Simple write %s\n", message)
	conn.Write([]byte(message))
}

func handleConnection(c net.Conn) {
	fmt.Printf("New connection\n")
	buffer := make([]byte, 1024)
	c.Read(buffer)
	fmt.Printf("Read %s\n", buffer)
}

func TestTcpServer(t *testing.T) {
	dockerAddr := os.Getenv("MIRU_SYSLOG_TCP_ADDR_PORT_SERVER")
	if len(dockerAddr) == 0 {
		t.Skip("MIRU_SYSLOG_TCP_ADDR_PORT_SERVER not found.")
	}

	fmt.Printf("Listen for tcp traffic on %s\n", dockerAddr)
	l, err := net.Listen("tcp", dockerAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		go handleConnection(c)
	}
}

func TestUdpServer(t *testing.T) {
	dockerAddr := os.Getenv("MIRU_SYSLOG_UDP_ADDR_PORT_SERVER")
	if len(dockerAddr) == 0 {
		t.Skip("MIRU_SYSLOG_UDP_ADDR_PORT_SERVER not found.")
	}

	fmt.Printf("Listen for udp traffic on %s\n", dockerAddr)
	pc, err := net.ListenPacket("udp", dockerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buffer := make([]byte, 1024)
		pc.ReadFrom(buffer)
		fmt.Printf("Read %s\n", buffer)
	}
}
