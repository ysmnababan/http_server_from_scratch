# Building HTTP Server from Scratch — Learning Summary

Repo for implementing HTTP server from scratch in Go
A summary of concepts covered before diving into echo framework internals.

---

## 1. TCP — The Foundation

Everything starts at TCP. HTTP is just text sent over a TCP connection. Go's `net` package exposes this directly.

```go
listener, _ := net.Listen("tcp", ":8080")
conn, _ := listener.Accept()
```

Key insight: TCP is a **raw byte stream** — it has no concept of messages or boundaries. Every protocol built on top of TCP must define its own message boundaries.

---

## 2. Reading from a TCP Connection

### The wrong way — `bytes.Buffer` with `ReadFrom`

```go
buf := &bytes.Buffer{}
buf.ReadFrom(conn) // blocks until EOF — never returns while client is connected
```

`ReadFrom` blocks until the connection closes. You never get to process the message in real-time.

### The manual way — byte by byte with `bytes.Buffer`

```go
buf := &bytes.Buffer{}
b := make([]byte, 1)
for {
    conn.Read(b)
    buf.Write(b)
    if b[0] == '\n' {
        // message complete
        conn.Write(buf.Bytes())
        buf.Reset()
    }
}
```

Works, but reads 1 byte per syscall — extremely inefficient.

### The correct way — `bufio.Reader`

```go
reader := bufio.NewReader(conn)
msg, _ := reader.ReadString('\n') // returns as soon as delimiter is found
```

`bufio.Reader` solves two problems at once:

- **Syscall efficiency** — reads 4096 bytes per syscall into an internal buffer, serves your reads from memory
- **Message boundaries** — `ReadString`, `ReadLine`, `ReadBytes` let you consume the stream in logical chunks

|                    | `bytes.Buffer`    | `bufio.Reader`                 |
| ------------------ | ----------------- | ------------------------------ |
| Purpose            | In-memory storage | Buffered reading from a source |
| Wraps              | Nothing           | Any `io.Reader`                |
| Syscall efficient  | No                | Yes — reads in chunks          |
| Message boundaries | No                | Yes                            |

---

## 3. HTTP — Protocol Structure

HTTP/1.1 is a text protocol with strict formatting rules.

### Request structure

```
GET /hello?name=world HTTP/1.1\r\n
Host: localhost:8080\r\n
Content-Type: application/json\r\n
\r\n
{"key": "value"}
```

### Response structure

```
HTTP/1.1 200 OK\r\n
Content-Type: application/json\r\n
Content-Length: 16\r\n
\r\n
{"key": "value"}
```

### `\r\n` — Why it matters

| Character | Name            | Decimal |
| --------- | --------------- | ------- |
| `\r`      | Carriage Return | 13      |
| `\n`      | Line Feed       | 10      |

HTTP mandates `\r\n` (CRLF) as the line terminator — not just `\n`. This was standardized so the protocol works consistently across all operating systems (Windows uses `\r\n`, Unix uses `\n`).

The **blank `\r\n` line** is the separator between headers and body. Your parser must detect this to know when headers end.

### Parsing headers manually

```go
line, _ := reader.ReadString('\n')
line = strings.TrimRight(line, "\r\n") // strip both characters
if line == "" {
    // headers done, body starts here
}
parts := strings.Split(line, " ")
method, path, version := parts[0], parts[1], parts[2]
```

### Three ways protocols define message boundaries

| Strategy      | Example                            | How                                 |
| ------------- | ---------------------------------- | ----------------------------------- |
| Delimiter     | HTTP headers, Redis simple strings | Read until `\r\n` or `\r\n\r\n`     |
| Length-prefix | gRPC, HTTP body                    | Read `Content-Length` bytes exactly |
| Fixed-size    | Some legacy protocols              | Always read N bytes                 |

HTTP uses **both** — headers are delimiter-based, body is length-prefixed via `Content-Length`.

---

## 4. Writing a Valid HTTP Response

### Go string literals — a critical gotcha

```go
// ❌ backticks = raw string — \r\n are literal characters, not control codes
`HTTP/1.1 200 OK\r\n`

// ✅ double quotes = interpreted string — \r\n become actual CR+LF bytes
"HTTP/1.1 200 OK\r\n"
```

### Correct response construction

```go
body := `{"message": "New User Created"}`
response := "HTTP/1.1 201 Created\r\n" +
    "Content-Type: application/json\r\n" +
    "Location: http://example.com/users/123\r\n" +
    fmt.Sprintf("Content-Length: %d\r\n", len(body)) +
    "\r\n" +
    body
conn.Write([]byte(response))
```

### Why `Content-Length` is required

Without it, the client doesn't know when the body ends. It keeps waiting for more bytes — Postman will show the request as "still connected" and never display the response.

### HTTP/1.1 keep-alive behavior

HTTP/1.1 keeps connections open by default after a response (keep-alive). To close immediately after responding:

```go
"Connection: close\r\n"
```

---

## 5. `net/http` — What It Does for You

Go's standard library wraps all of the above. Every line you wrote manually maps to something in `net/http`:

| Your manual code                       | `net/http` equivalent                  |
| -------------------------------------- | -------------------------------------- |
| `net.Listen` + `Accept` loop           | `http.ListenAndServe`                  |
| `bufio.NewReader(conn)`                | done internally in `conn.serve()`      |
| Parsing request line + headers         | `readRequest` in `server.go`           |
| Writing status line + headers + `\r\n` | `ResponseWriter.WriteHeader` + `Write` |
| Connection keep-alive loop             | `conn.serve()` read loop               |

The entire framework contract is one interface:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Echo, Gin, Chi, Fiber — they all just implement `ServeHTTP`. Everything else is routing and middleware on top.

### Go 1.22 — Method routing in `ServeMux`

Before Go 1.22 you had to check the HTTP method manually inside your handler. From 1.22 onwards, `ServeMux` accepts method + path patterns directly:

```go
// pre 1.22
http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
})

// 1.22+
http.HandleFunc("PUT /test", handler)
// GET /test now automatically returns 405 with correct Allow header
```

---

## 6. Roadmap to `net/http` Internals

Read the source in this order — each step builds on the last:

1. `ListenAndServe` → `Server.Serve` — recognize your TCP accept loop
2. `conn.serve()` — see `bufio.NewReader` and the per-connection goroutine
3. `readRequest` — map it to your manual parser
4. `http.Handler` interface — understand the framework handoff point
5. `response` struct — see your manual `conn.Write` abstracted into `ResponseWriter`

### Success indicators

You're ready to read echo's source when you can answer these without looking:

- Where exactly does `net/http` hand control to your handler?
- Why does calling `w.Write()` before `w.WriteHeader()` still produce a valid response?
- What happens if you never call `w.WriteHeader()` at all?
- Why does `http.ListenAndServe` spawn a goroutine per connection?
- What is `http.ServeMux` doing that a map-based router also does?
- Why does reading `r.Body` twice not work by default?
