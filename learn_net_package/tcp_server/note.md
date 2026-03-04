**`bytes.Buffer` — in-memory storage**

It's a growable byte slice. It has no concept of an external data source. You write data _into_ it, then read data _out of_ it. It's a container.

```go
buf := bytes.NewBuffer([]byte("hello"))
buf.Write([]byte(" world"))
fmt.Println(buf.String()) // "hello world"
```

When you call `buf.ReadFrom(conn)`, you're telling the buffer to pull everything from `conn` into itself — which as we discussed, blocks until EOF.

**`bufio.Reader` — buffered reading from a source**

It wraps an `io.Reader` (like a TCP connection) and adds an internal byte buffer on top of it. The key difference: it reads from the source in **chunks** into its internal buffer, then lets you consume that buffer piece by piece with smarter methods like `ReadString`, `ReadLine`, `ReadBytes`.

```go
reader := bufio.NewReader(conn) // wraps conn, adds internal buffer
msg, _ := reader.ReadString('\n') // reads chunk from conn, returns up to \n
```

**The real reason `bufio.Reader` exists — syscall cost**

Every time you read from a TCP connection, that's a syscall — asking the OS to copy bytes from the kernel's network buffer into your process. Syscalls are expensive relative to memory operations.

Without buffering, naive reading looks like this:

```go
// reading 1 byte at a time = 1 syscall per byte
b := make([]byte, 1)
for {
    conn.Read(b) // syscall every single time
}
```

`bufio.Reader` solves this by reading a large chunk at once (4096 bytes by default) in a single syscall, storing it internally, then serving your `ReadString` calls from that in-memory buffer without hitting the OS again until the buffer is exhausted.

```
One syscall → 4096 bytes land in bufio's internal buffer
ReadString('\n') → scans internal buffer, no syscall needed
ReadString('\n') → scans internal buffer, no syscall needed
ReadString('\n') → buffer exhausted, one more syscall → refill
```

**Analogy**

Think of a TCP connection like a water pipe and your program like someone filling cups:

- `bytes.Buffer` is a bucket — you wait for the pipe to fill the whole bucket, then use it
- Reading 1 byte at a time is putting your cup directly under the pipe for every sip
- `bufio.Reader` is a large tank that fills from the pipe in bulk, and you draw from the tank — fast and cheap per read

**The practical summary**

|                    | `bytes.Buffer`             | `bufio.Reader`                 |
| ------------------ | -------------------------- | ------------------------------ |
| Purpose            | In-memory storage          | Buffered reading from a source |
| Wraps              | Nothing                    | Any `io.Reader`                |
| Syscall efficiency | N/A                        | Yes, reads in chunks           |
| Good for           | Building/accumulating data | Reading streams efficiently    |
| Message boundaries | No                         | Yes (`ReadString`, `ReadLine`) |

When you later read `net/http` internals, you'll see Go wraps every TCP connection in a `bufio.Reader` immediately on accept — before any HTTP parsing happens. Now you know exactly why: it's both for efficiency and to get the delimiter-aware reading methods that HTTP parsing needs.
