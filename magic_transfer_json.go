package mjson

import (
	"strings"
)

// mJson holds the json object as a list of nodes
type mTransferJson struct {
	oldJson *mJson
	newJson *mJson
}

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

func (mj *mTransferJson) jsonNotationValueExtract() {
	mp := make(map[string]*interface{})
	mj.newJson.traversal(mj.newJson.head, func(node *node) {
		val, ok := node.value.(string)
		if ok {
			if strings.Contains(val, "{{") {
				node.value = "null"
				mp[strings.Replace(strings.Replace(val, "{{", "", 1), "}}", "", 1)] = &node.value
			}
		}
	})

	mj.oldJson.traversal(mj.oldJson.head, func(node *node) {
		if oldVal, ok := mp[strings.Join(node.path[1:], ".")]; ok {
			*oldVal = node.value
		}
	})
}

func (mj *mTransferJson) Release() (json string, err error) {
	mj.jsonNotationValueExtract()

	return mj.newJson.Release()
}
