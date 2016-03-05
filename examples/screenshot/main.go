package main

import (
	"fmt"
	"github.com/mkideal/cli"
	"os"
)

type argT struct {
	Help     bool   `cli:"!h,help" usage:"This is a help flag.'cli' tag must has prefix '!'"`
	Default  string `cli:"dft" usage:"This is a default flag" dft:"default value"`
	Home     string `cli:"home" usage:"This is a flag that using ENV as default" dft:"$HOME"`
	Required int    `cli:"*r" usage:"This is a required flag.'cli' tag must has prefix '*'"`
	Short    uint   `cli:"s" usage:"This is a flag that using short name '-s'"`
	Long     int8   `cli:"long" usage:"This is a flag that using long name '--long'"`
	Two      uint8  `cli:"T,two" usage:"This is a flag that supporting two names"`
	Many     int16  `cli:"x,y,xy,XY" usage:"This is a flag that supporting many names"`
	ShortDft uint16 `cli:"d" usage:"This is a default flag that using short name" dft:"65535"`
	UnName   int32  `usage:"This is a unnamed flag.Add prefix '--' as it's name"`
	Ignored  uint32 `cli:"-" usage:"This is a ignored field and you can not see it"`

	I64 int64   `cli:"i64" usage:"This is a int64 flag"`
	U64 uint64  `cli:"u64" usage:"This is a uint64 flag"`
	F32 float32 `cli:"f32" usage:"This is a float32 flag"`
	F64 float64 `cli:"f64" usage:"This is a flaot64 flag"`
}

var root = &cli.Command{
	Name: "screenshot",
	Argv: func() interface{} { return new(argT) },

	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Help {
			ctx.String(ctx.Usage())
		} else {
			ctx.String("I am very lazzy! I do nothing!\n")
		}
		return nil
	},
}

func main() {
	err := cli.Root(
		root,
		cli.Tree(help),
		cli.Tree(version),
		cli.Tree(build),
		cli.Tree(install),
	).Run(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
}

var log = func(ctx *cli.Context) error {
	ctx.String("path: %s\n", ctx.Path())
	return nil
}

var help = &cli.Command{Name: "help", Fn: log}
var version = &cli.Command{Name: "version", Fn: log}
var build = &cli.Command{Name: "build", Fn: log}
var install = &cli.Command{Name: "install", Fn: log}
