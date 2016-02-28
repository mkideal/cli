package main

import (
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

func main() {
	app := &cli.Command{
		Name: os.Args[0],
		Argv: func() interface{} { return new(arg_t) },
		Fn: func(ctx *cli.Context) error {
			argv := ctx.Argv().(*arg_t)
			if argv.Help {
				ctx.String(ctx.Usage())
			} else {
				ctx.String("argv=%v\n", *argv)
			}
			return nil
		},
	}

	app.RegisterFunc("help", func(ctx *cli.Context) error {
		ctx.String("show help: sub commands: help/version\n")
		return nil
	}, nil)

	app.Register(&cli.Command{
		// NOTE: Name is required, panic if ""
		Name: "version",

		// NOTE: Fn is required, panic if nil
		Fn: func(ctx *cli.Context) error {
			ctx.String("version: v0.0.1\n")
			return nil
		},

		// Argv is optional

		Desc: "Desc represent command's abstract, optional",
		Text: "Text represent command's detailed description, optional too",
	})

	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Printf("%v\n", err)
	}
}

type arg_t struct {
	Help bool   `cli:"h,help" usage:"show help"`
	Host string `cli:"H,host" usage:"specify host address" dft:"127.0.0.1"`
	Port uint16 `cli:"p,port" usage:"specify http port" dft:"8080"`
}
