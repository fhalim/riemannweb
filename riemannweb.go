package main

import (
	rm "github.com/amir/raidman/proto"
	pb "code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"net"
	"fmt"
	//"code.google.com/p/go.net/ipv4"
	"os"
	"io"
)

func main() {
	fmt.Printf("Hello world!")
	ln, err := net.Listen("tcp4", "0.0.0.0:9999")
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}
	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting request: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(c)
	}
	select {}
}
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	len := readLength(conn)
	buf := make([]byte, len)
	// Read the incoming connection into the buffer.
	err := readFully(conn, buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	evt := new(rm.Event)
	err = pb.Unmarshal(buf, evt)
	if(err != nil){
		fmt.Println("Could not unmarshall request: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Message received: ", evt)
	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}
func readLength(conn net.Conn) uint32 {
    lenbuf := make([]byte, 4)
	err := readFully(conn, lenbuf)
	if err != nil {
		fmt.Println("Error reading length:", err.Error())
	}
	return binary.BigEndian.Uint32(lenbuf)
}

func readFully(r io.Reader, p []byte) error {
	for len(p) > 0 {
		n, err := r.Read(p)
		p = p[n:]
		if err != nil {
			return err
		}
	}
	return nil
}
