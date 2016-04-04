package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	Help bool   `cli:"h,help" usage:"display help information"`
	Name string `cli:"name" usage:"tell me what your name is, default name is world" dft:"world"`
	Age  uint8  `cli:"a,age" usage:"tell me your age, default is 100" dft:"100"`
}

func main() {
	cli.Run(&argT{}, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Help {
			ctx.WriteUsage()
		} else {
			ctx.String("Hello, %s! Your age is %d?\n", argv.Name, argv.Age)
		}
		return nil
	})
}
