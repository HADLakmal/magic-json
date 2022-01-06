package mjson

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

// mJson holds the json object as a list of nodes
type mJson struct {
	list
	beginNode *node
}

// NewMagicJSON create the new mJson
func NewMagicJSON(json string) (MagicJSON, error) {
	return new(mJson).load(json)
}

// NewMagicJSONInterface create the new mJson
func NewMagicJSONInterface(jsonIntr interface{}) (MagicJSON, error) {
	return new(mJson).loadInterface(jsonIntr)
}

// extractor bind map of json information into list of nodes
func (mj *mJson) extractor(data interface{}, n *node, ancestorsName []string) {
	m, tMap := data.(map[string]interface{})
	d, tArray := data.([]interface{})
	switch {
	case tArray:
		for i, v := range d {
			tn := newNode(n.name, i, v, ancestorsName)
			n.next = append(n.next, tn)
			mj.extractor(v, tn, tn.path)
		}
	case tMap:
		for v, k := range m {
			_, typeOfMap := k.(map[string]interface{})
			_, typeOfArray := k.([]interface{})

			switch {
			case typeOfMap:
				tm := newNode(v, object, k, ancestorsName)
				n.next = append(n.next, tm)
				mj.extractor(k, tm, tm.path)
			case typeOfArray:
				tn := newNode(v, array, k, ancestorsName)
				n.next = append(n.next, tn)
				mj.extractor(k, tn, tn.path)
			default:
				fn := newNode(v, field, k, ancestorsName)
				n.next = append(n.next, fn)
			}
		}
	default:
		n.nt = field
		n.value = data
	}
}

// compound is called when lib need to release the json as a string
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
		m := make(map[string]interface{})
		for _, nd := range n.next {
			m[nd.name] = mj.compound(nd)
		}
		return m
	}

	return nil
}

// load parse the string json into a mJson
func (mj *mJson) load(j string) (*mJson, error) {
	// parse json string from gjson lib
	p := gjson.Parse(j).Value()
	if p == nil {
		return nil, fmt.Errorf(`json can't convert to json map'`)
	}
	// json object extract into node list
	n := mj.newHeader(p)
	mj.extractor(p, n, n.path)

	// attach beginning node as header
	mj.beginNode = mj.head

	return mj, nil
}

// loadInterface parse the interface object into a mJson
func (mj *mJson) loadInterface(p interface{}) (*mJson, error) {
	// json object extract into node list
	n := mj.newHeader(p)
	mj.extractor(p, n, n.path)

	// attach beginning node as header
	mj.beginNode = mj.head

	return mj, nil
}

// Key key finder
func (mj *mJson) Key(key string) JSONConverter {
	mj.traversal(mj.head, func(node *node) {
		if node.name == key {
			mj.beginNode = node
		}
		return
	})

	return mj
}

// Release mJson objects as string
func (mj *mJson) Release() (string, error) {
	m := mj.compound(mj.head)

	jsonBytes, errByte := json.Marshal(m)
	if errByte != nil {
		return "", errByte
	}

	return string(jsonBytes), nil
}

// ReleaseJson mJson objects as json
func (mj *mJson) ReleaseJson() interface{} {
	return mj.compound(mj.head)
}

func (mj *mJson) ReplaceCharInKey(oldCharacters, newCharacters string) JSONRelease {
	// old character is empty
	if oldCharacters == "" {
		return mj
	}

	mj.keyIdentifier(mj.beginNode, func(s string) string {
		return strings.Replace(s, oldCharacters, newCharacters, 1)
	})

	return mj
}

func (mj *mJson) ReplaceCharsInKey(oldCharacters, newCharacters string, count int) JSONRelease {
	// old character is empty
	if oldCharacters == "" {
		return mj
	}

	mj.keyIdentifier(mj.beginNode, func(s string) string {
		return strings.Replace(s, oldCharacters, newCharacters, count)
	})

	return mj
}

