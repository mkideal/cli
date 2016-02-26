package cli

import (
	"reflect"
	"strings"
)

const (
	tagCli    = "cli"
	tagUsage  = "usage"
	tagDefaut = "dft"

	dashOne = "-"
	dashTwo = "--"

	sepName     = ", "
	sepColSpace = 3
)

type cliTag struct {
	required     bool
	shortNames   []string
	longNames    []string
	usage        string
	defaultValue string
}

func parseTag(fieldName string, tag reflect.StructTag) (*cliTag, error) {
	clitag := &cliTag{
		shortNames: []string{},
		longNames:  []string{},
	}
	cli := tag.Get(tagCli)
	clitag.usage = tag.Get(tagUsage)
	clitag.defaultValue = tag.Get(tagDefaut)

	cli = strings.TrimSpace(cli)
	for strings.HasPrefix(cli, "*") {
		clitag.required = true
		cli = strings.TrimSpace(strings.TrimPrefix(cli, "*"))
	}

	names := strings.Split(cli, ",")
	isEmpty := true
	for _, name := range names {
		if name = strings.TrimSpace(name); name == "-" {
			return nil, nil
		}
		if len(name) == 0 {
			continue
		} else if len(name) == 1 {
			clitag.shortNames = append(clitag.shortNames, dashOne+name)
		} else {
			clitag.longNames = append(clitag.longNames, dashTwo+name)
		}
		isEmpty = false
	}
	if isEmpty {
		clitag.longNames = append(clitag.longNames, dashTwo+fieldName)
	}
	//TODO: validate tags
	return clitag, nil
}
