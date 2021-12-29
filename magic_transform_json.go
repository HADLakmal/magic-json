package mjson

import (
	"fmt"
	"strings"
)

// mTransferJson holds the old and new json object as a list of nodes
type mTransferJson struct {
	oldJson *mJson
	newJson *mJson
}

// NewTransferJSON transfer old json in to new json
func NewTransferJSON(oldJson, newJson string) (JSONTransfer, error) {
	mOld, err := new(mJson).load(oldJson)
	if err != nil {
		return nil, err
	}

	mNew, err := new(mJson).load(newJson)
	if err != nil {
		return nil, err
	}

	return &mTransferJson{
		oldJson: mOld,
		newJson: mNew,
	}, nil
}

// jsonNotationValueExtract transform new json in to desired format
func (mj *mTransferJson) jsonNotationValueExtract() {
	//
	mj.newJson.traversal(mj.newJson.head, func(node *node) {
		if node.nt == array && len(node.next) > 1 {
			val, ok := node.next[0].value.(string)
			if ok && strings.Contains(val, "<<") {
				mj.duplicate(node, val)
			}
		}
	})

	valueIndexes := make([]string, 0)
	valuePointers := make([]*interface{}, 0)
	mj.newJson.traversal(mj.newJson.head, func(node *node) {
		if node.nt == field {
			val, ok := node.value.(string)
			if ok {
				if strings.Contains(val, "{{") {
					node.value = "null"
					valueIndexes = append(valueIndexes, strings.Replace(strings.Replace(val, "{{", "", 1), "}}", "", 1))
					valuePointers = append(valuePointers, &node.value)
				}
			}
		}
	})

	mj.oldJson.traversal(mj.oldJson.head, func(node *node) {
		valuePath := strings.Join(node.path[1:], ".")
		for r, i := range valueIndexes {
			if valuePath == i {
				*valuePointers[r] = node.value
			}
		}
	})
}

// Release the new json string
func (mj *mTransferJson) Release() (json string, err error) {
	mj.jsonNotationValueExtract()

	return mj.newJson.Release()
}

// Release the new json
func (mj *mTransferJson) ReleaseJson() interface{} {
	mj.jsonNotationValueExtract()

	return mj.newJson.ReleaseJson()
}

// duplicate new nodes
func (mj *mTransferJson) duplicate(n *node, key string) {
	if n.nt != array || len(n.next) < 2 {
		return
	}

	// array element count of transform json
	var count int
	k := strings.Replace(strings.Replace(key, "<<", "", 1), ">>", "", 1)
	mj.oldJson.traversal(mj.oldJson.head, func(node *node) {
		if k == strings.Join(node.path[1:], ".") && n.nt == array {
			count = len(node.next)
		}
	})

	nodeVal := *n.next[1]
	n.next = make([]*node, 0)
	for i := 0; i < count; i++ {
		pathLength := len(nodeVal.path) - 1
		// newly create node to inject
		tempNode := new(node)
		tempNode.nt = i
		tempNode.path = make([]string, 0)
		tempNode.next = make([]*node, 0)
		tempNode.next = append(tempNode.next, nodeVal.next...)
		for x := 0; x < pathLength; x++ {
			tempNode.path = append(tempNode.path, nodeVal.path[x])
		}
		tempNode.path = append(tempNode.path, fmt.Sprintf(`%d`, tempNode.nt))

		// update each and every node with respect to array name
		mj.newJson.traversal(tempNode, func(n *node) {
			if n.next == nil {
				return
			}

			tempNext := n.next
			n.next = make([]*node, 0)
			for _, nd := range tempNext {
				var nt node
				nt.nt = nd.nt
				nt.name = nd.name
				v, ok := nd.value.(string)
				if nd.nt == field && ok {
					if strings.Contains(v, "<<") {
						nt.value = strings.Replace(nd.value.(string), key, fmt.Sprintf(`%d`, i), 1)
					} else {
						nt.value = v
					}
				}
				nt.next = nd.next
				n.next = append(n.next, &nt)
			}

		})

		n.next = append(n.next, tempNode)
	}
}
