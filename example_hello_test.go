package cli_test

import (
	"github.com/mkideal/cli"
)

// This is a HelloWorld example
func Example_hello() {
	type argT struct {
		cli.Helper
		Name string `cli:"name" usage:"tell me your name" dft:"world"`
		Age  uint8  `cli:"a,age" usage:"tell me your age" dft:"100"`
	}

	args := []string{"clitest", "--name=Cliper"}
	cli.RunWithArgs(&argT{}, args, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("Hello, %s! Your age is %d?\n", argv.Name, argv.Age)
		return nil
	})
	// Output:
	// Hello, Cliper! Your age is 100?
}
