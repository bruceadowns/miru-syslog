package main

import (
	"fmt"
	"log"
	"net"
	"testing"
)

func TestMain(t *testing.T) {
	t.Log("Testing main")
}

func TestTcpClient(t *testing.T) {
	fmt.Printf("Connect to tcp server at %s\n", listenAddr)

	conn, err := net.Dial("tcp", listenAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	message := "Hello from client"
	fmt.Printf("Simple write %s\n", message)
	conn.Write([]byte(message))
}

func TestUdpClient(t *testing.T) {
	fmt.Printf("Connect to udp server at %s\n", listenAddr)

	conn, err := net.Dial("udp", listenAddr)
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
	fmt.Printf("Listen for tcp traffic on %s\n", listenAddr)
	l, err := net.Listen("tcp", listenAddr)
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
	fmt.Printf("Listen for udp traffic on %s\n", listenAddr)
	pc, err := net.ListenPacket("udp", listenAddr)
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
