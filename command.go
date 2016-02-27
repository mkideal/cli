package cli

import (
	"errors"
	"fmt"
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
	}

	CommandFunc func(*Context) error

	ArgvFunc func() interface{}

	Command struct {
		Name     string
		Desc     string
		Fn       CommandFunc
		ArgvFn   ArgvFunc
		Children []*Command
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

	if len(ctx.router) == 0 {
		return nil, errEmptyCommand
	}
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
	if cmd.Children == nil {
		cmd.Children = []*Command{}
	}
	if cmd.findChild(child.Name) != nil {
		panic(fmt.Sprintf("repeat child `%s` of `%s`", child.Name, cmd.Name))
	}
	cmd.Children = append(cmd.Children, child)
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
		return errEmptyCommand
	}
	child := cmd.route(router)
	if child == nil {
		return errCommandNotFound
	}
	ctx, err := newContext(router, args[len(router):], child.ArgvFn())
	if err != nil {
		return err
	}
	return child.Fn(ctx)
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
	for _, sub := range cmd.Children {
		if sub.Name == name {
			return sub
		}
	}
	return nil
}
