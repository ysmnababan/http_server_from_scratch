package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":1323")
	if err != nil {
		panic(err)
	}
	defer func() { _ = ln.Close() }()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("client disconnecting")
			return
		}
		fmt.Print(msg)

		// if you use ` instead of " the \r or \n will be
		// treated as literal characters, so the response
		// will not be correctly formatted as an HTTP response.
		// _, _ = conn.Write([]byte(response))

		// body := "{\"message\": \"New User Created\"}"

		body := `{
		"message": "New User Created"
		}`
		response := "HTTP/1.1 201 Created\r\n" +
			fmt.Sprintf("Content-Length: %d\r\n", len(body)) +
			"Content-Type: application/json\r\n" +
			"Location: http://example.com/users/123\r\n" +
			"\r\n" +
			body
		_, _ = conn.Write([]byte(response))
	}
}
/*
		// response := "HTTP/1.1 200 OK\r\n" +
		// 	"Content-Type: text/plain\r\n" +
		// 	"Content-Length: 13\r\n" +
		// 	"\r\n" +
		// 	"Hello, World!"
*/[]
		// response := "HTTP/1.1 200 OK\r\n" +
		// 	"Content-Type: text/plain\r\n" +
		// 	"Content-Length: 13\r\n" +
		// 	"\r\n" +
		// 	"Hello, World!"
