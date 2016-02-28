# Command line interface [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/cli/master/LICENSE)

## License
[The MIT License (MIT)](https://zh.wikipedia.org/wiki/MIT許可證)

## Install
```sh
go get github.com/mkideal/cli
```

## Features

1. Based on golang tag. Support three tags: `cli`,`usage`,`dft`
2. Support specify default value: use `dft` tag
3. Support required declaration: `cli` tag with prefix `*`
4. Support multi flag name for same field: like `cli:"h,help"`
5. Type safty

## Usage
	
First, you should define a struct, like this:
```go
type Args struct {
	OnlySingle     bool    `cli:"v" usage:"only single char"`
	ManySingle     string  `cli:"X,Y" usage:"many single char"`
	SingleAndMulti int     `cli:"s,single-and-multi" usage:"single and multi"`
	OnlyMulti      uint    `cli:"only-multi" usage:"only multi"`
	Required       int8    `cli:"*required" usage:"required value"`
	Default        uint8   `cli:"id" usage:"default value" dft:"102"`
	Ignored        int16   `cli:"-" usage:"ignored field"`
	UnName         uint16  `usage:"unname field"`
	Int32          int32   `cli:"i32" usage:"type int32"`
	Uint32         uint32  `cli:"u32" usage:"type uint32"`
	Int64          int64   `cli:"i64" usage:"type int64"`
	Uint64         int64   `cli:"u64" usage:"type uint64"`
	Float32        float32 `cli:"f32" usage:"type float32"`
	Float64        float64 `cli:"f364" usage:"type float64"`
}
```

Then, call function `cli.Parse`:
```go
t := new(Args)
flagSet := cli.Parse(os.Args[1:], t)
if flagSet.Error != nil {
	//TODO: handle the error
}
//^REMOVE: show help
// fmt.Printf("Usage of `%s'`: \n%s", os.Args[0], flagSet.Usage)
```

If you only want to show help, you can directly call function `cli.Usage`:
```go
usage := cli.Usage(new(Args))
fmt.Printf("Usage of `%s'`: \n%s", os.Args[0], usage)
```

## Tags

### cli

`cli` tag support single-char format and multi-char format, e.g.

```go
Help    bool    `cli:"h,help"`
Version string  `cli:"version"`
Port    int     `cli:"p"`
XYZ     bool    `cli:"x,y,z,xyz,XYZ"` 
```

The argument is required if `cli` tag has prefix `*`, e.g.

```go
Required string `cli:"*required"`
```

### usage

`usage` tag describe the argument. If the argument is required, describe string has prefix `*` while show usage(`*` is red on unix-like os).

### dft
`dft` tag specify argument's default value.

## Cli/Command

`Cli` and `Command` define:

```go
type Cli struct {
	root *Command
}

type Command struct {
	Name   string
	Desc   string
	Fn     CommandFunc
	ArgvFn ArgvFunc

	parent   *Command
	children []*Command
}
```

Example:

```go
package main

import (
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

func main() {
	app := cli.New(os.Args[0], nil)

	app.Root().Fn = func(ctx *cli.Context) error {
		argv := ctx.Argv().(*arg_t)
		if argv.Help {
			fmt.Printf(ctx.Command().Usage())
		} else {
			fmt.Printf("argv=%v\n", *argv)
		}
		return nil
	}
	app.Root().ArgvFn = func() interface{} {
		return new(arg_t)
	}

	app.RegisterFunc("help", func(ctx *cli.Context) error {
		fmt.Println(`show help: sub commands: help/version`)
		return nil
	}, func() interface{} {
		return new(help_t)
	})

	app.Register(&cli.Command{
		// NOTE: Name is required, panic if ""
		Name: "version",
		// NOTE: Fn is required, panic if nil
		Fn: func(ctx *cli.Context) error {
			fmt.Println(`version: v0.0.1`)
			return nil
		},
		// NOTE: ArgvFn is required, panic if nil
		ArgvFn: func() interface{} {
			return new(version_t)
		},

		Desc: "Desc is optional",
	})

	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Printf("%v\n", err)
	}
}

type arg_t struct {
	Help bool   `cli:"h,help" usage:"show help"`
	Host string `cli:"H,host" usage:"specify host address" dft:"127.0.0.1"`
	Port uint16 `cli:"p,port" usage:"specify http port" dft:"8080"`
}

type help_t struct {
}

type version_t struct {
}
```
