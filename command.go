package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/labstack/gommon/color"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

var commandNameRegexp = regexp.MustCompile("a-zA-Z_0-9]+")

// IsValidCommandName validate name of command
func IsValidCommandName(commandName string) bool {
	return commandNameRegexp.MatchString(commandName)
}

type (
	// Context provide running context
	Context struct {
		router     []string
		path       string
		argv       interface{}
		nativeArgs []string
		flagSet    *flagSet
		command    *Command
		writer     io.Writer
		color      color.Color

		HTTPRequest  *http.Request
		HTTPResponse http.ResponseWriter
	}

	// Validator validate flag before running command
	Validator interface {
		Validate(*Context) error
	}

	// CommandFunc ...
	CommandFunc func(*Context) error

	// ArgvFunc ...
	ArgvFunc func() interface{}

	// Command is the top-level instance in command-line app
	Command struct {
		Name        string      // Command name
		Aliases     []string    // Command aliases name
		Desc        string      // Command abstract
		Text        string      // Command detailed description
		Fn          CommandFunc // Command handler
		Argv        ArgvFunc    // Command argument factory function
		CanSubRoute bool

		HTTPRouters []string
		HTTPMethods []string

		routersMap map[string]string

		parent   *Command
		children []*Command

		isServer bool

		locker     sync.Mutex // protect following data
		usage      string
		usageStyle UsageStyle

		NoHook bool

		// hooks for current command
		OnBefore func(*Context) error
		OnAfter  func(*Context) error

		// hooks for all commands if current command is root command
		OnRootPrepareError func(error) error
		OnRootBefore       func(*Context) error
		OnRootAfter        func(*Context) error
	}

	// CommandTree represents a tree of commands
	CommandTree struct {
		command *Command
		forest  []*CommandTree
	}
)

//---------
// Context
//---------

