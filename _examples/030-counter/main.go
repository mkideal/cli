package main

import (
	"os"

	"github.com/akeylesslabs/cli"
)

type argT struct {
	cli.Helper
	V cli.Counter `cli:"v" usage:"count verbose"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("v=%d\n", argv.V.Value())
		return nil
	}))
}
