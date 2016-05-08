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

type fieldTag struct {
	required      bool
	shortNames    []string
	longNames     []string
	usage         string
	defaultValue  string
	name          string
	prompt        string
	isPassword    bool
	parserCreator FlagParserCreator

	isHelp bool
}

func parseTag(fieldName string, tag reflect.StructTag) (*fieldTag, bool) {
	ftag := &fieldTag{
		shortNames: []string{},
		longNames:  []string{},
	}
	cli := tag.Get(tagCli)
	pw := tag.Get(tagPw)
	if pw != "" {
		ftag.isPassword = true
		cli = pw
	}
	ftag.usage = tag.Get(tagUsage)
	ftag.defaultValue = tag.Get(tagDefaut)
	ftag.name = tag.Get(tagName)
	ftag.prompt = tag.Get(tagPrompt)
	if parserName := tag.Get(tagParser); parserName != "" {
		if parserCreator, ok := parserCreators[parserName]; ok {
			ftag.parserCreator = parserCreator
		}
	}

	cli = strings.TrimSpace(cli)
	for {
		if strings.HasPrefix(cli, "*") {
			ftag.required = true
			cli = strings.TrimSpace(strings.TrimPrefix(cli, "*"))
		} else if strings.HasPrefix(cli, "!") {
			ftag.isHelp = true
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
			ftag.shortNames = append(ftag.shortNames, dashOne+name)
		} else {
			ftag.longNames = append(ftag.longNames, dashTwo+name)
		}
		isEmpty = false
	}
	if isEmpty {
		ftag.longNames = append(ftag.longNames, dashTwo+fieldName)
	}
	return ftag, isEmpty
}
