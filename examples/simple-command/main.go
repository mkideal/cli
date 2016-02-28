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
				fmt.Printf(ctx.Command().Usage())
			} else {
				fmt.Printf("argv=%v\n", *argv)
			}
			return nil
		},
	}

	app.RegisterFunc("help", func(ctx *cli.Context) error {
		fmt.Println(`show help: sub commands: help/version`)
		return nil
	}, func() interface{} {
		return new(help_t)
	})

	app.Register(&cli.Command{
		// NOTE: Name is required, panic if ""
		Name: "version",

		// NOTE: ArgvFn is required, panic if nil
		Argv: func() interface{} { return new(version_t) },

		// NOTE: Fn is required, panic if nil
		Fn: func(ctx *cli.Context) error {
			fmt.Println(`version: v0.0.1`)
			return nil
		},

		Desc: "Desc is optional",
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

type help_t struct {
}

type version_t struct {
}
