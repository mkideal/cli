package cli

import (
	"reflect"
	"strings"
)

const (
	tagCli  = "cli"
	tagPw   = "pw" // password
	tagEdit = "edit"

	tagUsage  = "usage"
	tagDefaut = "dft"
	tagName   = "name"
	tagPrompt = "prompt"
	tagParser = "parser"
	tagSep    = "sep" // used to seperate key/value pair of map, default is `=`

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
	isEdit       bool
	editFile     string

	sep string

	parserCreator FlagParserCreator
}

func parseTag(fieldName string, tag reflect.StructTag) (*tagProperty, bool) {
	p := &tagProperty{
		shortNames: []string{},
		longNames:  []string{},
	}
	cliLikeTagCount := 0

	// `cli` TAG
	cli := tag.Get(tagCli)
	if cli != "" {
		cliLikeTagCount++
	}

	// `pw` TAG
	if pw := tag.Get(tagPw); pw != "" {
		p.isPassword = true
		cli = pw
		cliLikeTagCount++
	}

	// `edit` TAG
	if edit := tag.Get(tagEdit); edit != "" {
		// specific filename for editor
		sepIndex := strings.Index(edit, ":")
		if sepIndex > 0 {
			p.editFile = edit[:sepIndex]
			edit = edit[sepIndex+1:]
		}
		p.isEdit = true
		cli = edit
		cliLikeTagCount++
	}

	if cliLikeTagCount > 1 {
		return nil, false
	}

	// `usage` TAG
	p.usage = tag.Get(tagUsage)

	// `dft` TAG
	p.defaultValue = tag.Get(tagDefaut)

	// `name` TAG
	p.name = tag.Get(tagName)

	// `prompt` TAG
	p.prompt = tag.Get(tagPrompt)

	// `parser` TAG
	if parserName := tag.Get(tagParser); parserName != "" {
		if parserCreator, ok := parserCreators[parserName]; ok {
			p.parserCreator = parserCreator
		}
	}

	// `sep` TAG
	p.sep = "="
	if sep := tag.Get(tagSep); sep != "" {
		p.sep = sep
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
