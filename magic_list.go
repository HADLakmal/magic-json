package mjson

// list of nodes which keep the given json
// as node objects
type list struct {
	head *node
}

// newHeader create new header
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

	n := newNode(`root`, t, nil, []string{})
	if lst.head == nil {
		lst.head = n
	}
	return n
}

// traversal through each node
func (lst *list) traversal(n *node, fn func(node *node)) {
	fn(n)
	if n.next == nil {
		return
	}

	for _, nd := range n.next {
		lst.traversal(nd, fn)
	}
}
