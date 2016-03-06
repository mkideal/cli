package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mkideal/cli"
)

func main() {
	if err := cli.Root(root,
		cli.Tree(help),
		cli.Tree(daemon),
		cli.Tree(ping),
		cli.Tree(api,
			cli.Tree(build),
			cli.Tree(install),
		),
	).Run(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}

//------
// root
//------
var root = &cli.Command{
	Fn: func(ctx *cli.Context) error {
		ctx.String(ctx.Usage())
		return nil
	},
}

//------
// help
//------
var help = &cli.Command{
	Name:        "help",
	Desc:        "display help",
	CanSubRoute: true,
	HTTPRouters: []string{"/help", "/v1/help"},
	HTTPMethods: []string{http.MethodGet},

	Fn: func(ctx *cli.Context) error {
		parent := ctx.Command().Parent()
		if len(ctx.Args()) == 0 {
			ctx.String(parent.Usage())
			return nil
		}
		child := parent.Route(ctx.Args())
		if child == nil {
			cmd := strings.Join(ctx.Args(), " ")
			return fmt.Errorf("command %s not found", ctx.Color().Yellow(cmd))
		}
		ctx.String(child.Usage(ctx))
		return nil
	},
}

//--------
// daemon
//--------
type daemonT struct {
	cli.Helper
	Port uint16 `cli:"p,port" usage:"http port" dft:"8080"`
}

func (t *daemonT) Validate() error {
	if t.Port == 0 {
		return fmt.Errorf("please don't use 0 as http port")
	}
	return nil
}

var daemon = &cli.Command{
	Name: "daemon",
	Desc: "startup app as daemon",
	Argv: func() interface{} { return new(daemonT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*daemonT)
		if argv.Help {
			ctx.String(ctx.Usage())
			return nil
		}
		cli.EnableDebug()
		addr := fmt.Sprintf(":%d", argv.Port)
		ctx.String("http addr: %s\n", addr)

		if err := ctx.Command().Root().RegisterHTTP(ctx); err != nil {
			return err
		}
		return http.ListenAndServe(addr, ctx.Command().Root())
	},
}

//------
// ping
//------
var ping = &cli.Command{
	Name: "ping",
	Desc: "ping server",
	Fn: func(ctx *cli.Context) error {
		ctx.String("pong\n")
		return nil
	},
}

//-----
// api
//-----
var api = &cli.Command{
	Name: "api",
	Desc: "display all api",
	Fn: func(ctx *cli.Context) error {
		ctx.String("Commands:\n")
		ctx.String("    build\n")
		ctx.String("    install\n")
		return nil
	},
}

//-------
// build
//-------
type buildT struct {
	cli.Helper
	Dir string `cli:"dir" usage:"dest path" dft:"./"`
}

var build = &cli.Command{
	Name: "build",
	Desc: "build application",
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

//---------
// install
//---------
var install = &cli.Command{
	Name: "install",
	Desc: "install application",
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