func (mj *mJson) KeyStringConverter(fn func(value string) string) JSONRelease {
	mj.keyIdentifier(mj.beginNode, func(s string) string {
		return fn(s)
	})

	return mj
}

func (mj *mJson) keyIdentifier(n *node, stFn func(s string) string) {
	mj.traversal(n, func(node *node) {
		node.name = stFn(node.name)
	})
}

func (mj *mJson) ReplaceCharInValue(oldCharacters, newCharacters string) JSONRelease {
	// old character is empty
	if oldCharacters == "" {
		return mj
	}

	mj.stringValueIdentifier(mj.beginNode, func(s string) interface{} {
		return strings.Replace(s, oldCharacters, newCharacters, 1)
	})

	return mj
}

func (mj *mJson) ReplaceCharsInValue(oldCharacters, newCharacters string, count int) JSONRelease {
	// old character is empty
	if oldCharacters == "" {
		return mj
	}

	mj.stringValueIdentifier(mj.beginNode, func(s string) interface{} {
		return strings.Replace(s, oldCharacters, newCharacters, count)
	})

	return mj
}

func (mj *mJson) ValueStringToInt() JSONRelease {
	mj.stringValueIdentifier(mj.beginNode, func(s string) interface{} {
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			return i
		}

		return s
	})

	return mj
}

func (mj *mJson) ValueStringToFloat() JSONRelease {
	mj.stringValueIdentifier(mj.beginNode, func(s string) interface{} {
		i, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return i
		}

		return s
	})

	return mj
}

func (mj *mJson) ValueStringConverter(fn func(value string) interface{}) JSONRelease {
	mj.stringValueIdentifier(mj.beginNode, func(s string) interface{} {
		return fn(s)
	})

	return mj
}

func (mj *mJson) stringValueIdentifier(n *node, stFn func(s string) interface{}) {
	mj.traversal(n, func(node *node) {
		if v, ok := node.value.(string); ok {
			node.value = stFn(v)
		}
	})
}

func (mj *mJson) ValueChecker(fn func(value interface{}) interface{}) JSONRelease {
	mj.valueChecker(mj.beginNode, func(s interface{}) interface{} {
		return fn(s)
	})

	return mj
}

func (mj *mJson) valueChecker(n *node, stFn func(s interface{}) interface{}) {
	mj.traversal(n, func(node *node) {
		node.value = stFn(node.value)
	})
}

func (mj *mJson) IntToString() JSONRelease {
	mj.intValueIdentifier(mj.beginNode, func(val int64) interface{} {
		return fmt.Sprintf(`%v`, val)
	})

	return mj
}

func (mj *mJson) IntConverter(fn func(value int64) interface{}) JSONRelease {
	mj.intValueIdentifier(mj.beginNode, func(val int64) interface{} {
		return fn(val)
	})

	return mj
}

func (mj *mJson) intValueIdentifier(n *node, intFn func(val int64) interface{}) {
	mj.traversal(n, func(node *node) {
		if node.value == nil {
			return
		}
		s := fmt.Sprintf(`%v`, node.value)
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			node.value = intFn(i)
		}
	})
}

func (mj *mJson) FloatToString() JSONRelease {
	mj.floatValueIdentifier(mj.beginNode, func(val float64) interface{} {
		return fmt.Sprintf(`%v`, val)
	})

	return mj
}

func (mj *mJson) FloatToInt() JSONRelease {
	mj.floatValueIdentifier(mj.beginNode, func(val float64) interface{} {
		return int(val)
	})

	return mj
}

func (mj *mJson) FloatConverter(fn func(value float64) interface{}) JSONRelease {
	mj.floatValueIdentifier(mj.beginNode, func(val float64) interface{} {
		return fn(val)
	})

	return mj
}

func (mj *mJson) floatValueIdentifier(n *node, floatFn func(val float64) interface{}) {
	mj.traversal(n, func(node *node) {
		if node.value == nil {
			return
		}
		s := fmt.Sprintf(`%v`, node.value)
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			node.value = floatFn(f)
		}
	})
}
