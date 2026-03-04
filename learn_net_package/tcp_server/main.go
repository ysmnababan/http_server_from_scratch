package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer func() { _ = ln.Close() }()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		// go handleConnWithBuffer(conn)
		go handleConn(conn)

	}
}

// use conn.Read to read each byte,
// and write to buffer until read a newline character,
// then write the buffer content back to client.
// This approach is no efficient because it reads one byte at a time,
// which can lead to many system calls and increased latency,
// especially for larger messages.
func handleConnWithBuffer(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	buf := &bytes.Buffer{}
	b := make([]byte, 1)
	for {
		_, err := conn.Read(b)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		_, _ = buf.Write(b)
		if b[0] == '\n' {
			_, _ = conn.Write([]byte("Echo: " + buf.String()))
			buf.Reset()
		}
	}
}

// same functionality as handleConnWithBuffer, but uses bufio.Reader to read lines of text,
// this approach is more efficient because bufio.Rader reads larger chuncks of data at a time,
// reducing the number of system calls
func handleConn(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		fmt.Printf("read message: %s", msg)
		_, _ = conn.Write([]byte(msg))
	}
}
