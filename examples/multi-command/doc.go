package main

import (
	"github.com/mkideal/cli"
)

var _ = app.Register(&cli.Command{
	Name: "doc",
	Desc: "Generate documents",
	Argv: func() interface{} { return new(doc_t) },
	Fn:   doc,
})

type doc_t struct {
	Help   bool   `cli:"h,help" usage:"display help information"`
	Suffix string `cli:"suffix" usage:"source file suffix" dft:".go,.c,.s"`
	Out    string `cli:"o,out" usage:"output filename"`
}

func doc(ctx *cli.Context) error {
	argv := ctx.Argv().(*doc_t)

	if argv.Help {
		ctx.String(ctx.Usage())
		return nil
	}
	ctx.String("%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}
