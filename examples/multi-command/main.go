package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

const VERSION = "v1.0.0"

var app = cli.NewWithCommand(&cli.Command{
	Name:   os.Args[0],
	Desc:   "Golang package manager",
	ArgvFn: func() interface{} { return new(gogo_t) },
	Fn:     gogo,
}, os.Stdout)

type gogo_t struct {
	Help    bool `cli:"h,help" usage:"display help information"`
	Version bool `cli:"v,version" usage:"display version"`
	List    bool `cli:"l,list" usage:"list all sub commands or not" dft:"false"`
}

func gogo(ctx *cli.Context) error {
	argv := ctx.Argv().(*gogo_t)
	if argv.Help {
		fmt.Fprintf(ctx.Writer(), ctx.Command().Usage())
		return nil
	}
	if argv.Version {
		fmt.Fprintf(ctx.Writer(), VERSION+"\n")
		return nil
	}

	if argv.List {
		length := 0
		for _, cmd := range ctx.Command().Children() {
			if len(cmd.Name) > length {
				length = len(cmd.Name)
			}
		}
		format := fmt.Sprintf("%%-%ds%%s\n", length+4)
		for _, cmd := range ctx.Command().Children() {
			fmt.Fprintf(ctx.Writer(), format, cmd.Name, cmd.Desc)
		}
		return nil
	}

	fmt.Fprintf(ctx.Writer(), "try `%s --help for more information'\n", ctx.Path())
	return nil
}

func jsonIndent(i interface{}) string {
	data, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		return ""
	}
	return string(data) + "\n"
}

func main() {
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Printf("%v\n", err)
	}
}
