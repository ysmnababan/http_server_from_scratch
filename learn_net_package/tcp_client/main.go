package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer func() { _ = conn.Close() }()

	for {
		fmt.Printf("input some text: \n")
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			panic(err)
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			panic(err)
		}
		// _ = conn.SetDeadline(time.Now().Add(5 * time.Second))

		line, err = bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("error read from conn: %v", err)
			return
		}
		fmt.Println("response: ", line)
	}
}
