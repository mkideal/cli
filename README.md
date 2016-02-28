# Command line interface [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/cli/master/LICENSE)

## License
[The MIT License (MIT)](https://en.wikipedia.org/wiki/MIT_License) Enlish wiki
[The MIT License (MIT)](https://zh.wikipedia.org/wiki/MIT許可證)  中文wiki

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

## Getting started
```go
package main

import (
	"github.com/mkideal/cli"
)

type arg_t struct {
	Help bool   `cli:"h,help" usage:"display help information"`
	Name string `cli:"name" usage:"your name" dft:"world"`
	Age  uint8  `cli:"a,age" usage:"your age" dft:"100"`
}

func main() {
	cli.Run(&arg_t{}, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*arg_t)
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
$> ./hello --name gopher
Hello, gopher! Your age is 100?
$> ./hello --name=clipher -a 9
Hello, clipher! Your age is 9?
$> ./hello -h
```

Try
```sh
./hello -a 256
```

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
	Float64        float64 `cli:"f64" usage:"type float64"`
}
```

Then, call function `cli.Parse`:
```go
t := new(Args)
flagSet := cli.Parse(os.Args[1:], t)
if flagSet.Error != nil {
	//TODO: handle the error
}
//^Try uncomment following line
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

`usage` tag describe the argument. If the argument is required, description string has prefix `*` while show usage(`*` is red on unix-like os).

### dft
`dft` tag specify argument's default value.

## Command

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

Examples:
[Hello](https://github.com/mkideal/cli/blob/master/examples/hello/main.go)
[Basic](https://github.com/mkideal/cli/blob/master/examples/basic/main.go)
[Simple Command](https://github.com/mkideal/cli/blob/master/examples/simple-command/main.go)
[Multi Command](https://github.com/mkideal/cli/blob/master/examples/multi-command/main.go)
