package cli_test

import (
	"github.com/mkideal/cli"
)

type helloT struct {
	cli.Helper
	Name string `cli:"name" usage:"tell me your name" dft:"world"`
	Age  uint8  `cli:"a,age" usage:"tell me your age" dft:"100"`
}

// This is a HelloWorld example
func Example_hello() {
	args := []string{"app", "--name=Cliper"}
	cli.RunWithArgs(new(helloT), args, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*helloT)
		ctx.String("Hello, %s! Your age is %d?\n", argv.Name, argv.Age)
		return nil
	})
	// Output:
	// Hello, Cliper! Your age is 100?
}
