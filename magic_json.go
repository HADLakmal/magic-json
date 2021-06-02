package mjson

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
)

type mJson struct {
	list
}

func NewMagicJson() JSONLoader {
	return new(mJson)
}

func (mj *mJson) extractor(data interface{}, n *node) {
	m, tMap := data.(map[string]interface{})
	d, tArray := data.([]interface{})
	switch {
	case tArray:
		for i, v := range d {
			tn := newNode(n.name, i, nil)
			n.next = append(n.next, tn)
			mj.extractor(v, tn)
		}
	case tMap:
		for v, k := range m {
			_, typeOfMap := k.(map[string]interface{})
			_, typeOfArray := k.([]interface{})

			switch {
			case typeOfMap:
				tm := newNode(v, object, nil)
				n.next = append(n.next, tm)
				mj.extractor(k, tm)
			case typeOfArray:
				tn := newNode(v, array, nil)
				n.next = append(n.next, tn)
				mj.extractor(k, tn)
			default:
				fn := newNode(v, field, k)
				n.next = append(n.next, fn)
			}
		}
	default:
		n.value = data
	}
}

func (mj *mJson) compound(n *node) interface{} {
	switch n.nt.(type) {
	case nodeType:
		switch n.nt {
		case field:
			return n.value
		case object:
			m := make(map[string]interface{})
			for _, nd := range n.next {
				m[nd.name] = mj.compound(nd)
			}
			return m
		case array:
			m := make([]interface{}, len(n.next))
			for i, nd := range n.next {
				m[i] = mj.compound(nd)
			}
			return m
		}
	case int:
		if n.value != nil {
			return n.value
		}
		m := make(map[string]interface{})
		for _, nd := range n.next {
			m[nd.name] = mj.compound(nd)
		}
		return m
	}

	return nil
}

func (mj *mJson) Load(j string) (JSONConverter, error) {
	// parse json string from gjson lib
	p := gjson.Parse(j).Value()
	if p == nil {
		return nil, fmt.Errorf(`json can't convert to json map'`)
	}
	// json object extract into node list
	mj.extractor(p, mj.newHeader(p))

	return mj, nil
}

func (mj *mJson) Release() (string, error) {
	m := mj.compound(mj.head)

	jsonBytes, errByte := json.Marshal(m)
	if errByte != nil {
		return "", errByte
	}

	return string(jsonBytes), nil
}

func (mj *mJson) ReplaceKey(oldCharacters, newCharacters string) {
	// old character is empty
	if oldCharacters == "" {
		return
	}

	mj.traversal(mj.head, func(node *node) {
		var reg = regexp.MustCompile(fmt.Sprintf(`%s*`, oldCharacters))
		if reg.MatchString(node.name) {
			node.name = strings.Replace(node.name, oldCharacters, newCharacters, 1)
		}
	})
}

func (mj *mJson) ReplaceValue(oldCharacters, newCharacters string) {
	// old character is empty
	if oldCharacters == "" {
		return
	}

	mj.traversal(mj.head, func(node *node) {
		if node.value == nil {
			return
		}
		if v, ok := node.value.(string); ok {
			var reg = regexp.MustCompile(fmt.Sprintf(`%s*`, oldCharacters))
			if reg.MatchString(v) {
				node.value = strings.Replace(v, oldCharacters, newCharacters, 1)
			}
		}
	})
}

//func (mj *mJson) LoadJson(j string) {
//	p := gjson.Parse(j).Value()
//	n := mj.newHeader(p)
//	mj.extractor(p, n)
//
//	mj.traversal(mj.head, func(node *node) {
//		if node.value == nil {
//			return
//		}
//		switch v := node.value.(type) {
//		case string:
//			var reg = regexp.MustCompile(`_*`)
//			if reg.MatchString(v) {
//				node.value = strings.Replace(v, `_`, ``, 1)
//			}
//		}
//	})
//
//	m := mj.compound(n)
//
//	jsonBytes, errByte := json.Marshal(m)
//	if errByte != nil {
//		return
//	}
//	println(string(jsonBytes))
//}
