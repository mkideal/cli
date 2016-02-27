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

	app.Register(&Command{
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

	err := app.Run([]string{
		"sub1",
		"-h",
	})
	if err != nil {
		t.Errorf("Run error: %v", err)
	}
}
