package mjson

// MagicJSON keep the information about given json string behavior
// json can be release as string at anytime
type MagicJSON interface {
	JSONRelease
	JSONConverter

	// Key traversal inserted json object
	// until find the input key
	// there can be multiple keys with same name
	// but that will be discarded and select the first searched key
	Key(key string) JSONConverter
}

// JSONConverter conversion of json key and value
type JSONConverter interface {
	StringConverter
	IntegerValueConverter
	FloatValueConverter
}

// StringConverter Convert key or value which is in string format
type StringConverter interface {
	// ReplaceKeyCharacter string key is replaced by the given characters
	// new characters will be replaced if old characters can be found in the json key
	// only replaced single match
	ReplaceKeyCharacter(oldCharacters, newCharacters string) JSONRelease

	// ReplaceCharsInKey string key is replaced by the given characters
	// new characters will be replaced if old characters can be found in the json key
	// replace with count of the match
	ReplaceCharsInKey(oldCharacters, newCharacters string, count int) JSONRelease

	// ReplaceCharInValue string value is replaced by the given characters
	// new characters will be replaced if old characters can be found in the json value
	// only replaced single match
	ReplaceCharInValue(oldCharacters, newCharacters string) JSONRelease

	// ReplaceValueCharacters string value is replaced by the given characters
	// new characters will be replaced if old characters can be found in the json value
	// replace with count of the match
	ReplaceValueCharacters(oldCharacters, newCharacters string, count int) JSONRelease

	// ValueStringToInt convert string values into integers
	// if that is possible to do the conversion
	ValueStringToInt() JSONRelease

	// ValueStringToFloat convert string values into integers
	// if that is possible to do the conversion
	ValueStringToFloat() JSONRelease

	// ValueStringConverter can provide input as a function
	// bind any value to existing position
	ValueStringConverter(fn func(value string) interface{}) JSONRelease
}

// IntegerValueConverter convert value of integer into any format
type IntegerValueConverter interface {
	// IntToString convert int64 value into string
	IntToString() JSONRelease

	// IntConverter can provide input as a function
	// bind any value to existing position
	IntConverter(fn func(value int64) interface{}) JSONRelease
}

// FloatValueConverter convert value of integer into any format
type FloatValueConverter interface {
	// FloatToString convert int64 value into string
	FloatToString() JSONRelease

	// FloatToInt convert int64 value into integer
	FloatToInt() JSONRelease

	// FloatConverter can provide input as a function
	// bind any value to existing position
	FloatConverter(fn func(value float64) interface{}) JSONRelease
}

// JSONRelease release the json object as a string
type JSONRelease interface {
	Release() (json string, err error)
}
