# Command line interface [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/cli/master/LICENSE)

## License

[The MIT License (MIT)](https://raw.githubusercontent.com/mkideal/cli/master/LICENSE)

## Install
```sh
go get github.com/mkideal/cli
```

## Features

* Simplest, fast to learn how to use.
* Safety. Support type check, range check, and custom validate function.
* Based on golang tag. Support three tags: `cli`,`usage`,`dft`.
* Support specify default value and required declaration.
* Support multiple flag name for same argument.
* Support command tree.

## TODOs
* Add HTTP router
* Support `[]string` as argument struct's field
* Support argument struct inherit
* Support command suggestion
* Support command completion

## Getting started

### Just run it!

```go
package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	Help bool   `cli:"h,help" usage:"display help information"`
	Name string `cli:"name" usage:"your name" dft:"world"`
	Age  uint8  `cli:"a,age" usage:"your age" dft:"100"`
}

func main() {
	cli.Run(&argT{}, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		if argv.Help {
			ctx.String(ctx.Usage())
		} else {
			ctx.String("Hello, %s! Your age is %d?\n", argv.Name, argv.Age)
		}
		return nil
	})
}
```

Type these in terminal
```sh
$> go build -o hello
$> ./hello
Hello, world! Your age is 100?
$> ./hello --name=clipher -a 9
Hello, clipher! Your age is 9?
$> ./hello -h
```

Try
```sh
$> ./hello -a 256
```

## Tags

### cli

`cli` tag supports short format and long format, e.g.

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

`usage` tag describe the argument. If the argument is required, description has prefix `*` while show usage(`*` is red on unix-like os).

### dft
`dft` tag specify argument's default value.

## Command

### Command struct
```go
type Command struct {
	Name   string
	Desc   string
	Text   string
	Fn     CommandFunc
	Argv   ArgvFunc

	parent   *Command
	children []*Command

	writer io.Writer
}
```

### Command tree

Command can register child command using `Register` method or `RegisterFunc` method.

```go
func (cmd *Command) Register(*Command) *Command
func (cmd *Command) RegisterFunc(string, CommandFunc, ArgvFunc) *Command
```

```sh
root
├── sub1
│   ├── sub11
│   └── sub12
└── sub2
```
	
```go
var root = &cli.Command{}

var sub1 = root.Register(&cli.Command{
	Name: "sub1",
	Fn: func(ctx *cli.Context) error {
		//do something
	},
})
var sub11 = sub1.Register(&cli.Command{
	Name: "sub11",
	Fn: func(ctx *cli.Context) error {
		//do something
	},
})
var sub12 = sub1.Register(&cli.Command{
	Name: "sub12",
	Fn: func(ctx *cli.Context) error {
		//do something
	},
})

var sub2 = root.Register(&cli.Command{
	Name: "sub2",
	Fn: func(ctx *cli.Context) error {
		//do something
	},
})
```

## Examples

* [Hello](https://github.com/mkideal/cli/blob/master/examples/hello/main.go)
* [Basic](https://github.com/mkideal/cli/blob/master/examples/basic/main.go)
* [Simple Command](https://github.com/mkideal/cli/blob/master/examples/simple-command/main.go)
* [Multi Command](https://github.com/mkideal/cli/blob/master/examples/multi-command/main.go)
