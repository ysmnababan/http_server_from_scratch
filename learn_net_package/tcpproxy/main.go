package main

import (
	"io"
	"net"
)

// forward copies data from src to dst, and closes both connections when done.
//
// TCP is full-duplex, meaning data flows in both directions simultaneously.
// Because io.Copy blocks until the src connection is closed or errors,
// we need TWO goroutines — one per direction — so both can run concurrently.
//
// Data path (kernel → userspace → kernel):
//
//	src.RecvBuffer → [io.Copy 32KB userspace buf] → dst.SendBuffer
//
// Each socket has two independent kernel buffers:
//
//	RecvBuffer: filled by the remote peer, drained by Read()
//	SendBuffer: filled by Write(), drained by OS when sending to remote peer
//
// So io.Copy(dst, src) means:
//
//	"read bytes from src's recv buffer, write into dst's send buffer"
//	 the OS handles the actual network transmission on both ends.
//
// Why close BOTH connections on return?
//
//	When one side disconnects, io.Copy returns on that goroutine.
//	But the other goroutine is still blocked on its own io.Copy.
//	Closing both connections here forces the other goroutine to unblock,
//	because its Read() will return an error on the now-closed connection.
//	Without this, the other goroutine leaks forever.
//
// Note: closing the same connection twice is safe in Go — net.Conn.Close()
// is idempotent. The second close just returns a ignored error.
func forward(dst, src net.Conn) {
	defer func() { _ = src.Close() }()
	defer func() { _ = dst.Close() }()
	_, _ = io.Copy(dst, src)
}

// main listens on :8080 and proxies each connection to example.com:80.
//
// For each client connection, we create ONE connection to the remote server,
// then spawn TWO goroutines to handle both directions of the TCP stream.
//
// There is only ever 1 conn and 1 server per client — both goroutines
// share the same two socket objects, just reading/writing different buffers.
//
// Traffic flow:
//
//	CLIENT                  PROXY                    EXAMPLE.COM
//	   │                      │                           │
//	   │   ── request ──►     │                           │
//	   │               conn.Read()                        │
//	   │                      │    ── request ──►         │
//	   │               server.Write()                     │
//	   │                      │                           │
//	   │                      │    ◄── response ──        │
//	   │               server.Read()                      │
//	   │   ◄── response ──    │                           │
//	   │               conn.Write()                       │
//	   │                      │                           │
//
// Goroutine ownership:
//
//	G1 forward(server, conn): CLIENT ──────────────────────────► EXAMPLE.COM
//	                               [conn.Read] ── [server.Write]
//
//	G2 forward(conn, server): CLIENT ◄────────────────────────── EXAMPLE.COM
//	                               [conn.Write] ── [server.Read]
//
// Socket buffer layout (per socket):
//
//	conn socket                        server socket
//	┌────────────────────┐             ┌────────────────────┐
//	│ recv buf (read)    │ ◄─ client   │ recv buf (read)    │ ◄─ example.com
//	│ send buf (write)   │ ─► client   │ send buf (write)   │ ─► example.com
//	└────────────────────┘             └────────────────────┘
//	  G2 writes here                     G1 writes here
//	  G1 reads here                      G2 reads here
//
// G1 and G2 never touch the same buffer simultaneously,
// so there is no data race even though both goroutines share the same sockets.
func main() {
	ln, _ := net.Listen("tcp", ":8080")
	defer ln.Close()

	for {
		// conn: the incoming client connection
		conn, _ := ln.Accept()

		// server: our outbound connection to the remote
		server, _ := net.Dial("tcp", "example.com:80")
		// NOTE: if Dial fails, server is nil and forward will panic.
		// Always check errors in production.

		go forward(server, conn) // G1: client → example.com
		go forward(conn, server) // G2: example.com → client
	}
}
