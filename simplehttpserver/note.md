# HTTP Behavior Notes

## Request / Response Structure

```
<start-line>\r\n
<header>: <value>\r\n
<header>: <value>\r\n
\r\n
<body>
```

- Headers and the blank line separator use `\r\n`
- The body is just raw bytes — HTTP imposes no formatting on it

---

## Body

- HTTP does **not** care how the body is formatted internally
- Rules come from `Content-Type`, not HTTP itself:
  - `application/json` → must be valid JSON
  - `text/html` → HTML
  - `application/octet-stream` → raw binary
- No `\r\n` required in the body

---

## How the Receiver Knows Where the Body Ends

| Method                       | How it works                               | When to use                                           |
| ---------------------------- | ------------------------------------------ | ----------------------------------------------------- |
| `Content-Length`             | Receiver reads exactly N bytes             | You know the size upfront (e.g. after `json.Marshal`) |
| `Transfer-Encoding: chunked` | Body sent in chunks, ends with `0\r\n\r\n` | You don't know the size upfront                       |
| Connection close             | Read until connection drops                | Fragile, mostly obsolete                              |

---

## Chunked Transfer Encoding

Used when you need to **start sending before you know the total size**.

### Format

```
<size in hex>\r\n
<data>\r\n
<size in hex>\r\n
<data>\r\n
0\r\n
\r\n
```

### Example

```
HTTP/1.1 200 OK
Transfer-Encoding: chunked

7\r\n
{"msg"\r\n
10\r\n
:"hello world"}\r\n
0\r\n
\r\n
```

### Real use cases

- LLM APIs streaming tokens (Claude, ChatGPT)
- Large file downloads (avoid buffering entire file)
- Streaming DB query results
- Live/real-time data (logs, stock prices)

---

## Go String Literals (reminder)

```go
body := `{"message": "hello"}`       // raw string literal (backticks)
body  = "{\"message\": \"hello\"}"   // interpreted string literal
```

Both produce the same string — HTTP only sees the final content.
