package mjson

import (
	"fmt"
	"regexp"
	"strings"
)

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

func (n *node) replaceCharacter(old, new string, count int) {
	var reg = regexp.MustCompile(fmt.Sprintf(`*%s*`, old))
	if reg.MatchString(n.name) {
		n.name = strings.Replace(n.name, old, new, count)
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
func (lst *list) traversal(n *node, fn func(node *node)) {
	fn(n)
	if n.next == nil {
		return
	}

	for _, nd := range n.next {
		lst.traversal(nd, fn)
	}
}