func newContext(path string, router, args []string, argv interface{}, clr color.Color) (*Context, error) {
	ctx := &Context{
		path:       path,
		router:     router,
		argv:       argv,
		nativeArgs: args,
		color:      clr,
		flagSet:    newFlagSet(),
	}
	if argv != nil {
		ctx.flagSet = parseArgv(args, argv, ctx.color)
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
// `./app hello world -a --xyz=1` will return ["-a" "--xyz=1"]
func (ctx *Context) Args() []string {
	return ctx.nativeArgs
}

// FreedomArgs return freedom args
// `./app hello world -a=1 abc xyz` will return ["abc" "xyz"]
func (ctx *Context) FreedomArgs() []string {
	return ctx.flagSet.args
}

// Argv returns parsed args object
func (ctx *Context) Argv() interface{} {
	return ctx.argv
}

// FormValues returns parsed args as url.Values
func (ctx *Context) FormValues() url.Values {
	if ctx.flagSet == nil {
		Panicf("ctx.flagSet == nil")
	}
	return ctx.flagSet.values
}

// Command returns current command instance
func (ctx *Context) Command() *Command {
	return ctx.command
}

// Usage returns current command's usage with current context
func (ctx *Context) Usage() string {
	return ctx.command.Usage(ctx)
}

// WriteUsage writes usage to writer
func (ctx *Context) WriteUsage() {
	ctx.String(ctx.Usage())
}

// Writer returns writer
func (ctx *Context) Writer() io.Writer {
	if ctx.writer == nil {
		ctx.writer = colorable.NewColorableStdout()
	}
	return ctx.writer
}

// Write implements io.Writer
func (ctx *Context) Write(data []byte) (n int, err error) {
	return ctx.Writer().Write(data)
}

// Color returns color instance
func (ctx *Context) Color() *color.Color {
	return &ctx.color
}

// String writes formatted string to writer
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

// Register registers a child command
func (cmd *Command) Register(child *Command) *Command {
	if child == nil {
		Panicf("command `%s` try register a nil command", cmd.Name)
	}
	if child.Name == "" {
		Panicf("command `%s` try register a empty command", cmd.Name)
	}
	if cmd.children == nil {
		cmd.children = []*Command{}
	}
	if child.parent != nil {
		Panicf("command `%s` has been child of `%s`", child.Name, child.parent.Name)
	}
	if cmd.findChild(child.Name) != nil {
		Panicf("repeat register child `%s` for command `%s`", child.Name, cmd.Name)
	}
	if child.Aliases != nil {
		for _, alias := range child.Aliases {
			if cmd.findChild(alias) != nil {
				Panicf("repeat register child `%s` for command `%s`", alias, cmd.Name)
			}
		}
	}
	cmd.children = append(cmd.children, child)
	child.parent = cmd

	return child
}

// RegisterFunc registers handler as child command
func (cmd *Command) RegisterFunc(name string, fn CommandFunc, argvFn ArgvFunc) *Command {
	return cmd.Register(&Command{Name: name, Fn: fn, Argv: argvFn})
}

// RegisterTree registers a command tree
func (cmd *Command) RegisterTree(forest ...*CommandTree) {
	for _, tree := range forest {
		cmd.Register(tree.command)
		if tree.forest != nil && len(tree.forest) > 0 {
			tree.command.RegisterTree(tree.forest...)
		}
	}
}

// Parent returns command's parent
func (cmd *Command) Parent() *Command {
	return cmd.parent
}

// IsServer returns command whether if run as server
func (cmd *Command) IsServer() bool {
	return cmd.isServer
}

// IsClient returns command whether if run as client
func (cmd *Command) IsClient() bool {
	return !cmd.IsServer()
}

// SetIsServer sets command running mode(server or not)
func (cmd *Command) SetIsServer(yes bool) {
	cmd.Root().isServer = yes
}

// Run runs the command with args
func (cmd *Command) Run(args []string) error {
	return cmd.RunWith(args, nil, nil)
}

// RunWith runs the command with args and writer,httpMethods
func (cmd *Command) RunWith(args []string, writer io.Writer, resp http.ResponseWriter, httpMethods ...string) error {
	fds := []uintptr{}
	if writer == nil {
		writer = colorable.NewColorableStdout()
		fds = append(fds, os.Stdout.Fd())
	}
	clr := color.Color{}
	colorSwitch(&clr, writer, fds...)

	var ctx *Context
	var suggestion string
	ctx, suggestion, err := cmd.prepare(clr, args, writer, resp, httpMethods...)

	if err != nil {
		if cmd.OnRootPrepareError != nil {
			err = cmd.OnRootPrepareError(err)
		}
		if err != nil {
			return wrapErr(err, suggestion, clr)
		}
		return nil
	}

	if argv := ctx.Argv(); argv != nil {
		Debugf("command %s ready exec with argv %v", ctx.command.Name, argv)
	} else {
		Debugf("command %s ready exec", ctx.command.Name)
	}

	if ctx.command.NoHook {
		return ctx.command.Fn(ctx)
	}

	funcs := []func(*Context) error{
		ctx.command.OnBefore,
		cmd.OnRootBefore,
		ctx.command.Fn,
		cmd.OnRootAfter,
		ctx.command.OnAfter,
	}
	for _, f := range funcs {
		if f != nil {
			if err := f(ctx); err != nil {
				if err == ExitError {
					return nil
				}
				return err
			}
		}
	}
	return nil
}

func (cmd *Command) prepare(clr color.Color, args []string, writer io.Writer, resp http.ResponseWriter, httpMethods ...string) (ctx *Context, suggestion string, err error) {
	// split args
	router := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, dashOne) {
			break
		}
		router = append(router, arg)
	}
	if len(router) == 0 && cmd.Fn == nil {
		err = throwCommandNotFound(clr.Yellow(cmd.Name))
		return
	}
	path := strings.Join(router, " ")
	child, end := cmd.SubRoute(router)

	// if route fail
	if child == nil || (!child.CanSubRoute && end != len(router)) {
		suggestions := cmd.Suggestions(path)
		buff := bytes.NewBufferString("")
		if suggestions != nil && len(suggestions) > 0 {
			if len(suggestions) == 1 {
				fmt.Fprintf(buff, "\nDid you mean %s?", clr.Bold(suggestions[0]))
			} else {
				fmt.Fprintf(buff, "\n\nDid you mean one of these?\n")
				for _, sug := range suggestions {
					fmt.Fprintf(buff, "    %s\n", sug)
				}
			}
		}
		suggestion = buff.String()
		err = throwCommandNotFound(clr.Yellow(path))
		return
	}

	methodAllowed := false
	if len(httpMethods) == 0 ||
		child.HTTPMethods == nil ||
		len(child.HTTPMethods) == 0 {
		methodAllowed = true
	} else {
		method := httpMethods[0]
		for _, m := range child.HTTPMethods {
			if method == m {
				methodAllowed = true
				break
			}
		}
	}
	if !methodAllowed {
		err = throwMethodNotAllowed(clr.Yellow(httpMethods[0]))
		return
	}

	// create argv
	var argv interface{}
	if child.Argv != nil {
		argv = child.Argv()
	}

	// create Context
	ctx, err = newContext(path, router[:end], args[end:], argv, clr)
	if err != nil {
		return
	}

	// validate argv if argv implements interface Validator
	if argv != nil && !ctx.flagSet.dontValidate {
		if validator, ok := argv.(Validator); ok {
			err = validator.Validate(ctx)
			if err != nil {
				return
			}
		}
	}

	ctx.command = child
	ctx.writer = writer
	ctx.HTTPResponse = resp
	return
}

