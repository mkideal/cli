package cli

import (
	"testing"
)

func TestCommand(t *testing.T) {
	app := New("test")

	type arg_t struct {
		Help    bool   `cli:"h,help" usage:"show help"`
		Version string `cli:"v,version" usage:"show version" dft:"v0.0.0"`
	}

	sub1 := app.Register(&Command{
		Name: "sub1",
		Fn: func(ctx *Context) error {
			if ctx.Path() != "sub1" {
				t.Errorf("path: `%s` vs `%s`", ctx.Path(), "sub1")
			}
			argv := ctx.Argv().(*arg_t)
			if argv.Help != true || argv.Version != "v0.0.0" {
				t.Errorf("argv=%v", *argv)
			}
			return nil
		},
		ArgvFn: func() interface{} { return new(arg_t) },
	})

	sub1.Register(&Command{
		Name: "sub11",
		Fn: func(ctx *Context) error {
			if ctx.Path() != "sub1 sub11" {
				t.Errorf("path: `%s` vs `%s`", ctx.Path(), "sub1 sub11")
			}
			argv := ctx.Argv().(*arg_t)
			if argv.Help != false || argv.Version != "v1.0.0" {
				t.Errorf("argv=%v", *argv)
			}
			return nil
		},
		ArgvFn: func() interface{} { return new(arg_t) },
	})

	if err := app.Run([]string{
		"sub1",
		"-h",
	}); err != nil {
		t.Errorf("Run `sub1` error: %v", err)
	}

	if err := app.Run([]string{
		"sub1",
		"sub11",
		"--version=v1.0.0",
	}); err != nil {
		t.Errorf("Run `sub1 sub11` error: %v", err)
	}
}
