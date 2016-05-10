package cli

import (
	"encoding/json"
	"io/ioutil"
)

// FlagParser represents a parser for parsing flag
type FlagParser interface {
	Parse(s string) error
}

// FlagParserCreator represents factory function of FlagParser
type FlagParserCreator func(interface{}) FlagParser

var parserCreators = map[string]FlagParserCreator{}

// RegisterFlagParser registers FlagParserCreator by name
func RegisterFlagParser(name string, creator FlagParserCreator) {
	if _, ok := parserCreators[name]; ok {
		panic("RegisterFlagParser has registered " + name)
	}
	parserCreators[name] = creator
}

func init() {
	RegisterFlagParser("json", newJSONParser)
	RegisterFlagParser("jsonfile", newJSONFileParser)
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

// JSON file parser
type JSONFileParser struct {
	i interface{}
}

func newJSONFileParser(i interface{}) FlagParser {
	return &JSONFileParser{i}
}

func (p JSONFileParser) Parse(s string) error {
	data, err := ioutil.ReadFile(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, p.i)
}
