package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var (
	startlineStatus string = "startline"
	headerstatus    string = "header"
	bodyStatus      string = "body"
	undefined       string = "undefined"
	completeStatus  string = "complete"
)

type header struct {
	requestHeader map[string]string
}

type body struct {
	raw string
}

func (b *body) ContentLenght() int {
	return len(b.raw)
}

func (b *body) Append(in string) {
	b.raw = b.raw + in
}

type request struct {
	method   string
	target   string
	protocol string
	header   map[string]string
	body     body

	parsingStatus string
	currBodyLen   int
}

func NewRequest() *request {
	// requestHeader, representationHeader := make(map[string]string), make(map[string]string)
	req := &request{
		header:        make(map[string]string),
		parsingStatus: startlineStatus,
	}
	return req
}

func (r *request) parse(in string) {
	switch r.parsingStatus {
	case startlineStatus:
		in = strings.TrimSpace(in)
		sl := strings.Split(in, " ")
		// fmt.Println(in, sl)
		if len(sl) != 3 {
			r.parsingStatus = undefined
			return
		}
		r.method = sl[0]
		r.target = sl[1]
		r.protocol = sl[2]
		r.parsingStatus = headerstatus
		return
	case headerstatus:
		in = strings.TrimSpace(in)
		if len(in) == 0 {
			r.parsingStatus = bodyStatus
			return
		}
		pair := strings.Split(in, ":")
		r.header[pair[0]] = pair[1]
	case bodyStatus:
		r.body.Append(in)
		if lenStr, ok := r.header["Content-Length"]; ok {
			contentLength, err := strconv.Atoi(strings.TrimSpace(lenStr))
			if err != nil {
				panic(err)
			}
			if r.body.ContentLenght() >= contentLength {
				r.parsingStatus = completeStatus
				return
			}
		} else if val, ok := r.header["Transfer-Encoding"]; ok && val == "chunked" {
			log.Print("chunked data")
		} else {
			r.parsingStatus = completeStatus
		}
	}
}

func (r *request) debug() {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>")
	fmt.Printf("<%s> <%s> <%s>\n", r.method, r.target, r.protocol)
	for k, val := range r.header {
		fmt.Printf("%s: %s\n", k, val)
	}
	fmt.Println()
	fmt.Println(r.body.raw)
	fmt.Println("<<<<<<<<<<<<<<<<<<<<")
	fmt.Println()
}
