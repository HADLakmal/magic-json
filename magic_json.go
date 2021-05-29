package mjson

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type MJson struct {
	list
}

func (mj *MJson) extractor(data interface{}, n *node) {
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

func (mj *MJson) compound(n *node) interface{} {
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

func (mj *MJson) LoadJson(j string) {
	p := gjson.Parse(j).Value()
	n := mj.newHeader(p)
	mj.extractor(p, n)

	m := mj.compound(n)

	jsonBytes, errByte := json.Marshal(m)
	if errByte != nil {
		return
	}
	println(string(jsonBytes))
}
