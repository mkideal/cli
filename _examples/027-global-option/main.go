package main

import (
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

func main() {
	if err := cli.Root(root, cli.Tree(sub)).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

type rootT struct {
	cli.Helper
	Hello string `cli:"H" usage:"hello is a global flag"`
}

type subT struct {
	World string `cli:"w" usage:"world is a sub flag"`
}

var root = &cli.Command{
	Name:   "app",
	Desc:   "application",
	Global: true,
	Argv:   func() interface{} { return new(rootT) },
	Fn: func(ctx *cli.Context) error {
		ctx.JSONIndentln(ctx.RootArgv(), "", "    ")
		ctx.JSONIndentln(ctx.Argv(), "", "    ")
		return nil
	},
}

var sub = &cli.Command{
	Name: "sub",
	Desc: "subcommand",
	Argv: func() interface{} { return new(subT) },
	Fn: func(ctx *cli.Context) error {
		ctx.JSONIndentln(ctx.RootArgv(), "", "    ")
		ctx.JSONIndentln(ctx.Argv(), "", "    ")

		var argv = &subT{}
		var parentArgv = &rootT{}
		if err := ctx.GetArgvList(argv, parentArgv); err != nil {
			return err
		}
		ctx.JSONIndentln(parentArgv, "", "    ")
		ctx.JSONIndentln(argv, "", "    ")
		return nil
	},
}
