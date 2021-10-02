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
        "first": "{{information.first_name}}",
        "last": "{{information.last_name}}",
"filter_value" : "{{filters.#.0}}"
    }
}`, oldJsonBody: `{
  "information": {
    "first_name": "magic",
    "last_name": "json"
  },
  "filters": [
    {
      "key": "type",
      "operator": 1,
      "values": [
        "a",
        "b"
      ]
    }
  ]
}`,
			expected: `{"name":{"filter_value":{"key":"type","operator":1,"values":["a","b"]},"first":"magic","last":"json"}}`},

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
		`Json array transform`: {newJsonBody: `{"name": {
            "first": "{{information.first_name}}",
            "last": "{{information.last_name}}"
        },
        "filter_value" : ["<<filters>>",
            {
                "type": "{{filters.#.<<filters>>.key.type}}",
                "operator": "{{filters.#.<<filters>>.operator}}"
            }]
    }`, oldJsonBody: `{
  "information": {
    "first_name": "magic",
    "last_name": "json"
  },
  "filters": [
    {
      "key": {
        "type" : "binary"
        },
      "operator": 1,
      "values": [
        "a",
        "b"
      ]
    },
    {
      "key": {
        "type" : "secondary"
        },
      "operator": 2,
      "values": [
        1232,
        4666
      ]
    }
  ]
}`,
			expected: `{"filter_value":[{"operator":1,"type":"binary"},{"operator":2,"type":"secondary"}],"name":{"first":"magic","last":"json"}}`},

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
