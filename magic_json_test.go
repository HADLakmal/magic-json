package mjson_test

import (
	"fmt"
	mjson "github.com/HADLakmal/magic-json"
	"reflect"
	"testing"
)

func TestMJson_ReplaceValue(t *testing.T) {
	tests := map[string]struct {
		jsonBody  string
		wantError bool
		expected  string
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
            "operator": "MATCH__VALUE",
            "values": [
                "ASSET_"
            ]
        }
    ],
    "order_by": [{
	"key": "type_",
	}
	],
    "paging": {
        "offset": 0,
        "size": 20
    }
}`, expected: `{"filters":[{"key":"type","operator":"MATCH","values":["ASSET"]},{"key":"type","operator":"MATCH_VALUE","values":["ASSET"]}],"order_by":[{"key":"type"}],"paging":{"offset":0,"size":20}}`},
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
]`, expected: `[{"headers":[{"key":"accountid","value":"0a40e6a9-1216-426a-977a-7d13a36dc64e"},{"key":"createdby","value":"0a40e6a9-1216-426a-977a-7d13a36dc64e"},{"key":"traceid","value":"e463fe2f-8dcc-41a5-999c-6b886c9101fa"}]}]`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSON(test.jsonBody)
			if err != nil {
				panic(err)
			}
			p.ReplaceCharInValue("_", "")
			str, err := p.Release()
			if err != nil && !test.wantError {
				panic(err)
			}
			if str != test.expected {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}

func TestMJson_ObjectKeyReplaceValue(t *testing.T) {
	tests := map[string]struct {
		jsonBody  string
		wantError bool
		expected  string
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
            "operator": "MATCH__VALUE",
            "values": [
                "ASSET_"
            ]
        }
    ],
    "order_by": [{
	"key": "type_",
	}
	],
    "paging": {
        "offset": 0,
        "size": 20
    }
}`, expected: `{"filters":[{"key":"type","operator":"MATCH","values":["ASSET"]},{"key":"type","operator":"MATCH_VALUE","values":["ASSET"]}],"order_by":[{"key":"type_"}],"paging":{"offset":0,"size":20}}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSON(test.jsonBody)
			if err != nil {
				panic(err)
			}
			p.Key("filters").ReplaceCharInValue("_", "")
			str, err := p.Release()
			if err != nil && !test.wantError {
				panic(err)
			}
			if str != test.expected {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}

