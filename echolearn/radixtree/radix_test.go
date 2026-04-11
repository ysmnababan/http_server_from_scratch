package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: test add path "////"

func TestSplit(t *testing.T) {
	path := "/"
	slices := strings.Split(path, "/")
	for i, s := range slices {
		fmt.Printf("%d=> %s\n", i, s)
	}
}

func TestAddAllMatch_NewRoot(t *testing.T) {
	r := NewRadixTree()

	r.Add("a/b/c", 1)
	assert.Equal(t, 1, r.handler)
	r.Add("d/e/f", 2)
	assert.Equal(t, defaultHandler, r.handler)
	childA, ok := r.router[segment{key: "a"}]
	assert.True(t, ok)
	assert.Equal(t, 1, childA.handler)
	childD, ok := r.router[segment{key: "d"}]
	assert.True(t, ok)
	assert.Equal(t, 2, childD.handler)
	r.Add("x/y/z", 3)
	childX, ok := r.getChild("x")
	assert.True(t, ok)
	assert.Equal(t, 3, childX.handler)

	r.Add("a/b/t", 4)
	childA, ok = r.getChild("a")
	assert.True(t, ok)
	assert.Len(t, childA.pattern, 2)
	child, ok := childA.getChild("c")
	assert.True(t, ok)
	assert.Equal(t, 1, child.handler)
	child, ok = childA.getChild("t")
	assert.True(t, ok)
	assert.Equal(t, 4, child.handler)

	assert.Equal(t, 1, r.GetHandler("a/b/c"))
	assert.Equal(t, 2, r.GetHandler("d/e/f"))
	assert.Equal(t, 3, r.GetHandler("x/y/z"))
	assert.Equal(t, 4, r.GetHandler("a/b/t"))
}

func TestAddAllMatch_UnmatchIntheMiddle(t *testing.T) {
	r := NewRadixTree()

	r.Add("a/b/c", 1)
	assert.Equal(t, 1, r.handler)
	r.Add("a/b/d", 2)
	assert.Equal(t, defaultHandler, r.handler)
	childC, ok := r.router[segment{key: "c"}]
	assert.True(t, ok)
	assert.Equal(t, 1, childC.handler)
	childD, ok := r.router[segment{key: "d"}]
	assert.True(t, ok)
	assert.Equal(t, 2, childD.handler)
}

func TestAddAllMatch_Overwrite(t *testing.T) {
	r := NewRadixTree()

	r.Add("a/b/c", 1)
	assert.Equal(t, 1, r.handler)
	r.Add("a/b/c", 2)
	assert.Equal(t, 2, r.handler)
}

func TestAddAllMatch_Case1(t *testing.T) {
	r := NewRadixTree()

	r.Add("a/b/c", 1)
	assert.Equal(t, 1, r.handler)
	r.Add("a/b/c/d/e", 2)
	assert.Equal(t, 1, r.handler)
	child := r.router[segment{key: "d"}]
	assert.Equal(t, 2, child.handler)
}

func TestAddAllMatch_Case2(t *testing.T) {
	r := NewRadixTree()

	r.Add("a/b/c/d/e", 1)
	assert.Equal(t, 1, r.handler)
	r.Add("a/b/c", 2)
	assert.Equal(t, 2, r.handler)
	child := r.router[segment{key: "d"}]
	assert.Equal(t, 1, child.handler)
}
