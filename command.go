package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

var (
	errEmptyCommand = errors.New("empty command")
)

type (
	// Context provider running context
	Context struct {
		router     []string
		path       string
		argv       interface{}
		nativeArgs []string
		flagSet    *flagSet
		command    *Command
	}

	// Validator validate flag before run command
	Validator interface {
		Validate() error
	}

	// CommandFunc aliases command handle function
	CommandFunc func(*Context) error

	// ArgvFunc aliases command argv factory function
	ArgvFunc func() interface{}

	// Command is the main object in command-line app
	Command struct {
		Name string      // Command name
		Desc string      // Command abstract
		Text string      // Command detailed description
		Fn   CommandFunc // Command handler
		Argv ArgvFunc    // Command argument factory function

		parent   *Command
		children []*Command

		writer io.Writer
		usage  string
	}
)

//---------
// Context
//---------

func newContext(path string, router, args []string, argv interface{}) (*Context, error) {
	ctx := &Context{
		path:       path,
		router:     router,
		argv:       argv,
		nativeArgs: args,
	}
	if argv != nil {
		ctx.flagSet = parseArgv(args, argv)
		if ctx.flagSet.err != nil {
			return nil, ctx.flagSet.err
		}
	}
	return ctx, nil
}

// Path returns full command name
// `./app hello world -a --xyz=1` will returns "hello world"
func (ctx *Context) Path() string {
	return ctx.path
}

// Router returns full command name with string array
// `./app hello world -a --xyz=1` will returns ["hello" "world"]
func (ctx *Context) Router() []string {
	return ctx.router
}

// Args returns native args
// `./app hello world -a --xyz=1` will returns ["-a" "--xyz=1"]
func (ctx *Context) Args() []string {
	return ctx.nativeArgs
}

// Argv returns parsed args object
func (ctx *Context) Argv() interface{} {
	return ctx.argv
}

// FormValues returns parsed args as url.Values
func (ctx *Context) FormValues() url.Values {
	return ctx.flagSet.values
}

// Command returns current command object
func (ctx *Context) Command() *Command {
	return ctx.command
}

// Usage returns current command's usage
func (ctx *Context) Usage() string {
	return ctx.command.Usage()
}

// Writer returns current command's writer
func (ctx *Context) Writer() io.Writer {
	return ctx.command.Writer()
}

// String writes format string to writer
func (ctx *Context) String(format string, args ...interface{}) *Context {
	fmt.Fprintf(ctx.Writer(), format, args...)
	return ctx
}

// JSON writes json string of obj to writer
func (ctx *Context) JSON(obj interface{}) *Context {
	data, err := json.Marshal(obj)
	if err == nil {
		fmt.Fprintf(ctx.Writer(), string(data))
	}
	return ctx
}

// JSONln writes json string of obj end with "\n" to writer
func (ctx *Context) JSONln(obj interface{}) *Context {
	return ctx.JSON(obj).String("\n")
}

// JSONIndent writes pretty json string of obj to writer
func (ctx *Context) JSONIndent(obj interface{}, prefix, indent string) *Context {
	data, err := json.MarshalIndent(obj, prefix, indent)
	if err == nil {
		fmt.Fprintf(ctx.Writer(), string(data))
	}
	return ctx
}

// JSONIndentln writes pretty json string of obj end with "\n" to writer
func (ctx *Context) JSONIndentln(obj interface{}, prefix, indent string) *Context {
	return ctx.JSONIndent(obj, prefix, indent).String("\n")
}

//---------
// Command
//---------

// Writer sets default writer(os.Stdout) if nil and returns writer
func (cmd *Command) Writer() io.Writer {
	if cmd.writer == nil {
		cmd.writer = os.Stdout
	}
	return cmd.writer
}

// SetWriter sets sepcify writer
func (cmd *Command) SetWriter(writer io.Writer) {
	cmd.writer = writer
}

// Register registers a child command
func (cmd *Command) Register(child *Command) *Command {
	if child.Name == "" {
		panic(`child.Name == ""`)
	}
	if cmd.children == nil {
		cmd.children = []*Command{}
	}
	if cmd.findChild(child.Name) != nil {
		panic(fmt.Sprintf("repeat register child `%s` for command `%s`", child.Name, cmd.Name))
	}
	cmd.children = append(cmd.children, child)
	child.parent = cmd

	// inherit parent's writer if nil
	if child.writer == nil {
		child.writer = child.parent.writer
	}
	return child
}

// RegisterFunc registers handler as child command
func (cmd *Command) RegisterFunc(name string, fn CommandFunc, argvFn ArgvFunc) *Command {
	return cmd.Register(&Command{Name: name, Fn: fn, Argv: argvFn})
}

// Run runs the command with args
func (cmd Command) Run(args []string) error {
	router := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, dashOne) {
			break
		}
		router = append(router, arg)
	}
	if len(router) == 0 {
		if cmd.Fn == nil {
			return errEmptyCommand
		}
	}
	path := strings.Join(router, " ")
	child := cmd.route(router)
	if child == nil {
		return fmt.Errorf("command `%s` not found", path)
	}

	var argv interface{}
	if child.Argv != nil {
		argv = child.Argv()
	}
	ctx, err := newContext(path, router, args[len(router):], argv)
	if err != nil {
		return err
	}
	if argv != nil {
		// validate argv if argv implements Validator interface
		if validator, ok := argv.(Validator); ok {
			if err := validator.Validate(); err != nil {
				return err
			}
		}
	}
	ctx.command = child
	return child.Fn(ctx)
}

// Usage sets usage and returns it
func (cmd *Command) Usage() string {
	// get usage form cache
	if cmd.usage != "" {
		return cmd.usage
	}
	buff := bytes.NewBufferString("")
	if cmd.Desc != "" {
		fmt.Fprintf(buff, "%s\n\n", cmd.Desc)
	}
	if cmd.Text != "" {
		fmt.Fprintf(buff, "%s\n\n", cmd.Text)
	}
	fmt.Fprintf(buff, "Usage:\n%s\n", usage(cmd.Argv()))
	if cmd.children != nil && len(cmd.children) > 0 {
		fmt.Fprintf(buff, "Commands:\n%v", cmd.ListChildren("  ", "   "))
	}
	cmd.usage = buff.String()
	return cmd.usage
}

// Path returns command full name
func (cmd *Command) Path() string {
	path := cmd.Name
	cur := cmd
	for cur.parent != nil {
		cur = cur.parent
		path = cur.Name + " " + path
	}
	return path
}

func (cmd *Command) route(router []string) *Command {
	cur := cmd
	for _, name := range router {
		child := cur.findChild(name)
		if child == nil {
			return nil
		}
		cur = child
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

// ListChildren returns all children's brief infos
func (cmd *Command) ListChildren(prefix, indent string) string {
	if cmd.children == nil || len(cmd.children) == 0 {
		return ""
	}
	buff := bytes.NewBufferString("")
	length := 0
	for _, child := range cmd.children {
		if len(child.Name) > length {
			length = len(child.Name)
		}
	}
	format := fmt.Sprintf("%s%%-%ds%s%%s\n", prefix, length, indent)
	for _, child := range cmd.children {
		fmt.Fprintf(buff, format, child.Name, child.Desc)
	}
	return buff.String()
}
