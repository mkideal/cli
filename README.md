# Command line interface [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/cli/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/mkideal/cli)](https://goreportcard.com/report/github.com/mkideal/cli) [![GoDoc](https://godoc.org/github.com/mkideal/cli?status.svg)](https://godoc.org/github.com/mkideal/cli)

## Screenshot

![screenshot2](http://www.mkideal.com/images/screenshot2.png)

## Key features

* Lightweight and easy to use.
* Defines flag by tag, e.g. flag name(short or/and long), description, default value, password, prompt and so on.
* Type safety.
* Output looks very nice.
* Supports custom Validator.
* Supports slice and map as a flag.
* Supports any type as a flag field which implements cli.Decoder interface.
* Supports any type as a flag field which use FlagParser.
* Suggestions for command.(e.g. `hl` => `help`, "veron" => "version").
* Supports default value for flag, even expression about env variable(e.g. `dft:"$HOME/dev"`).
* Supports editor like `git commit` command.(See example [21](http://www.mkideal.com/golang/cli-examples.html#example-21-editor) and [22](http://www.mkideal.com/golang/cli-examples.html#example-22-custom-editor))

## Getting started

```go
package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Name string `cli:"name" usage:"your name" dft:"world"`
	Age  uint8  `cli:"a,age" usage:"your age" dft:"100"`
}

func main() {
	cli.Run(&argT{}, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("Hello, %s! Your age is %d?\n", argv.Name, argv.Age)
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

## API documentation

See [**godoc**](https://godoc.org/github.com/mkideal/cli)

## Examples

See [**_examples**](https://github.com/mkideal/cli/tree/master/_examples), example_*test.go files or site below:

* [Example 1: Hello world](http://www.mkideal.com/golang/cli-examples.html#example-1-hello)
* [Example 2: How to use **flag**](http://www.mkideal.com/golang/cli-examples.html#example-2-flag)
* [Example 3: How to use **required** flag](http://www.mkideal.com/golang/cli-examples.html#example-3-required-flag)
* [Example 4: How to use **default** flag](http://www.mkideal.com/golang/cli-examples.html#example-4-default-flag)
* [Example 5: How to use **slice**](http://www.mkideal.com/golang/cli-examples.html#example-5-slice)
* [Example 6: How to use **map**](http://www.mkideal.com/golang/cli-examples.html#example-6-map)
* [Example 7: Usage of **force** flag](http://www.mkideal.com/golang/cli-examples.html#example-7-force-flag)
* [Example 8: Usage of **child command**](http://www.mkideal.com/golang/cli-examples.html#example-8-child-command)
* [Example 9: **Auto help**](http://www.mkideal.com/golang/cli-examples.html#example-9-auto-help)
* [Example 10: Usage of **Validator**](http://www.mkideal.com/golang/cli-examples.html#example-10-usage-of-validator)
* [Example 11: **Prompt** and **Password**](http://www.mkideal.com/golang/cli-examples.html#example-11-prompt-and-password)
* [Example 12: How to use **Decoder**](http://www.mkideal.com/golang/cli-examples.html#example-12-decoder)
* [Example 13: Builtin Decoder: **PidFile**](http://www.mkideal.com/golang/cli-examples.html#example-13-pid-file)
* [Example 14: Builtin Decoder: **Time** and **Duration**](http://www.mkideal.com/golang/cli-examples.html#example-14-time-and-duration)
* [Example 15: Builtin Decoder: **File**](http://www.mkideal.com/golang/cli-examples.html#example-15-file)
* [Example 16: **Parser**](http://www.mkideal.com/golang/cli-examples.html#example-16-parser)
* [Example 17: Builtin Parser: **JSONFileParser**](http://www.mkideal.com/golang/cli-examples.html#example-17-json-file)
* [Example 18: How to use **custom parser**](http://www.mkideal.com/golang/cli-examples.html#example-18-custom-parser)
* [Example 19: How to use **Hooks**](http://www.mkideal.com/golang/cli-examples.html#example-19-hooks)
* [Example 20: How to use **Daemon**](http://www.mkideal.com/golang/cli-examples.html#example-20-daemon)
* [Example 21: How to use **Editor**](http://www.mkideal.com/golang/cli-examples.html#example-21-editor)
* [Example 22: Custom **Editor**](http://www.mkideal.com/golang/cli-examples.html#example-22-custom-editor)

## Projects which use CLI

* [onepw](https://github.com/mkideal/onepw) - A lightweight tool for managing passwords
* [rqlite CLI](https://github.com/rqlite/rqlite/tree/master/cmd/rqlite) - A command line tool for connecting to a rqlited node

## Argument object of cli

Supported tags in argument object of cli:

* cli - supports short format and long format, e,g, `-p`, `--port`
* pw - similar to `cli`, but used to input password
* usage - description of flag
* dft - default value, supports constant value, env variable, and even expression
* name - placeholder for flag
* prompt - prompt string
* parser - builtin parsers: `json`, `jsonfile`, supports custom parsers.

### Supported types of flag

* All basic types: int,uint,...,flaot32,float64,string,bool
* Slice of basic type: []int, []uint, []string,...
* Map of basic type: map[uint]int, map[string]string,...
* Any type which implments `cli.Decoder`
* Any type which use correct parser(`json`,`jsonfile`, or your registered parser)

### AutoHelper

```go
// AutoHelper represents interface for showing help information automaticly
AutoHelper interface {
	AutoHelp() bool
}
```

If your `argT` implments AutoHelper, it will show help if argT.AutoHelp return true.

For example:

```go
package main

import "github.com/mkideal/cli"

type argT struct {
	Help bool `cli:"h,help" usage:"show help"`
}

func (argv *argT) AutoHelp() bool {
	return argv.Help
}

func main() {
	cli.Run(&argT{}, func(ctx *cli.Context) error {
		return nil
	})
}
```

Build and run:
```sh
go build -o app
./app -h
```

This will print help information.

For convenience, builtin type `cli.Helper` implements cli.AutoHelper.

So, you just need to inherit it!

```go
type argT struct {
	cli.Helper
}
```

### Validator - validate argument

```go
// Validator validates flag before running command
type Validator interface {
	Validate(*Context) error
}
```

### Context.IsSet
Determin whether the flag is set

### Hooks

* *OnRootPrepareError* - Function of root command which should be invoked if occur error while prepare period
* *OnBefore* - Function which should be invoked before `Fn`
* *OnRootBefore* - Function of root command which should be invoked before `Fn`
* *OnRootAfter* - Function of root command which should be invoked after `Fn`
* *OnAfter* - Function which should be invoked after `Fn`

Invoked sequence: OnRootPrepareError => OnBefore => OnRootBefore => Fn => OnRootAfter => OnAfter

Don't invoke any hook function if `NoHook` property of command is true, e.g.

```go
var helpCmd = &cli.Command{
	Name: "help",
	Desc: "balabala...",
	NoHook: true,
	Fn: func(ctx *cli.Context) error { ... },
}
```

### Command tree

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

Each command have one default router: slashes-separated router path of command. e.g.
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
