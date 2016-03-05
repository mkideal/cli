package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mkideal/cli"
)

var isClient = true

func main() {
	if err := cli.Root(root,
		cli.Tree(daemon),
		cli.Tree(ping),
		cli.Tree(app,
			cli.Tree(build),
			cli.Tree(install),
		),
	).Run(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}

type daemonT struct {
	Help bool `cli:"!h,help" usage:"Dispaly help"`
	Port int  `cli:"p,port" usage:"http port" dft:"8080"`
}

var root = &cli.Command{}

var daemon = &cli.Command{
	Name: "daemon",
	Argv: func() interface{} { return new(daemonT) },

	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*daemonT)
		if argv.Help {
			ctx.String(ctx.Usage())
			return nil
		}
		isClient = false
		cli.EnableDebug()
		addr := fmt.Sprintf(":%d", argv.Port)
		ctx.String("http addr: %s\n", addr)
		return http.ListenAndServe(addr, ctx)
	},
}

var ping = &cli.Command{
	Name: "ping",
	Fn: func(ctx *cli.Context) error {
		ctx.String("pong\n")
		return nil
	},
}

var app = &cli.Command{
	Name: "app",
	Fn: func(ctx *cli.Context) error {
		ctx.String("Commands:\n")
		ctx.String("    build\n")
		ctx.String("    install\n")
		return nil
	},
}

type buildT struct {
	Help bool   `cli:"h,help" usage:"Dispaly help"`
	Dir  string `cli:"dir" usage:"Dest path" dft:"./"`
}

var build = &cli.Command{
	Name: "build",
	Argv: func() interface{} { return new(buildT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*buildT)
		if argv.Help {
			ctx.String(ctx.Usage())
			return nil
		}
		ctx.JSONIndentln(argv, "", "    ")
		return nil
	},
}

var install = &cli.Command{
	Name: "install",
	Argv: func() interface{} { return new(buildT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*buildT)
		if argv.Help {
			ctx.String(ctx.Usage())
			return nil
		}
		ctx.JSONIndentln(argv, "", "    ")
		return nil
	},
}
