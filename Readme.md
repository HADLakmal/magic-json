# Magic Json Library

Magic Json has ability to change the key or value in the JSON. JSON is treated as a list of node to accommodate the
library requirements.

Service is capable of,

* Change or Replace the characters of JSON key
* Change characters of a JSON value
* Replace JSON value by entire different type of value
* Travel to unique key object and change behaviours as describe above

# Getting Start

### Installing

To start using Magic Json, install Go and run go get:

```sh
$ go get -u github.com/HADLakmal/magic-json
```

This will fetch the library.

### Example Usage

Load json to magic json library and release the unchanged json

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

Remove `_` from JSON keys

```go
package main

import "github.com/HADLakmal/magic-json"

const mJson = `{"first_name":"magic","last_name":"json","age":21}`

func main() {
	m, err := mjson.NewMagicJSON(mJson)
	if err != nil {
		panic(err)
	}
	r, _ := m.ReplaceKeyCharacter("_", "").Release()
	println(r)
}
```