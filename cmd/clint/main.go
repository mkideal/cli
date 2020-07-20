package main

import (
	"os"

	"github.com/akeylesslabs/cli"
)

type argT struct {
	cli.Helper
	Version bool `cli:"!v,version" usage:"display version info"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Version {
			ctx.String("%v\n", appVersion)
			return nil
		}
		return nil
	}))
}
