package mjson

import "fmt"

type nodeType string

const (
	object nodeType = `object`
	array  nodeType = `array`
	field  nodeType = `field`
)

type node struct {
	next  []*node
	name  string
	nt    interface{}
	value interface{}
}

func newNode(name string, nt interface{}, value interface{}) *node {
	return &node{
		name:  name,
		nt:    nt,
		value: value,
	}
}

type list struct {
	head *node
}

func (lst *list) newHeader(data interface{}) *node {
	var t nodeType
	_, tMap := data.(map[string]interface{})
	_, tArray := data.([]interface{})
	switch {
	case tArray:
		t = array
	case tMap:
		t = object
	}

	n := newNode(`root`, t, nil)
	if lst.head == nil {
		lst.head = n
	}
	return n
}

func (lst *list) display() {
	ls := lst.head
	for ls != nil {
		fmt.Printf("%s (%s) = %v ->", ls.name, ls.nt, ls.value)
		ls = ls.next[0]
	}
	fmt.Println()
}
