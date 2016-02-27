package main

import (
	"fmt"

	"github.com/mkideal/cli"
)

var _ = app.Register(&cli.Command{
	Name:   "clean",
	Desc:   "Clean build data",
	ArgvFn: func() interface{} { return new(clean_t) },
	Fn:     clean,
})

type clean_t struct {
	Help      bool `cli:"h,help" usage:"display help information"`
	Recursion bool `cli:"R,recursion" usage:"clean recursion or not" dft:"true"`
}

func clean(ctx *cli.Context) error {
	argv := ctx.Argv().(*clean_t)

	if argv.Help {
		fmt.Fprintf(ctx.Writer(), ctx.Command().Usage())
		return nil
	}
	fmt.Fprintf(ctx.Writer(), "%s: %v", ctx.Path(), jsonIndent(argv))
	return nil
}
