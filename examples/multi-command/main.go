package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

const version = "v1.0.0"

var app = &cli.Command{
	Name: os.Args[0],
	Desc: "Golang package manager",
	Text: `gogo is a new golang package manager
very very good`,
	Argv: func() interface{} { return new(gogoT) },
	Fn:   gogo,
}

type gogoT struct {
	Help    bool `cli:"h,help" usage:"display help information"`
	Version bool `cli:"v,version" usage:"display version"`
	List    bool `cli:"l,list" usage:"list all sub commands or not" dft:"false"`
}

func gogo(ctx *cli.Context) error {
	argv := ctx.Argv().(*gogoT)
	if argv.Help {
		ctx.String(ctx.Command().Usage())
		return nil
	}
	if argv.Version {
		ctx.String(version + "\n")
		return nil
	}

	if argv.List {
		ctx.String(ctx.Command().ListChildren(" ", "  =>  "))
		return nil
	}

	ctx.String("try `%s --help for more information'\n", ctx.Path())
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
	//NOTE: You can set any writer implements io.Writer
	// default writer is os.Stdout
	app.SetWriter(os.Stderr)

	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Printf("%v\n", err)
	}
}
