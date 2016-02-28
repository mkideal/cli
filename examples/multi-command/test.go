package main

import (
	"github.com/mkideal/cli"
)

var _ = app.Register(&cli.Command{
	Name: "test",
	Desc: "Test golang application",
	Argv: func() interface{} { return new(testT) },
	Fn:   test,
})

type testT struct {
	Help   bool   `cli:"h,help" usage:"display help information"`
	Dir    string `cli:"dir" usage:"source code root dir" dft:"./"`
	Suffix string `cli:"suffix" usage:"source file suffix" dft:".go,.c,.s"`
	Out    string `cli:"o,out" usage:"output filename"`
}

func test(ctx *cli.Context) error {
	argv := ctx.Argv().(*testT)

	if argv.Help {
		ctx.String(ctx.Usage())
		return nil
	}
	ctx.String("%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}
