package mjson

import (
	"fmt"
	"testing"
)

func TestMJson_NewTransferJSON(t *testing.T) {
	tests := map[string]struct {
		oldJsonBody string
		newJsonBody string
		wantError   bool
		expected    string
	}{
		`Json value and array transfer`: {newJsonBody: `{
    "name": {
        "first": "{{first_name}}",
        "last": "{{last_name}}",
"filter_value" : "{{filters.#}}"
    }
}`, oldJsonBody: `{
    "first_name" : "dam",
    "last_name" : "lak",
"filters":[
{
            "key": "type",
            "operator": "MATCH",
            "values": [
                "ASSET"
            ]
        }
]

}`,
			expected: `{"name":{"filter_value":[{"key":"type","operator":"MATCH","values":["ASSET"]}],"first":"dam","last":"lak"}}`},

		`Json value and array single value transfer`: {newJsonBody: `{
    "name": {
        "first": "{{first_name}}",
        "last": "{{last_name}}"
    },
"filter_value" : "{{filters.#.0}}"
}`, oldJsonBody: `{
    "first_name" : "dam",
    "last_name" : "lak",
"filters":[
{
            "key": "type",
            "operator": "MATCH",
            "values": [
                1
            ]
        }
]

}`,
			expected: `{"filter_value":{"key":"type","operator":"MATCH","values":[1]},"name":{"first":"dam","last":"lak"}}`},

		`Json not existed transfer`: {newJsonBody: `{
    "name": {
        "first": "{{first_name}}",
        "last": "{{last_name}}"
    },
"filter_value" : "{{filters.#.0.name}}"
}`, oldJsonBody: `{
    "first_name" : "dam",
    "last_name" : "lak",
"filters":[
{
            "key": "type",
            "operator": "MATCH",
            "values": [
                1
            ]
        }
]

}`,
			expected: `{"filter_value":"null","name":{"first":"dam","last":"lak"}}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := NewTransferJSON(test.oldJsonBody, test.newJsonBody)
			if err != nil {
				panic(err)
			}
			s, err := p.Release()
			if err != nil && !test.wantError {
				panic(err)
			}
			if s != test.expected {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, s, test.expected))
			}
		})
	}
}
