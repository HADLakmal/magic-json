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
	newPath := make([]string, 0)
	for _, p := range path {
		newPath = append(newPath, p)
	}

	pathName := name

	switch nt.(type) {
	case int:
		newPath = append(newPath,"#")
		pathName = fmt.Sprintf(`%d`, nt)
	}

	return &node{
		name:  name,
		nt:    nt,
		value: value,
		path:  append(newPath, pathName),
	}
}