// Usage returns the usage string of command
func (cmd *Command) Usage(ctx *Context) string {
	style := GetUsageStyle()
	clr := color.Color{}
	clr.Disable()
	if ctx != nil {
		clr = ctx.color
	}
	// get usage form cache
	cmd.locker.Lock()
	tmpUsage := cmd.usage
	usageStyle := cmd.usageStyle
	cmd.locker.Unlock()
	if tmpUsage != "" && usageStyle == style {
		Debugf("get usage of command %s from cache", clr.Bold(cmd.Name))
		return tmpUsage
	}
	buff := bytes.NewBufferString("")
	if cmd.Desc != "" {
		fmt.Fprintf(buff, "%s\n\n", cmd.Desc)
	}
	if cmd.Text != "" {
		fmt.Fprintf(buff, "%s\n\n", cmd.Text)
	}
	if cmd.Argv != nil {
		fmt.Fprintf(buff, "%s:\n\n%s", clr.Bold("Options"), usage(cmd.Argv(), clr, style))
	}
	if cmd.children != nil && len(cmd.children) > 0 {
		if cmd.Argv != nil {
			buff.WriteByte('\n')
		}
		fmt.Fprintf(buff, "%s:\n%v", clr.Bold("Commands"), cmd.ChildrenDescriptions("  ", "   "))
	}
	tmpUsage = buff.String()
	cmd.locker.Lock()
	cmd.usage = tmpUsage
	cmd.usageStyle = style
	cmd.locker.Unlock()
	return tmpUsage
}

// Path returns space-separated command full name
func (cmd *Command) Path() string {
	return cmd.pathWithSep(" ")
}

func (cmd *Command) pathWithSep(sep string) string {
	var (
		path = ""
		cur  = cmd
	)
	for cur.parent != nil {
		if cur.Name != "" {
			if path == "" {
				path = cur.Name
			} else {
				path = cur.Name + sep + path
			}
		}
		cur = cur.parent
	}
	return path
}

// Root returns command's ancestor
func (cmd *Command) Root() *Command {
	ancestor := cmd
	for ancestor.parent != nil {
		ancestor = ancestor.parent
	}
	return ancestor
}

