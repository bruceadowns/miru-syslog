package main

import (
	"fmt"
	"log"
	"net"
	"testing"
)

const dockerAddr = "192.168.99.100:514"

func TestMain(t *testing.T) {
	t.Log("Testing main")
}

func TestTcpClient(t *testing.T) {
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
