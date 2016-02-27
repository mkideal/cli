package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	errEmptyCommand    = errors.New("empty command")
	errCommandNotFound = errors.New("command not found")
)

type (
	Context struct {
		router     []string
		path       string
		argv       interface{}
		nativeArgs []string
		flagSet    *FlagSet
		command    *Command
	}

	CommandFunc func(*Context) error

	ArgvFunc func() interface{}

	Command struct {
		Name   string
		Desc   string
		Fn     CommandFunc
		ArgvFn ArgvFunc

		parent   *Command
		children []*Command

		writer io.Writer
	}
)

//---------
// Context
//---------
func newContext(router, args []string, argv interface{}) (*Context, error) {
	ctx := &Context{
		router: router,
		argv:   argv,
	}
	ctx.nativeArgs = args

	ctx.path = strings.Join(ctx.router, " ")
	ctx.flagSet = Parse(args, argv)
	if ctx.flagSet.Error != nil {
		return nil, ctx.flagSet.Error
	}

	return ctx, nil
}

func (ctx *Context) FlagSet() *FlagSet {
	return ctx.flagSet
}

func (ctx *Context) Path() string {
	return ctx.path
}

func (ctx *Context) Router() []string {
	return ctx.router
}

func (ctx *Context) Args() []string {
	return ctx.nativeArgs
}

func (ctx *Context) Argv() interface{} {
	return ctx.argv
}

func (ctx *Context) Command() *Command {
	return ctx.command
}

func (ctx *Context) Writer() io.Writer {
	if ctx.command.writer == nil {
		return os.Stdout
	}
	return ctx.command.writer
}

//---------
// Command
//---------
func (cmd *Command) Register(child *Command) *Command {
	if child.Name == "" {
		panic(`child.Name == ""`)
	}
	if child.ArgvFn == nil {
		panic(`child.ArgvFn == nil`)
	}
	if cmd.children == nil {
		cmd.children = []*Command{}
	}
	if cmd.findChild(child.Name) != nil {
		panic(fmt.Sprintf("repeat child `%s` of `%s`", child.Name, cmd.Name))
	}
	cmd.children = append(cmd.children, child)
	child.parent = cmd
	if child.writer == nil {
		child.writer = child.parent.writer
	}
	return child
}

func (cmd *Command) RegisterFunc(name string, fn CommandFunc, argvFn ArgvFunc) *Command {
	return cmd.Register(&Command{Name: name, Fn: fn, ArgvFn: argvFn})
}

func (cmd *Command) Run(args []string) error {
	router := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, dashOne) {
			break
		}
		router = append(router, arg)
	}
	if len(router) == 0 {
		if cmd.Fn == nil || cmd.ArgvFn == nil {
			return errEmptyCommand
		}
	}
	child := cmd.route(router)
	if child == nil {
		return errCommandNotFound
	}
	ctx, err := newContext(router, args[len(router):], child.ArgvFn())
	if err != nil {
		return err
	}
	ctx.command = child
	return child.Fn(ctx)
}

func (cmd *Command) Usage() string {
	buff := bytes.NewBufferString("")
	fmt.Fprintf(buff, "Usage of `%s':\n", cmd.Path())
	if cmd.Desc != "" {
		fmt.Fprintf(buff, "\n%s\n\n", cmd.Desc)
	}
	fmt.Fprintf(buff, Usage(cmd.ArgvFn()))
	return buff.String()
}

func (cmd *Command) Path() string {
	path := cmd.Name
	cur := cmd
	for cur.parent != nil {
		cur = cur.parent
		path = cur.Name + " " + path
	}
	return path
}

func (cmd *Command) Parent() *Command {
	return cmd.parent
}

func (cmd *Command) Children() []*Command {
	return cmd.children
}

func (cmd *Command) route(router []string) *Command {
	cur := cmd
	for _, name := range router {
		if child := cur.findChild(name); child == nil {
			return nil
		} else {
			cur = child
		}
	}
	return cur
}

func (cmd *Command) findChild(name string) *Command {
	for _, child := range cmd.children {
		if child.Name == name {
			return child
		}
	}
	return nil
}