// Route finds command full matching router
func (cmd *Command) Route(router []string) *Command {
	child, end := cmd.SubRoute(router)
	if end != len(router) {
		return nil
	}
	return child
}

// SubRoute finds command partial matching router
func (cmd *Command) SubRoute(router []string) (*Command, int) {
	cur := cmd
	for i, name := range router {
		child := cur.findChild(name)
		if child == nil {
			return cur, i
		}
		cur = child
	}
	return cur, len(router)
}

// findChild finds child command by name
func (cmd *Command) findChild(name string) *Command {
	if cmd.nochild() {
		return nil
	}
	for _, child := range cmd.children {
		if child.Name == name {
			return child
		}
		if child.Aliases != nil {
			for _, alias := range child.Aliases {
				if alias == name {
					return child
				}
			}
		}
	}
	return nil
}

// ListChildren returns all names of command children
func (cmd *Command) ListChildren() []string {
	if cmd.nochild() {
		return []string{}
	}

	ret := make([]string, 0, len(cmd.children))
	for _, child := range cmd.children {
		ret = append(ret, child.Name)
	}
	return ret
}

// ChildrenDescriptions returns all children's brief infos by one string
func (cmd *Command) ChildrenDescriptions(prefix, indent string) string {
	if cmd.nochild() {
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

func (cmd *Command) nochild() bool {
	return cmd.children == nil || len(cmd.children) == 0
}

// Suggestions returns all similar commands
func (cmd *Command) Suggestions(path string) []string {
	if cmd.parent != nil {
		return cmd.Root().Suggestions(path)
	}

	var (
		cmds    = []*Command{cmd}
		targets = []string{}
	)
	for len(cmds) > 0 {
		if cmds[0].nochild() {
			cmds = cmds[1:]
		} else {
			for _, child := range cmds[0].children {
				targets = append(targets, child.Path())
			}
			cmds = append(cmds[0].children, cmds[1:]...)
		}
	}

	dists := []editDistanceRank{}
	for i, size := 0, len(targets); i < size; i++ {
		if d, ok := match(path, targets[i]); ok {
			dists = append(dists, editDistanceRank{s: targets[i], d: d})
		}
	}
	sort.Sort(editDistanceRankSlice(dists))
	for i := 0; i < len(dists); i++ {
		targets[i] = dists[i].s
	}
	return targets[:len(dists)]
}

func colorSwitch(clr *color.Color, w io.Writer, fds ...uintptr) {
	clr.Disable()
	if len(fds) > 0 {
		if isatty.IsTerminal(fds[0]) {
			clr.Enable()
		}
	} else if w, ok := w.(*os.File); ok && isatty.IsTerminal(w.Fd()) {
		clr.Enable()
	}
}

// HelpCommandFn implements buildin help command function
func HelpCommandFn(ctx *Context) error {
	var (
		args   = ctx.Args()
		parent = ctx.Command().Parent()
	)
	if len(args) == 0 {
		ctx.String(parent.Usage(ctx))
		return nil
	}
	var (
		child = parent.Route(args)
		clr   = ctx.Color()
	)
	if child == nil {
		return fmt.Errorf("command %s not found", clr.Yellow(strings.Join(args, " ")))
	}
	ctx.String(child.Usage(ctx))
	return nil
}

// HelpCommand returns a buildin help command
func HelpCommand(desc string) *Command {
	return &Command{
		Name:        "help",
		Desc:        desc,
		CanSubRoute: true,
		NoHook:      true,
		Fn:          HelpCommandFn,
	}
}

// Daemon startup app as a daemon process, success if result from stderr has prefix successPrefix
func Daemon(ctx *Context, successPrefix string) error {
	cmd := exec.Command(os.Args[0], ctx.Args()...)
	serr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(serr)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	if strings.HasPrefix(line, successPrefix) {
		ctx.String(line)
		cmd.Process.Release()
	} else {
		cmd.Process.Kill()
		return fmt.Errorf(line)
	}
	return nil
}
