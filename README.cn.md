# Command line interface [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mkideal/cli/master/LICENSE)

## 许可协议

[MIT 许可协议](https://raw.githubusercontent.com/mkideal/cli/master/LICENSE)

## 安装获取
```sh
go get github.com/mkideal/cli
```

## 惊鸿一瞥

![screenshot](http://www.mkideal.com/images/screenshot.png)
![screenshot2](http://www.mkideal.com/images/screenshot2.png)

## 有些什么特色呢

* 简单,容易上手
* 有安全检查. 包括类型检查,值范围检查,更强大的是自定义的验证函数.
* 基于go语言的tag属性实现,参数定义结构化,简洁优雅.支持4种tag: `cli`,`usage`,`dft`, `name`.
* 支持默认值标签`dft`,可以用环境变量做默认值,支持required声明.
* 支持单个flag多个名字,像 -h --help 这样的.
* 支持命令树形结构，N层子命令随意玩.
* 支持子命令错误纠正提示,`hlp`会问你是要`help`吗
* 天然的命令树形结构摇身一变就可以变成HTTP路由了,像`$app hello world` -> `/hello/world`
* 支持命令行参数结构体的继承
* 支持短flag的组合式.`-x -y -z` -> `-xyz`, 不过必须全是bool型的才可以组合
* 支持长这样的`-Fvalue`的用法,它就等于`-F value`,不过`-F`不能是bool型
* 可以用 `--` 来隔离flags和arguments
* 支持使用数组和map了,数组的这样用:`-Fv1 -Fv2`(或`-F v1 -F v2`),map的`-Fk1=v1 -Fk2=v2`(或`-F k1=v1 -F k2=v2`)
* 使用帮助支持两种显示风格:默认的`左标签右说明`,manual风格的`上标签下说明`

## 还想做点什么
* 要有命令补全就好了
* 想把`dft`标签做强大点,能够支持表达式
* 还想在参数解析前/后,命令执行前加钩子
* ...... 还有什么呢

## 快速入门

### 迅速运行下面的代码!

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

在终端上,去到你的代码目录,然后敲入下面的命令:
```sh
$> go build -o hello
$> ./hello
Hello, world! Your age is 100?
$> ./hello --name=clipher -a 9
Hello, clipher! Your age is 9?
$> ./hello -h
```

不妨再试一试下面这个,看看有什么不对
```sh
$> ./hello -a 256
```

## clier - 命令生成其

`clier` 是一个用`cli`写成命令生成工具,你可以使用它快速创建新的命令.

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

## 标签

这个命令行程序工具库的重要特征之一就是使用tag来实现命令行程序的诸多特性.
一次说明支持的4种tag的用法.

### cli

`cli` 这个tag主要用来定义命令行参数的flag,同时还可以定义flag是否是required.

`cli` 支持短格式和长格式, 所谓短格式就是单字母前面只加一个横杠(`-h`),长格式是多字母,前面加两个横杠(`--help`).

```go
Help    bool    `cli:"h,help"`
Version string  `cli:"version"`
Port    int     `cli:"p"`
XYZ     bool    `cli:"x,y,z,xyz,XYZ"` 
```

给`cli` 加上一个前缀 `*` 就能将此flag声明为必须的(默认所有参数都是可选的), 比如:

```go
Required string `cli:"*required"`
```

还有一种前缀: `!`.加上这个前缀的flag必须是bool型,并且如果这个flag被赋值为true的话,将阻止参数必要性检查和自定义验证函数

```go
Help bool `cli:"!h,help"`
```

### usage

`usage`标签用于描述flag的作用/用法.如果flag被标记为required,那么usage显示文本会自动加上红色的`*`

NOTE: `cli.SetUsageStyle`函数用于设置帮助的显示风格.比如要设置成manual的:

```go
...
func main() {
	cli.SetUsageStyle(cli.ManualStyle)
	...
}
```

### dft

`dft`标签用来指定参数默认值,支持使用环境变量做默认值,比如下面这个

```go
Port   int    `cli:"p,port" usage:"http port" dft:"8080"` 
Home   string `cli:"home" usage:"home dir" dft:"$HOME"`
GoRoot string `cli:"goroot" usage:"go root dir" dft:"$GOROOT"`
```

### name
`name` 标签也许是用武之地最小的tag,用来给flag定义一个在帮助中的引用名.

## Command 对象

`Command`是`cli`的两大重要对象之一(另一个是`Context`),它定义了一条命令的信息和执行函数.

### Command 的定义

```go
type Command struct {
	Name        string		// 命令的名字
	Desc        string		// 命令的简要描述
	Text        string		// 命令的详细描述
	Fn          CommandFunc // 命令的执行函数
	Argv        ArgvFunc	// 命令接受的参数对象的工厂函数

	CanSubRoute bool		// 是否可以部分路由匹配,什么意思呢?举个栗子:
							//  假设编译了一个程序叫app
							//  它有一个命令叫hello,而hello没有world
							//  子命令,那么如果hello命令的CanSubRoute
							//	为false(这是默认行为),下面的执行
							//  就会提示你找不到`hello world`命令
							//		$> ./app hello world
							//	而如果CanSubRoute为true的话,上面
							//  的程序就会执行hello命令,而world是hello
							//  命令执行时接受到的第一个参数
							//	如果hello有world子命令,那么不管
							//  CanSubRoute为true还是false,都会执行
							//  world子命令

	HTTPRouters []string	// 为命令定义额外的http路由.这里说额外,
							// 是因为每个命令都有一个标准的路由,是根据
							// 命令在树形结构中的位置生成的一个路由
							// 上面`./app hello world`的栗子里,
							// hello的标准路由就是`/hello`
							// world若是hello的子命令,那么
							// world的标准路由就是`/hello/world`
							// 加入给hello命令的HTTPRouters加一个
							// `/v1/hello`路由,那么http访问`/v1/hello`
							// 就如同访问`/hello`

	HTTPMethods []string	// 为命令指定接受的HTTP方法.如果HTTPMethods
							// 为空,那么命令接受任何HTTP方法
	...
}
```

### Command 树结构

#### 构建Command 树的方法一(推荐): 使用Root和Tree函数构建

```go
// Tree 为cmd注册若干颗子树,返回以cmd为根的一颗命令树
// - 它会创建一个新的CommandTree对象
func Tree(cmd *Command, forest ...*CommandTree) *CommandTree

// Root 为root注册若干颗子树,返回root
// - 它不会创建一个新的CommandTree对象
func Root(root *Command, forest ...*CommandTree) *Command
```

下面是一个具有多层深度的命令树

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

#### 构建Command 树的方法一: 手动给每条命令注册子命令

`Command`有个两个方法用来注册子命令: `Register` 和`RegisterFunc`

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

下面的代码按上面这个树形结构构建命令树
	
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

## HTTP 路由

`Command` 对象实现了 `ServeHTTP` 方法.

```go
func (cmd *Command) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

`Command` 两个用于HTTP的属性: `HTTPMethods` 和 `HTTPRouters`(参见`Command`结构定义的说明)

```go
HTTPRouters []string
HTTPMethods []string
```

一个设置HTTP属性的栗子:

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

**NOTE**: 如果要是用自定义的额外HTTP路由必须调用根命令的`RegisterHTTP`方法.

```go
...
if err := ctx.Command().Root().RegisterHTTP(ctx); err != nil {
	return err
}
return ctx.Command().Root().ListenAndServeHTTP(addr)
...
```
想了解更多,可查看示例 [HTTP](./examples/http/main.go).

## 远程HTTP调用(RPC)

`Command` 有一个RPC方法实现了基于HTTP的远程调用.

想要了解如何使用的话,可以查看示例 [RPC](./examples/rpc/main.go).

## 更多栗子

* [Hello](./examples/hello/main.go)
* [Screenshot](./examples/screenshot/main.go)
* [Basic](./examples/basic/main.go)
* [Simple Command](./examples/simple-command/main.go)
* [Multi Command](./examples/multi-command)
* [Tree](./examples/tree/main.go)
* [HTTP](./examples/http/main.go)
* [RPC](./examples/rpc/main.go)

## 辅助工具

* [goplus](https://github.com/mkideal/goplus) - `可用来生成cli型的程序支架`
