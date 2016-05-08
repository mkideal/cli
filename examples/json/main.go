package main

import (
	"encoding/json"

	"github.com/mkideal/cli"
)

type jsonT struct {
	Int    int
	String string
}

func (j *jsonT) Decode(s string) error {
	return json.Unmarshal([]byte(s), j)
}

func (j *jsonT) Encode() string {
	if data, err := json.Marshal(j); err != nil {
		return ""
	} else {
		return string(data)
	}
}

type argT struct {
	cli.Helper
	JSON jsonT `cli:"json" usage:"json argument"`
}

func run(ctx *cli.Context, argv *argT) error {
	ctx.JSONIndentln(argv.JSON, "", "    ")
	return nil
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		return run(ctx, argv)
	})
}
