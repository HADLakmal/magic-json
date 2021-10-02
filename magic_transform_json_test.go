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

		`Json not existed transform`: {newJsonBody: `{
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
		`Json array transform`: {newJsonBody: `{
"records":["<<filters>>" ,
{"name": {
        "first": "{{filters.#.<<filters>>.operator}}",
        "last": "{{last_name}}"
    }}],
"filter_value" : "{{filters.#.0.key}}"
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
        },
{
            "key": "iterator",
            "operator": "EQUAL",
            "values": [
                1
            ]
        }
]

}`,
			expected: `{"filter_value":"type","records":[{"name":{"first":"MATCH","last":"lak"}},{"name":{"first":"EQUAL","last":"lak"}}]}`},

		`Json array transform extend`: {newJsonBody: `{
"records":["<<payload.filters>>" ,
{"name": {
        "first": "{{payload.filters.#.<<payload.filters>>.values}}",
        "last": "{{last_name}}"
    }}],
"filter_value" : "{{payload.filters.#.0.values.#.0}}"
}`, oldJsonBody: `{
    "first_name" : "dam",
    "last_name" : "lak",
"payload":{
"filters":[
{
            "key": "type",
            "values": [
                {
				"name" : "inner_array"
				},
{
				"name" : "inner_array_2"
				}
            ]
        },
{
            "key": "iterator",
            "values": [
                {
				"name" : "inner_array"
				}
            ]
        }
]
},


}`,
			expected: `{"filter_value":{"name":"inner_array"},"records":[{"name":{"first":[{"name":"inner_array"},{"name":"inner_array_2"}],"last":"lak"}},{"name":{"first":[{"name":"inner_array"}],"last":"lak"}}]}`},
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
