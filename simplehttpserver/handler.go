package main

import (
	"bufio"
	"log"
	"net"
)

func handleWithReq(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	reader := bufio.NewReader(conn)
	httpReq := NewRequest()

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("client disconnected...", err)
			return
		}
		httpReq.parse(msg)
		if httpReq.parsingStatus == completeStatus {
			httpReq.debug()
			sendResponse(conn)
		}
	}
}
