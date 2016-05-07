package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Name string `cli:"name" usage:"tell me what your name is, default world" dft:"world"`
	Age  uint8  `cli:"a,age" usage:"tell me your age, default 100" dft:"100"`
}

func main() {
	cli.Run(&argT{}, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("Hello, %s! Your age is %d?\n", argv.Name, argv.Age)
		return nil
	})
}
