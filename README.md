# Command line interface [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/cli/master/LICENSE)

## License

[The MIT License (MIT)](https://raw.githubusercontent.com/mkideal/cli/master/LICENSE)

## [中文文档](./README.cn.md)

## Install
```sh
go get github.com/mkideal/cli
```

## Screenshot

![screenshot](http://www.mkideal.com/images/screenshot.png)
![screenshot2](http://www.mkideal.com/images/screenshot2.png)

## Features

* Simplest, fast to learn.
* Safety. Support type check, range check, and custom validate function.
* Based on golang tag. Support four tags: `cli`,`usage`,`dft`, `name`.
* Support default value and required declaration.
* Support multiple flag name for same argument.
* Support command tree.
* Support command suggestion
* Support HTTP router
* Support struct field
* Support `-F<value>` flag format
* Support separated flags and arguments by `--`
* Distinguish flags and arguments - `app cmd --flag -b=c arg1 arg2` 
* Support array flag - `-F v1 -F v2` or `-Fv1 -Fv2`
* Support map flag - `-F k1=v1 -F k2=v2` or `-F<k1=v1> -F<k2=v2>`

## TODOs
* Support command completion
* Support Before/After hooks
* Support expr for `dft` tag

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

## clier - command generator

`clier` is command generator write by `cli`. You can use it for fast-creating a new command.

```sh
Usage: clier [OPTIONS] COMMAND-NAME

Examples:
	clier hello
	clier -f -s "balabalabala" hello
	clier -p balabala hello

Options:

  -h, --help
      display help

  -F, --file=NAME
      create source file for new command, default <commandName>.go

  -f, --force[=false]
      force create file if exists

  -p, --package
      dest package name, default <basedir FILE>

  -s, --desc
      command description

  --csr, --can-sub-route[=false]
      set CanSubRoute attribute for new command

  --argv-type-name
      argv type, default <commandName>T, e.g. command name is hello, then defaut argv type is helloT
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

The argument is marked as `force` flag if `cli` tag has prefix `!`, e.g

```go
Help bool `cli:"!h,help"`
```

Don't validate flags if `force` assigned with true.

### usage

`usage` tag describes the argument. If the argument is required, description has prefix `*` while show usage(`*` is red on unix-like os).

### dft
`dft` tag specifies argument's default value.You can specify ENV as default value. e.g.

```go
Port   int    `cli:"p,port" usage:"http port" dft:"8080"` 
Home   string `cli:"home" usage:"home dir" dft:"$HOME"`
GoRoot string `cli:"goroot" usage:"go root dir" dft:"$GOROOT"`
```

### name
`name` tag give a name for show.

## Command

### Command struct

```go
type Command struct {
	Name        string
	Desc        string
	Text        string
	Fn          CommandFunc
	Argv        ArgvFunc
	CanSubRoute bool
	HTTPRouters []string
	HTTPMethods []string
	...
}
```

### Command tree

#### Method 1: Construct command tree

```go
// Tree creates a CommandTree
func Tree(cmd *Command, forest ...*CommandTree) *CommandTree

// Root registers forest for root and return root
func Root(root *Command, forest ...*CommandTree) *Command
```

```go
func log(ctx *cli.Context) error {
	ctx.String("path: `%s`\n", ctx.Path())
	return nil
}
var (
	cmd1  = &cli.Command{Name: "cmd1", Fn: log}
	cmd11 = &cli.Command{Name: "cmd11", Fn: log}
	cmd12 = &cli.Command{Name: "cmd12", Fn: log}

	cmd2   = &cli.Command{Name: "cmd2", Fn: log}
	cmd21  = &cli.Command{Name: "cmd21", Fn: log}
	cmd22  = &cli.Command{Name: "cmd22", Fn: log}
	cmd221 = &cli.Command{Name: "cmd221", Fn: log}
	cmd222 = &cli.Command{Name: "cmd222", Fn: log}
	cmd223 = &cli.Command{Name: "cmd223", Fn: log}
)

root := cli.Root(
	&cli.Command{Name: "tree"},
	cli.Tree(cmd1,
		cli.Tree(cmd11),
		cli.Tree(cmd12),
	),
	cli.Tree(cmd2,
		cli.Tree(cmd21),
		cli.Tree(cmd22,
			cli.Tree(cmd221),
			cli.Tree(cmd222),
			cli.Tree(cmd223),
		),
	),
)
```

#### Method 2: Register child command
Command registers child command using `Register` method or `RegisterFunc` method.

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

## HTTP router

Context implements ServeHTTP method.

```go
func (ctx *Context) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

`Command` has two http properties: `HTTPMethods` and `HTTPRouters`.

```go
HTTPRouters []string
HTTPMethods []string
```

Each command have one default router: slashes-separated router path of command. e.g. command
`$app sub1 sub11` has default router `/sub1/sub11`. The `HTTPRouters` property will add some new extra routers, it would'nt replace the default router. The `HTTPMethods` is a string array. It will handle any method if `HTTPMethods` is empty.

```go
var help = &cli.Command{
	Name:        "help",
	Desc:        "display help",
	CanSubRoute: true,
	HTTPRouters: []string{"/v1/help", "/v2/help"},
	HTTPMethods: []string{"GET", "POST"},

	Fn: cli.HelpCommandFn,
}
```

**NOTE**: Must call root command's RegisterHTTP method if you want to use custom `HTTPRouters`

```go
...
if err := ctx.Command().Root().RegisterHTTP(ctx); err != nil {
	return err
}
return http.ListenAndServe(addr, ctx.Command().Root())
...
```
See example [HTTP](./examples/http/main.go).

## RPC

See example [RPC](./examples/rpc/main.go).

## Examples

* [Hello](./examples/hello/main.go)
* [Screenshot](./examples/screenshot/main.go)
* [Basic](./examples/basic/main.go)
* [Simple Command](./examples/simple-command/main.go)
* [Multi Command](./examples/multi-command)
* [Tree](./examples/tree/main.go)
* [HTTP](./examples/http/main.go)
* [RPC](./examples/rpc/main.go)

## Who use

* [goplus](https://github.com/mkideal/goplus)

## External tools

* [goplus](https://github.com/mkideal/goplus) - `generate go application skeleton`

