package cli

import (
	"reflect"
	"strings"
)

const (
	tagCli    = "cli"
	tagUsage  = "usage"
	tagDefaut = "dft"
	tagName   = "name"
	tagPw     = "pw" // password
	tagPrompt = "prompt"
	tagParser = "parser"

	dashOne = "-"
	dashTwo = "--"

	sepName = ", "
)

type tagProperty struct {
	// is a required flag?(if `cli` or `pw` tag has prefix `*`)
	required bool

	// is a force flag?(if `cli` or `pw` tag has prefix `!`)
	isForce bool

	// flag names
	shortNames []string
	longNames  []string

	usage        string
	defaultValue string
	name         string
	prompt       string
	isPassword   bool

	parserCreator FlagParserCreator
}

func parseTag(fieldName string, tag reflect.StructTag) (*tagProperty, bool) {
	p := &tagProperty{
		shortNames: []string{},
		longNames:  []string{},
	}
	cli := tag.Get(tagCli)
	pw := tag.Get(tagPw)
	if pw != "" {
		p.isPassword = true
		cli = pw
	}
	p.usage = tag.Get(tagUsage)
	p.defaultValue = tag.Get(tagDefaut)
	p.name = tag.Get(tagName)
	p.prompt = tag.Get(tagPrompt)
	if parserName := tag.Get(tagParser); parserName != "" {
		if parserCreator, ok := parserCreators[parserName]; ok {
			p.parserCreator = parserCreator
		}
	}

	cli = strings.TrimSpace(cli)
	for {
		if strings.HasPrefix(cli, "*") {
			p.required = true
			cli = strings.TrimSpace(strings.TrimPrefix(cli, "*"))
		} else if strings.HasPrefix(cli, "!") {
			p.isForce = true
			cli = strings.TrimSpace(strings.TrimPrefix(cli, "!"))
		} else {
			break
		}
	}

	names := strings.Split(cli, ",")
	isEmpty := true
	for _, name := range names {
		if name = strings.TrimSpace(name); name == dashOne {
			return nil, false
		}
		if len(name) == 0 {
			continue
		} else if len(name) == 1 {
			p.shortNames = append(p.shortNames, dashOne+name)
		} else {
			p.longNames = append(p.longNames, dashTwo+name)
		}
		isEmpty = false
	}
	if isEmpty {
		p.longNames = append(p.longNames, dashTwo+fieldName)
	}
	return p, isEmpty
}
