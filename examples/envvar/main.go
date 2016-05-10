package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	String string  `cli:"s,string" usage:"string env var" dft:"$CLI_ENVVAR_STRING/hello"`
	Int    int     `cli:"i,int" usage:"int env var" dft:"$CLI_ENVVAR_INT*2+13"`
	Float  float64 `cli:"f,float" usage:"float env var" dft:"$CLI_ENVVAR_FLOAT*2.5-0.5"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		ctx.JSONIndentln(ctx.Argv(), "", "    ")
		return nil
	})
}
