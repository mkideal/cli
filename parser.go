package cli

import (
	"encoding/json"
)

type FlagParser interface {
	Parse(s string) error
}

type FlagParserCreator func(interface{}) FlagParser

var parserCreators = map[string]FlagParserCreator{}

func RegisterFlagParser(name string, creator FlagParserCreator) {
	if _, ok := parserCreators[name]; ok {
		return
	}
	parserCreators[name] = creator
}

func init() {
	RegisterFlagParser("json", newJSONParser)
}

// JSON parser
type JSONParser struct {
	i interface{}
}

func newJSONParser(i interface{}) FlagParser {
	return &JSONParser{i}
}

func (p JSONParser) Parse(s string) error {
	return json.Unmarshal([]byte(s), p.i)
}
