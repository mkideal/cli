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

## Getting started

### Just run it!

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

## Argument object of cli

Supported tags in argument object of cli:

* cli
* pw
* usage
* dft
* name
* prompt
* parser

### cli

`cli` tag supports short format and long format, e.g.

```go
Port    int     `cli:"p"`       // -p
Version string  `cli:"version"` // --version
Help    bool    `cli:"h,help"`  // -h OR --help
XYZ     bool    `cli:"x,y,xy"`  // -x OR -y OR --xy

Required string `cli:"*r"`
Help bool `cli:"!h,help"`
```

The argument is required if ``cli`` tag has prefix ``*``, e.g.
The argument is marked as a `force` flag if `cli` tag has prefix `!`.
Will prevent validating arguments if `force` flag assigned with true.

How to use flag?
Assume, argT has flags like this:

```go
type argT {
	Flag  string         `cli:"f,flag"`
	Slice []string       `cli:"D,slice"`
	Map   map[string]int `cli:"M,map"`
}
```

Now, you can use it:
```sh
-f value
-f=value
-fvalue ### NOTE: in this case, -f MUST not be a boolean
--flag value
--flag=value

-D1 -D2
-D 1 -D 2
-D1 -D2
-D1 --slice 2
--slice 1 --slice 2

-Mx=1 -My=2
-Mx=1 -M y=2
-M x=1 -M y=2
-Mx=1 --map y=2
--map x=1 --map y=2
```

### pw

`pw` tag similar to `cli` but used for typing password. You can type the password in prompt, e.g.

### usage

`usage` tag describes the flag.

### dft

`dft` tag specifies argument's default value. You can use a env variable as default value, even expression. e.g.

```go
Port   int    `cli:"p,port" usage:"http port" dft:"8080"` 
GoRoot string `cli:"goroot" usage:"go root dir" dft:"$GOROOT"`
Home   string `cli:"home" usage:"home dir" dft:"$HOME"`
Devdir string `cli:"dev" usage:"dev dir" dft:"$HOME/dev"`
Port   int    `cli:"p,port" usage:"listening port" dft:"$HTTP_PORT+1000"`
```

### name

`name` tag give a reference name.

```go
Password string `pw:"p,password" usage:"type the password" prompt:"type the password"`
```

```sh
$> ./app -p
type the password:
```

### prompt

`prompt` is the prompt string.

### parser

If `parser` is set, will using specific parser parses flag. Here are some Builtin parsers:

* json
* jsonfile

You can implements your parser, and register it.

```go
// FlagParser interface
type FlagParser interface {
	Parse(s string) error
}

// FlagParserCreator
type FlagParserCreator func(interface{}) FlagParser

// Register parser creator
func RegisterFlagParser(name string, creator FlagParserCreator)
```

See example [JSON](./examples/json/main.go)

### Supported types of flag

* All basic types: int,uint,...,flaot32,float64,string,bool
* Slice of basic type: []int, []uint, []string,...
* Map of basic type: map[uint]int, map[string]string,...
* Any type which implments `cli.Decoder`
* Any type which use correct parser

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

	// custom type
	JSON jsonT `cli:"json" usage:"custom type which implements Decoder"`
	Config Config `cli:"config" usage:"custom type which use parser" parser:"jsonfile"`
}

type jsonT struct {
	K1 string
	K2 int
}

func (j *jsonT) Decode(s string) error {
	return json.Unmarshal([]byte(s), j)
}

type Config struct {
	//...
}
```

Some predefine types:

* [cli.Helper](https://github.com/mkideal/cli/blob/master/builtin.go#L8)
* [cli/ext.Time](https://github.com/mkideal/cli/blob/master/ext/types.go#L15)
* [cli/ext.Duration](https://github.com/mkideal/cli/blob/master/ext/types.go#L69)
* [cli/ext.File](https://github.com/mkideal/cli/blob/master/ext/types.go#L87)
* [cli/ext.PidFile](https://github.com/mkideal/cli/blob/master/ext/types.go#L126)

**NOTE**: `Parser` vs `Decoder`

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

### Validator - custom validating argument

```go
// Validator validates flag before running command
type Validator interface {
	Validate(*Context) error
}
```

See example [Basic](./examples/hello/main.go).

### Context.IsSet
Determin wether the flag is set

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
See example [HTTP](./examples/http/main.go).

## Context

```go
type Context struct {
	...
}
```

Methods of Context:
```go
// Path of "app cmd sub -f arg1 arg2" is "cmd/sub"
func (ctx *Context) Path() string

// Router of "app cmd sub -f arg1 arg2" is ["cmd" "sub"]
func (ctx *Context) Router() []string

// NativeArgs of "app cmd sub -f arg1 arg2" is ["-f" "arg1" "arg2"]
func (ctx *Context) NativeArgs() []string

// Args of "app cmd sub -f arg1 arg2" is ["arg1" "arg2"]
func (ctx *Context) Args() []string

// Argument object
func (ctx *Context) Argv() interface{}

// Encode flags to url.Values
func (ctx *Context) FormValues() url.Values

func (ctx *Context) Command() *Command
func (ctx *Context) Usage() string
func (ctx *Context) WriteUsage()
func (ctx *Context) Writer() io.Writer
func (ctx *Context) Write(data []byte) (n int, err error)
func (ctx *Context) Color() *color.Color

func (ctx *Context) String(format string, args ...interface{}) *Context
func (ctx *Context) JSON(obj interface{}) *Context
func (ctx *Context) JSONln(obj interface{}) *Context
func (ctx *Context) JSONIndent(obj interface{}, prefix, indent string) *Context
func (ctx *Context) JSONIndentln(obj interface{}, prefix, indent string) *Context
```

## Examples

* [Hello](./examples/hello/main.go)
* [Screenshot](./examples/screenshot/main.go)
* [Basic](./examples/basic/main.go)
* [JSON](./examples/json/main.go)
* [Simple Command](./examples/simple-command/main.go)
* [Multi Command](./examples/multi-command)
* [Tree](./examples/tree/main.go)
* [HTTP](./examples/http/main.go)
* [RPC](./examples/rpc/main.go)
* [Daemon](./examples/daemon/main.go)
* [JSON](./examples/json/main.go)

## Projects
* [onepw](https://github.com/mkideal/onepw) - A lightweight tool for managing passwords
* [rqlite CLI](https://github.com/rqlite/rqlite/tree/master/cmd/rqlite) - A command line tool for connecting to a rqlited node
