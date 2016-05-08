package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mkideal/cli"
)

func main() {
	if err := cli.Root(root,
		cli.Tree(daemon),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

const prefix = "start ok"

type argT struct {
	cli.Helper
	Echo string `cli:"e,echo" usage:"echo message"`
}

var root = &cli.Command{
	Desc: "daemon test",
	Argv: func() interface{} { return new(argT) },

	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		fmt.Fprintf(os.Stderr, "%s: %s, kill me with \"kill %d\"\n", prefix, argv.Echo, os.Getpid())

		<-time.After(time.Second * 30)
		return nil
	},
}

var daemon = &cli.Command{
	Name: "daemon",
	Desc: "startup as a daemon process",
	Argv: func() interface{} { return new(argT) },

	Fn: func(ctx *cli.Context) error {
		return cli.Daemon(ctx, prefix)
	},
}
