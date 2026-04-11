package main

import (
	"fmt"
	"strings"
)

// segment is something like this /segment1/segment2/segment3
// so there are 3 segments is the above path.

type segment struct {
	key       string
	isDynamic bool
}

func isDynamic(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] == ':'
}

func NewSegment(s string) segment {
	return segment{
		key:       s,
		isDynamic: isDynamic(s),
	}
}

func (s *segment) isSameSegment(t segment) bool {
	return s.key == t.key
}

func getSegment(path string) []segment {
	s := strings.Split(path, "/")
	if len(s) == 0 {
		return nil
	}
	out := make([]segment, 0, len(s))
	for _, item := range s {
		out = append(out, NewSegment(item))
	}
	return out
}

var defaultHandler int = 999

type RadixTree struct {
	handler int
	router  map[segment]*RadixTree
	pattern []segment
}

func NewRadixTree() *RadixTree {
	return &RadixTree{
		router: make(map[segment]*RadixTree),
	}
}

func (r *RadixTree) getChild(key string) (*RadixTree, bool) {
	node, ok := r.router[segment{key: key}]
	return node, ok
}

func (r *RadixTree) GetHandler(path string) int {
	fmt.Println(path)
	segments := getSegment(path)
	if len(r.pattern) == 0 {
		first := segments[0]
		node, ok := r.router[first]
		if !ok {
			return defaultHandler
		}
		fmt.Println("here")
		return node.GetHandler(path)
	}
	minLen := min(len(r.pattern), len(segments))
	i := 0
	for i < minLen {
		if !r.pattern[i].isSameSegment(segments[i]) {
			return defaultHandler
		}
		i++
	}
	// ab in abc
	// abc in abc => no other
	// abc in ab
	// abcde in abc + de
	if len(segments) < len(r.pattern) {
		return defaultHandler
	}
	if len(segments) == len(r.pattern) {
		return r.handler
	}
	node, ok := r.router[segments[i]]
	if !ok {
		return defaultHandler
	}
	return node.GetHandler(toPath(segments[i:]))
}

func toPath(s []segment) string {
	out := make([]string, 0, len(s))
	for _, i := range s {
		out = append(out, i.key)
	}
	return strings.Join(out, "/")
}

func (r *RadixTree) Add(path string, hf int) {
	segments := getSegment(path)

	if len(r.pattern) == 0 {
		// search in router
		first := segments[0]
		node, ok := r.router[first]
		if !ok {
			if len(r.router) == 0 {
				// empty pattern and router
				r.pattern = append(r.pattern, segments...)
				r.handler = hf
				return
			}
			// no node found but already has sibling,
			// so add as a child
			newNode := &RadixTree{
				pattern: segments,
				handler: hf,
				router:  make(map[segment]*RadixTree),
			}
			r.router[first] = newNode
			return
		}

		node.Add(path, hf)
		return
	}
	minLen := min(len(r.pattern), len(segments))
	allMatch := true
	i := 0
	for i < minLen {
		if !r.pattern[i].isSameSegment(segments[i]) {
			allMatch = false
			break
		}
		i++
	}
	if allMatch {
		if len(r.pattern) == len(segments) {
			// rewrite the handler
			r.handler = hf
			return
		}
		// a/b/c     => registered
		// a/b/c/d/e
		if len(r.pattern) < len(segments) {
			newNode := &RadixTree{
				handler: hf,
				pattern: segments[i:],
				router:  (make(map[segment]*RadixTree)),
			}
			r.router[segments[i]] = newNode
			return
		}
		// a/b/c/d/e => registered
		// a/b/c
		newNode := &RadixTree{
			handler: r.handler,
			pattern: r.pattern[i:],
			router:  (make(map[segment]*RadixTree)),
		}
		r.handler = hf
		r.router[r.pattern[i]] = newNode
		r.pattern = r.pattern[:i]
		return
	}

	if i == 0 {
		// difference found in root
		// a,b,c =>registered
		// d,e,f
		newNode := &RadixTree{
			pattern: segments,
			handler: hf,
			router:  make(map[segment]*RadixTree),
		}
		prevNode := &RadixTree{
			pattern: r.pattern,
			handler: r.handler,
			router:  r.router,
		}
		r.router[prevNode.pattern[0]] = prevNode
		r.router[newNode.pattern[0]] = newNode
		r.handler = defaultHandler
		r.pattern = r.pattern[:0]
		// r.pattern = append(r.pattern, &segment{}) // empty segment
		return
	}

	// add a,b,c => registered
	// add a,b,d
	newNode := &RadixTree{
		pattern: segments[i:],
		handler: hf,
		router:  make(map[segment]*RadixTree),
	}
	prevNode := &RadixTree{
		pattern: r.pattern[i:],
		handler: r.handler,
		router:  r.router,
	}
	r.router = make(map[segment]*RadixTree)
	r.router[prevNode.pattern[0]] = prevNode
	r.router[newNode.pattern[0]] = newNode
	r.handler = defaultHandler
	r.pattern = r.pattern[:i]
}
