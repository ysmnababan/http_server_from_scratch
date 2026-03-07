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
	// representationHeader map[string]string
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

// func (r *request) updateState(in string) {
// 	if r.parsingStatus == headerstatus && in == "" {
// 		r.parsingStatus = bodyStatus
// 		return
// 	}
// 	if lenStr, ok := r.header["Content-Length"]; ok {
// 		lenInt, _ := strconv.Atoi(lenStr)
// 		if r.parsingStatus == headerstatus && r.body.ContentLenght() == lenInt {
// 			r.parsingStatus = bodyStatus
// 		}
// 	}
// }

func (r *request) parse(in string) {
	// r.updateState(in)
	// if r.parsingStatus == undefined {
	// 	log.Print("error parsing")
	// 	return
	// }

	switch r.parsingStatus {
	case startlineStatus:
		sl := strings.Split(in, " ")
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
		if len(in) == 0 {
			r.parsingStatus = bodyStatus
			return
		}
		pair := strings.Split(in, ":")
		r.header[pair[0]] = pair[1]
	case bodyStatus:
		if lenStr, ok := r.header["Content-Length"]; ok {
			contentLength, _ := strconv.Atoi(lenStr)
			if r.body.ContentLenght() >= contentLength {
				r.parsingStatus = completeStatus
				return
			}
			r.body.Append(in)
		} else if val, ok := r.header["Transfer-Encoding"]; ok && val == "chunked" {
			log.Print("chunked data")
			r.body.Append(in)
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
