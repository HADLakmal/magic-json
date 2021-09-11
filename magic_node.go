package mjson

import "fmt"

type nodeType string

const (
	object nodeType = `object`
	array  nodeType = `array`
	field  nodeType = `field`
)

// node compact the json key value relationship
type node struct {
	next  []*node
	name  string
	nt    interface{}
	value interface{}
	path  []string
}

// newNode create new node
func newNode(name string, nt interface{}, value interface{}, path []string) *node {
	pathName := name
	switch nt.(type) {
	case nodeType:
		switch nt {
		case array:
			path = append(path, name)
			pathName = "#"
		}
	case int:
		pathName = fmt.Sprintf(`%d`, nt)
	}

	return &node{
		name:  name,
		nt:    nt,
		value: value,
		path:  append(path, pathName),
	}
}
