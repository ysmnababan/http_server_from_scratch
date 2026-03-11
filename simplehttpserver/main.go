package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type UserResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

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
		go handleWithReq(conn)
		// go handleSimple(conn)
	}
}

func handleSimple(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			sendResponse(conn)
			log.Println("CLIENT DISCONNECTING...")
			return
		}
		fmt.Print(">>", msg)
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
		fmt.Print(">>", msg)

		// if you use ` instead of " the \r or \n will be
		// treated as literal characters, so the response
		// will not be correctly formatted as an HTTP response.

		b := struct {
			Message string `json:"message"`
			User    struct {
				Id        int    `json:"id"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Email     string `json:"email"`
			}
		}{
			Message: "New User Created",
			User: struct {
				Id        int    `json:"id"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Email     string `json:"email"`
			}{
				Id:        10,
				FirstName: "Test",
				LastName:  "Name",
				Email:     "exampleemail",
			},
		}

		bstring, _ := json.Marshal(&b)
		body := string(bstring)
		log.Printf("Response body: %s\n>>>>>>\n", body)
		response := "HTTP/1.1 201 Created\r\n" +
			fmt.Sprintf("Content-Length: %d\r\n", len(body)) +
			"Content-Type: application/json\r\n" +
			"Location: http://example.com/users/123\r\n" +
			"\r\n" +
			body
		_, _ = conn.Write([]byte(response))
	}
}

func sendResponse(conn net.Conn) {
	b := UserResponse{
		Message: "New User Created",
		User: User{
			FirstName: "first",
			LastName:  "last",
			ID:        321,
			Email:     "example@student.co.id",
		},
	}
	bstring, _ := json.Marshal(&b)
	body := string(bstring)
	// log.Printf("Response body: %s\n>>>>>>\n", body)
	response := "HTTP/1.1 201 Created\r\n" +
		fmt.Sprintf("Content-Length: %d\r\n", len(body)) +
		"Content-Type: application/json\r\n" +
		"Location: http://example.com/users/123\r\n" +
		"\r\n" +
		body
	_, _ = conn.Write([]byte(response))
}