func TestMJson_ValueStringConverter(t *testing.T) {
	tests := map[string]struct {
		jsonBody  string
		wantError bool
		expected  string
	}{
		`object Json `: {jsonBody: `{
    "filters": [
        {
            "key": "type",
            "operator": "MATCH",
            "values": [
                "ASSET"
            ]
        }
    ]
}`, expected: `{"filters":[{"key":"replaced","operator":"replaced","values":["replaced"]}]}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSON(test.jsonBody)
			if err != nil {
				panic(err)
			}
			p.ValueStringConverter(func(value string) interface{} {
				return "replaced"
			})
			str, err := p.Release()
			if err != nil && !test.wantError {
				panic(err)
			}
			if str != test.expected {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}

func TestMJson_ReplaceKey(t *testing.T) {
	tests := map[string]struct {
		jsonBody  string
		wantError bool
		expected  string
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
}`, expected: `{"filters":[{"key":"type","operator":"MATCH","values":["ASSET"]},{"key":"type","operator":"MATCH","values":["ASSET"]}],"orderby":[{"key":"type"}],"paging":{"offset":0,"size":20}}`},
		`sample Json `: {jsonBody: `{
    "filters_": [
        {
            "key": "type",
            "operator": "MATCH",
            "values": [
                "ASSET"
            ]
        }
    ]
}`, expected: `{"filters":[{"key":"type","operator":"MATCH","values":["ASSET"]}]}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSON(test.jsonBody)
			if err != nil {
				panic(err)
			}
			str, err := p.ReplaceCharInKey("_", "").Release()
			if err != nil && !test.wantError {
				panic(err)
			}
			if str != test.expected {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}

func TestMJson_IntToString(t *testing.T) {
	tests := map[string]struct {
		jsonBody  string
		wantError bool
		expected  string
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
}`, expected: `{"filters":[{"key":"type","operator":"MATCH","values":["ASSET"]},{"key":"type","operator":"MATCH","values":["ASSET"]}],"order_by":[{"key":"type"}],"paging":{"offset":"0","size":"20"}}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSON(test.jsonBody)
			if err != nil {
				panic(err)
			}
			str, err := p.IntToString().Release()
			if err != nil && !test.wantError {
				panic(err)
			}
			if str != test.expected {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}

func TestMJson_IntConverter(t *testing.T) {
	tests := map[string]struct {
		jsonBody  string
		wantError bool
		expected  string
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
}`, expected: `{"filters":[{"key":"type","operator":"MATCH","values":["ASSET"]},{"key":"type","operator":"MATCH","values":["ASSET"]}],"order_by":[{"key":"type"}],"paging":{"offset":"replace","size":"replace"}}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSON(test.jsonBody)
			if err != nil {
				panic(err)
			}
			str, err := p.IntConverter(func(value int64) interface{} {
				return "replace"
			}).Release()
			if err != nil && !test.wantError {
				panic(err)
			}
			if str != test.expected {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}

func TestMJson_ReleaseJson(t *testing.T) {
	tests := map[string]struct {
		jsonBody  string
		wantError bool
		expected  interface{}
	}{
		`object Json `: {jsonBody: `{
    "paging": {
        "offset": 1.1,
        "size": 20
    }
}`, expected: map[string]interface{}{
			`paging`: map[string]interface{}{
				`offset`: "1.1",
				`size`:   "20",
			},
		}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSON(test.jsonBody)
			if err != nil {
				panic(err)
			}
			str := p.FloatToString().ReleaseJson()
			if !reflect.DeepEqual(str, test.expected) {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}

func TestMJson_ValueChecker(t *testing.T) {
	tests := map[string]struct {
		jsonBody  string
		wantError bool
		expected  interface{}
	}{
		`object Json `: {jsonBody: `{
    "paging": {
        "offset": 1.1,
        "size": 20
    }
}`, expected: map[string]interface{}{
			`paging`: map[string]interface{}{
				`offset`: `20`,
				`size`:   `20`,
			},
		}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSON(test.jsonBody)
			if err != nil {
				panic(err)
			}
			p.ValueChecker(func(value interface{}) interface{} {
				return `20`
			})
			str := p.ReleaseJson()
			if !reflect.DeepEqual(str, test.expected) {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}

func TestMJson_FloatToString(t *testing.T) {
	tests := map[string]struct {
		jsonBody  string
		wantError bool
		expected  string
	}{
		`object Json `: {jsonBody: `{
    "paging": {
        "offset": 1.1,
        "size": 20
    }
}`, expected: `{"paging":{"offset":"1.1","size":"20"}}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSON(test.jsonBody)
			if err != nil {
				panic(err)
			}
			str, err := p.FloatToString().Release()
			if err != nil && !test.wantError {
				panic(err)
			}
			if str != test.expected {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}

func TestMJson_FloatToInt(t *testing.T) {
	tests := map[string]struct {
		jsonBody  string
		wantError bool
		expected  string
	}{
		`object Json `: {jsonBody: `{
    "paging": {
        "offset": 1.1,
        "size": 20
    }
}`, expected: `{"paging":{"offset":1,"size":20}}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSON(test.jsonBody)
			if err != nil {
				panic(err)
			}
			str, err := p.FloatToInt().Release()
			if err != nil && !test.wantError {
				panic(err)
			}
			if str != test.expected {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}

func TestMJson_FloatConverter(t *testing.T) {
	tests := map[string]struct {
		jsonBody  string
		wantError bool
		expected  string
	}{
		`object Json `: {jsonBody: `{
    "paging": {
        "offset": 1.1,
        "size": 20
    }
}`, expected: `{"paging":{"offset":"float convert","size":"float convert"}}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSON(test.jsonBody)
			if err != nil {
				panic(err)
			}
			str, err := p.FloatConverter(func(value float64) interface{} {
				return "float convert"
			}).Release()
			if err != nil && !test.wantError {
				panic(err)
			}
			if str != test.expected {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}

func TestMJson_NewMagicJSONInterface(t *testing.T) {
	tests := map[string]struct {
		jsonInter map[string]interface{}
		wantError bool
		expected  string
	}{
		`object Json `: {jsonInter: map[string]interface{}{
			`paging`: map[string]interface{}{
				`offset`: 1,
				`size`:   2,
			},
		}, expected: `{"paging":{"offset":1,"size":2}}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := mjson.NewMagicJSONInterface(test.jsonInter)
			if err != nil {
				panic(err)
			}
			str, err := p.Release()
			if err != nil && !test.wantError {
				panic(err)
			}
			if str != test.expected {
				t.Fatal(fmt.Sprintf(`got : %s, \n and expected :%s`, str, test.expected))
			}
		})
	}
}
