# Magic Json Library

Magic Json has ability to change the key or value in the JSON. JSON is treated as a list of node to accommodate the
library requirements.

Service is capable of,

* Change or Replace the characters of JSON key
* Change characters of a JSON value
* Replace JSON value by entire different type of value
* Check value through the function
* Travel to unique key object and change behaviours as describe above

# Getting Start

### Example Usage

#### Load JSON String

##### Result as String

Load json string to library and release the unchanged json string

```go
package main

import "github.com/HADLakmal/magic-json"

const mJson = `{"name":{"first":"magic","last":"json"},"age":21}`

func main() {
	m, err := mjson.NewMagicJSON(mJson)
	if err != nil {
		// error happen due to mJson is not in proper format
		panic(err)
	}
	// release the loaded json as string
	r, _ := m.Release()
	println(r)
}
```

##### Result as Json

Load json string to library and release the unchanged json

```go
package main

import "github.com/HADLakmal/magic-json"

const mJson = `{"name":{"first":"magic","last":"json"},"age":21}`

func main() {
	m, err := mjson.NewMagicJSON(mJson)
	if err != nil {
		// error happen due to mJson is not in proper format
		panic(err)
	}
	// release the loaded json as string
	r := m.ReleaseJson()
	println(r)
}
```

#### Load JSON Interface

##### Result as String

Load json string to library and release the unchanged json string

```go
package main

import "github.com/HADLakmal/magic-json"

var mJson = map[string]interface{}{
	`name`: map[string]interface{}{
		`first`: `magic`,
		`last`:  `json`,
	},
}

func main() {
	m, err := mjson.NewMagicJSONInterface(mJson)
	if err != nil {
		// error happen due to mJson is not in proper format
		panic(err)
	}
	// release the loaded json as string
	r, _ := m.Release()
	println(r)
}
```

##### Result as Json

Load json string to library and release the unchanged json

```go
package main

import "github.com/HADLakmal/magic-json"

var mJson = map[string]interface{}{
	`name`: map[string]interface{}{
		`first`: `magic`,
		`last`:  `json`,
	},
}

func main() {
	m, err := mjson.NewMagicJSONInterface(mJson)
	if err != nil {
		// error happen due to mJson is not in proper format
		panic(err)
	}
	// release the loaded json as string
	r := m.ReleaseJson()
	println(r)
}
```

#### JSON Key Modify

Replace character in JSON keys. If you want to replace multiple characters then you can use ```ReplaceCharsInKey()```
with exact count. Below example replace one character in two different ways.

```go
package main

import (
	"github.com/HADLakmal/magic-json"
	"strings"
)

const mJson = `{"first_name":"magic","last_name":"json","age":21}`

func main() {
	m, err := mjson.NewMagicJSON(mJson)
	if err != nil {
		panic(err)
	}
	// remove _ from key
	r, _ := m.ReplaceCharInKey("_", "").Release()
	println(r)

	// key replace by input function
	mf, err := mjson.NewMagicJSON(mJson)
	if err != nil {
		panic(err)
	}
	// remove _ key
	fr, _ := mf.KeyStringConverter(func(key string) string {
		return strings.Replace(key, "_", "", 1)
	}).Release()
	println(fr)
}
```

#### JSON Value Modify

Replace character in JSON value. Replace will be happened if value filed is string and value should contain the input
character. Below example replace one character in two different ways.

```go
package main

import (
	"github.com/HADLakmal/magic-json"
	"strings"
)

const mJson = `{"first_name":"magic_json","last_name":"library","age":21}`

func main() {
	m, err := mjson.NewMagicJSON(mJson)
	if err != nil {
		panic(err)
	}
	// remove _ from value
	r, _ := m.ReplaceCharInValue("_", "").Release()
	println(r)

	// value replace by input function
	mf, err := mjson.NewMagicJSON(mJson)
	if err != nil {
		panic(err)
	}
	// remove _ from value
	rf, _ := mf.ValueStringConverter(func(value string) interface{} {
		return strings.Replace(value, "_", "", 1)
	}).Release()
	println(rf)
}
```

#### Library Functions

| Function      | Return        | Explain                                                                                                                                                                                     |
| -------------------- | ------------- |---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ```Key(key string)```  | JSONConverter  | find specific JSON key then your changes will affected under this JSON object/array. There can be multiple keys with same name but that will be discarded and select the first searched key |
| ```ReplaceCharInKey(oldCharacters, newCharacters string)```  | JSONRelease  | string key is replaced by the given characters                                                                                                                                              |
| ```ReplaceCharsInKey(oldCharacters, newCharacters string, count int)```  | JSONRelease  | string key is replaced by the given characters and replace with count of the match                                                                                                          |
| ```KeyStringConverter(fn func(value string) string)```  | JSONRelease  | Key character replacement can be provide as a function, that will bind key to desired string                                                                                                |
| ```ReplaceCharInValue(oldCharacters, newCharacters string)```  | JSONRelease  | string value is replaced by the given characters                                                                                                                                            |
| ```ReplaceCharsInValue(oldCharacters, newCharacters string, count int)```  | JSONRelease  | string value is replaced by the given characters and replace with count of the match                                                                                                        |
| ```ValueStringConverter(fn func(value string) string)```  | JSONRelease  | value character replacement can be provide as a function, that will bind key to desired string                                                                                              |
| ```ValueStringToInt()```  | JSONRelease  | convert string values into integer                                                                                                                                                          |
| ```ValueStringToFloat()```  | JSONRelease  | convert string values into float                                                                                                                                                            |
| ```IntToString()```  | JSONRelease  | convert integer value into string                                                                                                                                                           |
| ```IntConverter(fn func(value int64) interface{}) string)```  | JSONRelease  | convert integer value into any desired value                                                                                                                                                |
| ```FloatToString()```  | JSONRelease  | convert float value into string                                                                                                                                                             |
| ```FloatToInt()```  | JSONRelease  | convert float value into integer                                                                                                                                                            |
| ```FloatConverter(fn func(value float64) interface{})```  | JSONRelease  | convert float value by providing input as a function                                                                                                                                        |
| ```ValueChecker(fn func(value interface{}) interface{})```  |  JSONRelease  | check the value format input as a function                                                                                                                                                  |