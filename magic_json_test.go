package mjson

import "testing"

func TestMJson_ReplaceValue(t *testing.T) {
	tests := map[string]struct {
		jsonBody string
	}{
		`object Json `: {jsonBody: `{
    "filters": [
        {
            "key": "type",
            "operator": "MATCH",
            "values": [
                "ASSET"
            ]
        },
		{
            "key": "type",
            "operator": "MATCH_VALUE",
            "values": [
                "ASSET"
            ]
        }
    ],
    "order_by": [{
	"key": "type",
	}
	],
    "paging": {
        "offset": 0,
        "size": 20
    }
}`},
		`array Json `: {jsonBody: `[
    {
        "headers": [
            {
                "value": "0a40e6a9-1216-426a-977a-7d13a36dc64e",
                "key": "account_id"
            },
            {
                "value": "0a40e6a9-1216-426a-977a-7d13a36dc64e",
                "key": "created_by"
            },
            {
                "value": "e463fe2f-8dcc-41a5-999c-6b886c9101fa",
                "key": "trace_id"
            }
        ]
    }
]`},
		`sample Json `: {jsonBody: `{
    "filters": [
        {
            "key": "type",
            "operator": "MATCH",
            "values": [
                "ASSET"
            ]
        }
    ]
}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := NewMagicJson()
			p, err := m.Load(test.jsonBody)
			if err != nil {
				panic(err)
			}
			p.ReplaceValue("_", "")
			str, err := m.Release()
			if err != nil {
				panic(err)
			}
			t.Log(str)
		})
	}
}

func TestMJson_ReplaceKey(t *testing.T) {
	tests := map[string]struct {
		jsonBody string
	}{
		`object Json `: {jsonBody: `{
    "filters": [
        {
            "key": "type",
            "operator": "MATCH",
            "values": [
                "ASSET"
            ]
        },
		{
            "key": "type",
            "operator": "MATCH",
            "values": [
                "ASSET"
            ]
        }
    ],
    "order_by": [{
	"key": "type",
	}
	],
    "paging": {
        "offset": 0,
        "size": 20
    }
}`},
		`array Json `: {jsonBody: `[
    {
        "headers": [
            {
                "value": "0a40e6a9-1216-426a-977a-7d13a36dc64e",
                "key": "account_id"
            },
            {
                "value": "0a40e6a9-1216-426a-977a-7d13a36dc64e",
                "key": "created_by"
            },
            {
                "value": "e463fe2f-8dcc-41a5-999c-6b886c9101fa",
                "key": "trace_id"
            }
        ]
    }
]`},
		`sample Json `: {jsonBody: `{
    "filters": [
        {
            "key": "type",
            "operator": "MATCH",
            "values": [
                "ASSET"
            ]
        }
    ]
}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := NewMagicJson()
			p, err := m.Load(test.jsonBody)
			if err != nil {
				panic(err)
			}
			p.ReplaceKey("_", "")
			str, err := m.Release()
			if err != nil {
				panic(err)
			}
			t.Log(str)
		})
	}
}
