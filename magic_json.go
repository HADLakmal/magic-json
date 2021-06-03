package mjson

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"regexp"
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

// extractor bind map of json information into list of nodes
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

// load parse the string json into a mJson
func (mj *mJson) load(j string) (MagicJSON, error) {
	// parse json string from gjson lib
	p := gjson.Parse(j).Value()
	if p == nil {
		return nil, fmt.Errorf(`json can't convert to json map'`)
	}
	// json object extract into node list
	mj.extractor(p, mj.newHeader(p))

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
	})

	return mj
}

// Release mJson objects
func (mj *mJson) Release() (string, error) {
	m := mj.compound(mj.head)

	jsonBytes, errByte := json.Marshal(m)
	if errByte != nil {
		return "", errByte
	}

	return string(jsonBytes), nil
}

func (mj *mJson) ReplaceKeyCharacter(oldCharacters, newCharacters string) JSONRelease {
	// old character is empty
	if oldCharacters == "" {
		return mj
	}

	mj.keyIdentifier(mj.beginNode, oldCharacters, newCharacters, 1)

	return mj
}

func (mj *mJson) ReplaceCharsInKey(oldCharacters, newCharacters string, count int) JSONRelease {
	// old character is empty
	if oldCharacters == "" {
		return mj
	}

	mj.keyIdentifier(mj.beginNode, oldCharacters, newCharacters, count)

	return mj
}

func (mj *mJson) keyIdentifier(n *node, oldCharacters, newCharacters string, count int) {
	mj.traversal(n, func(node *node) {
		var reg = regexp.MustCompile(fmt.Sprintf(`%s*`, oldCharacters))
		if reg.MatchString(node.name) {
			node.name = strings.Replace(node.name, oldCharacters, newCharacters, count)
		}
	})
}

func (mj *mJson) ReplaceCharInValue(oldCharacters, newCharacters string) JSONRelease {
	// old character is empty
	if oldCharacters == "" {
		return mj
	}

	mj.stringValueIdentifier(mj.beginNode, func(s string) interface{} {
		var reg = regexp.MustCompile(fmt.Sprintf(`%s*`, oldCharacters))
		if reg.MatchString(s) {
			return strings.Replace(s, oldCharacters, newCharacters, 1)
		}
		return s
	})

	return mj
}

func (mj *mJson) ReplaceValueCharacters(oldCharacters, newCharacters string, count int) JSONRelease {
	// old character is empty
	if oldCharacters == "" {
		return mj
	}

	mj.stringValueIdentifier(mj.beginNode, func(s string) interface{} {
		var reg = regexp.MustCompile(fmt.Sprintf(`%s*`, oldCharacters))
		if reg.MatchString(s) {
			return strings.Replace(s, oldCharacters, newCharacters, count)
		}
		return s
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
		if node.value == nil {
			return
		}
		if v, ok := node.value.(string); ok {
			node.value = stFn(v)
		}
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
