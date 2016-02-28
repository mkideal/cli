package main

import (
	"github.com/mkideal/cli"
)

var _ = publishCmd.Register(&cli.Command{
	Name: "cn",
	Desc: "Publish golang application to CN",
	Argv: func() interface{} { return new(publishCnT) },
	Fn:   publishCn,
})

type publishCnT struct {
	Help   bool   `cli:"h,help" usage:"display help information"`
	Dir    string `cli:"dir" usage:"source code root dir" dft:"./"`
	Suffix string `cli:"suffix" usage:"source file suffix" dft:".go,.c,.s"`
	Out    string `cli:"o,out" usage:"output filename"`
}

func publishCn(ctx *cli.Context) error {
	argv := ctx.Argv().(*publishCnT)

	if argv.Help {
		ctx.String(ctx.Usage())
		return nil
	}
	ctx.String("%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}
