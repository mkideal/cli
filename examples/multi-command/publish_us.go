package main

import (
	"github.com/mkideal/cli"
)

var _ = publishCmd.Register(&cli.Command{
	Name: "us",
	Desc: "Publish golang application to US",
	Argv: func() interface{} { return new(publish_us_t) },
	Fn:   publish_us,
})

type publish_us_t struct {
	Help   bool   `cli:"h,help" usage:"display help information"`
	Dir    string `cli:"dir" usage:"source code root dir" dft:"./"`
	Suffix string `cli:"suffix" usage:"source file suffix" dft:".go,.c,.s"`
	Out    string `cli:"o,out" usage:"output filename"`
}

func publish_us(ctx *cli.Context) error {
	argv := ctx.Argv().(*publish_us_t)

	if argv.Help {
		ctx.String(ctx.Usage())
		return nil
	}
	ctx.String("%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}
