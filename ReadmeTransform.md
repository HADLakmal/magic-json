# Transform JSON

Magic Json library has ability to transform JSON into new JSON format. New JSON should have configured properly at value
field.

* Transform flat JSON
* Transform JSON with array
* Static values not be replaced

# Getting Start

There are several notations used for identify between new JSON and the old JSON. Those notations will be mapped the
transformed values.

```json
{
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
}
```

Above JSON have two fields and one array. JSON have multiple notations to map values into new JSON while using ```.```
in between flow.

| Notation      |   Value      | Explain   |
| -------------------- | ------------- |-----------------|
| ```"{{information.first_name}}"```  | ```magic```  | Exact value as field |
| ```"{{filters}}"```  | ```[{"key": "type","operator": 1,"values": ["a","b"]}]```  | Exact value as array |
| ```"{{filters.#.0.key}}"```  | ```type```  | Exact value as field from array |
| ```"{{filters.#.0}}"```  | ```{"key": "type","operator": 1,"values": ["a","b"]}```  | Exact value as object from first element of array |
| ```"{{filters.#.<<filters>>.key}}"```  | ```[{"key": "type","operator": 1,"values": ["a","b"]}]```  | Map JSON array to new JSON array (explain in example) |

### Example Usage

#### JSON Field Map

Map old JSON fields to new JSON as a flat JSON

```go
package main

import (
	"github.com/HADLakmal/magic-json"
)

const oldJson = `{
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
}`

func main() {
	newJson := `"{name": {
            "first": "{{information.first_name}}",
            "last": "{{information.last_name}}"
        },
        "filter_value" : "{{filters.#.0}}"
    }`
	m, err := mjson.NewTransferJSON(oldJson, newJson)
	if err != nil {
		panic(err)
	}

	json, _ := m.Release()
	println(json)
	// {"name":{"filter_value":{"key":"type","operator":1,"values":["a","b"]},"first":"magic","last":"json"}}
}
```

#### JSON Array Map

Map old JSON fields to new JSON as array JSON

```go
package main

import (
	"github.com/HADLakmal/magic-json"
)

const oldJson = `{
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
}`

func main() {
	newJson := `{"name": {
            "first": "{{information.first_name}}",
            "last": "{{information.last_name}}"
        },
        "filter_value" : ["<<filters>>",
            {
                "type": "{{filters.#.<<filters>>.key.type}}",
                "operator": "{{filters.#.<<filters>>.operator}}"
            }]
    }`
	m, err := mjson.NewTransferJSON(oldJson, newJson)
	if err != nil {
		panic(err)
	}

	json, _ := m.Release()
	println(json)
	// {"filter_value":[{"operator":1,"type":"binary"},{"operator":2,"type":"secondary"}],"name":{"first":"magic","last":"json"}}
}
```
