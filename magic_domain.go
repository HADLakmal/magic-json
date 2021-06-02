package mjson

type JSONLoader interface {
	Load(json string) (JSONConverter, error)
	Release() (json string, err error)
}

type JSONConverter interface {
	ReplaceKey(oldCharacters, newCharacters string)
	ReplaceValue(oldCharacters, newCharacters string)
}
