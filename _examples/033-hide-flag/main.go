// main.go
// This example hides Gender and InternalUsage flags.
package main

import (
	"os"

	"github.com/mkideal/cli"
)

type helloT struct {
	cli.Helper
	Name          string `cli:"name" usage:"tell me your name" dft:"world"`
	Gender        string `cli:"-"` // deprecated
	InternalUsage string `cli:"-"` // hide
	Age           uint8  `cli:"a,age" usage:"tell me your age" dft:"100"`
}

func main() {
	os.Exit(cli.Run(new(helloT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*helloT)
		ctx.String("Hello, %s! Your age is %d?\n", argv.Name, argv.Age)
		return nil
	}))
}
