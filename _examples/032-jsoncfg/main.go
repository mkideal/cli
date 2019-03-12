package main

import (
	"os"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Self *argT `json:"-" cli:"c,config" usage:"parse json config from file" parser:"jsoncfg" dft:"test.json"`
	A    string
	B    int
	C    bool
}

func newArgT() *argT {
	var argv = new(argT)
	argv.Self = argv
	return argv
}

func main() {
	os.Exit(cli.Run(newArgT(), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.JSONIndentln(argv, "", "    ")
		return nil
	}))
}
