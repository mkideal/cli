# Command line interface [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/cli/master/LICENSE)

## Screenshot

![screenshot2](http://www.mkideal.com/images/screenshot2.png)

## Getting started

### Just run it!

```go
package main

import (
	"github.com/mkideal/cli"
)

type argT struct {
	Help bool   `cli:"!h,help" usage:"display help information"`
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
Port    int     `cli:"p"`       // -p
Version string  `cli:"version"` // --version
Help    bool    `cli:"h,help"`  // -h OR --help
XYZ     bool    `cli:"x,y,xy"`  // -x OR -y OR --xy
```

The argument is required if ``cli`` tag has prefix ``*``, e.g.

```go
Required string `cli:"*r"`
```

The argument is marked as `force` flag if `cli` tag has prefix `!`, e.g

```go
Help bool `cli:"!h,help"`
```

Will prevent validating arguments if `force` flag assigned with true.

### usage

`usage` tag describes the flag.

### dft

`dft` tag specifies argument's default value. You can specify ENV as default value. e.g.

```go
Port   int    `cli:"p,port" usage:"http port" dft:"8080"` 
Home   string `cli:"home" usage:"home dir" dft:"$HOME"`
GoRoot string `cli:"goroot" usage:"go root dir" dft:"$GOROOT"`
```

### name

`name` tag give a name for show.

### pw

`pw` tag used for typing password. Enter the password in prompt, e.g.

```go
Password string `pw:"p,password" usage:"type the password" prompt:"type the password"`
```

```sh
$> ./app -p
type the password:
```

### prompt

`prompt` is the prompt string.


## Supported types for flag

* All basic types: int,uint,...,flaot32,float64,string,bool
* Slice of basic type: []int, []uint, []string,...
* Map of basic type: map[uint]int, map[string]string,...

```go
type argT struct {
	Bool    bool    `cli:"b,bool" usage:"-b OR --bool OR -b=true OR -b false OR --bool false OR --bool=false"`
	Int     int     `cli:"i,int" usage:"-i6 OR -i=-7 OR -i 8 OR --int=9 OR --int 9"`
	Uint    uint    `cli:"u,uint" usage:"-u1 OR -u=2 OR -u 3 OR --uint=4 OR --uint 5"`
	Int8    int8    `cli:"i8" usage:"int8 type"`
	Uint8   uint8   `cli:"u8" usage:"uint8 type"`
	Int16   int16   `cli:"i16" usage:"int16 type"`
	Uint16  uint16  `cli:"u16" usage:"uint16 type"`
	Int32   int32   `cli:"i32" usage:"int32 type"`
	Uint32  uint32  `cli:"u32" usage:"uint32 type"`
	Int64   int64   `cli:"i64" usage:"int64 type"`
	Uint64  uint64  `cli:"u64" usage:"uint64 type"`
	Float32 float32 `cli:"f32" usage:"float32 type"`
	Float64 float64 `cli:"f64" usage:"float64 type"`
	String  string  `cli:"s,string" usage:"string type"`

	// fold flags for bool
	BoolX bool `cli:"x" usage:"bool x"`
	BoolY bool `cli:"y" usage:"bool y"`
	BoolZ bool `cli:"z" usage:"bool z"`
	// You can use it like this:
	// -x OR -xy OR -xyz

	Slice []uint32          `cli:"S,slice" usage:"-S 6 -S 7 -S8 --slice 9 (results: [6 7 8 9])"`
	Map   map[string]uint16 `cli:"M,map" usage:"-Mx=1 -M y=2 --map z=3 (results:[x:1 y:2 z:3])"`
}
```

## Command

```go
type Command struct {
	Name        string		// name
	Aliases     []string    // aliases names
	Desc        string		// description
	Text        string		// detail description
	Fn          CommandFunc // handler
	Argv        ArgvFunc	// argument factory function

	NoHook      bool
	CanSubRoute bool
	HTTPRouters []string
	HTTPMethods []string

	// hooks for current command
	OnBefore func(*Context) error
	OnAfter  func(*Context) error

	// hooks for all commands if current command is root command
	OnRootPrepareError func(error) error
	OnRootBefore       func(*Context) error
	OnRootAfter        func(*Context) error
	...
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

## Examples

* [Hello](./examples/hello/main.go)
* [Screenshot](./examples/screenshot/main.go)
* [Basic](./examples/basic/main.go)
* [Simple Command](./examples/simple-command/main.go)
* [Multi Command](./examples/multi-command)
* [Tree](./examples/tree/main.go)
* [HTTP](./examples/http/main.go)
* [RPC](./examples/rpc/main.go)
* [Daemon](./examples/daemon/main.go)

## Projects
* [onepw](https://github.com/mkideal/onepw) - A lightweight tool for managing passwords
* [rqlite CLI](https://github.com/rqlite/rqlite/tree/master/cmd/rqlite) - A command line tool for connecting to a rqlited node
