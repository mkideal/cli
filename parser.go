package cli

import (
	"encoding/json"
	"net/url"
)

// FlagParser represents a parser for parsing flag
type FlagParser interface {
	Parse(s string) error
}

// FlagParserCreator represents factory function of FlagParser
type FlagParserCreator func(ptr interface{}) FlagParser

var parserCreators = map[string]FlagParserCreator{}

// RegisterFlagParser registers FlagParserCreator by name
func RegisterFlagParser(name string, creator FlagParserCreator) {
	if _, ok := parserCreators[name]; ok {
		panic("RegisterFlagParser has registered: " + name)
	}
	parserCreators[name] = creator
}

func init() {
	RegisterFlagParser("json", newJSONParser)
	RegisterFlagParser("jsonfile", newJSONFileParser)
	RegisterFlagParser("jsoncfg", newJSONConfigFileParser)
	RegisterFlagParser("url", newURLParser)
}

// JSON parser
type JSONParser struct {
	ptr interface{}
}

func newJSONParser(ptr interface{}) FlagParser {
	return &JSONParser{ptr}
}

func (p JSONParser) Parse(s string) error {
	return json.Unmarshal([]byte(s), p.ptr)
}

// JSON file parser
type JSONFileParser struct {
	ptr interface{}
}

func newJSONFileParser(ptr interface{}) FlagParser {
	return &JSONFileParser{ptr}
}

func (p JSONFileParser) Parse(s string) error {
	return ReadJSONFromFile(s, p.ptr)
}

// JSON config file parser
type JSONConfigFileParser struct {
	ptr interface{}
}

func newJSONConfigFileParser(ptr interface{}) FlagParser {
	return &JSONConfigFileParser{ptr}
}

func (p JSONConfigFileParser) Parse(s string) error {
	return ReadJSONConfigFromFile(s, p.ptr)
}

// URL parser
type URLParser struct {
	ptr *url.URL
}

func newURLParser(ptr interface{}) FlagParser {
	return &URLParser{ptr: ptr.(*url.URL)}
}

func (p *URLParser) Parse(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	*p.ptr = *u
	return nil
}
