package main

import (
	"fmt"

	"github.com/mkideal/cli"
)

var _ = app.Register(&cli.Command{
	Name:   "test",
	Desc:   "Test golang application",
	ArgvFn: func() interface{} { return new(test_t) },
	Fn:     test,
})

type test_t struct {
	Help   bool   `cli:"h,help" usage:"display help information"`
	Dir    string `cli:"dir" usage:"source code root dir" dft:"./"`
	Suffix string `cli:"suffix" usage:"source file suffix" dft:".go,.c,.s"`
	Out    string `cli:"o,out" usage:"output filename"`
}

func test(ctx *cli.Context) error {
	argv := ctx.Argv().(*test_t)

	if argv.Help {
		fmt.Fprintf(ctx.Writer(), ctx.Command().Usage())
		return nil
	}
	fmt.Fprintf(ctx.Writer(), "%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}
