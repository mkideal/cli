package cli

import (
	"testing"
)

func TestCommandTree(t *testing.T) {
	app := &Command{}

	type argT struct {
		Help    bool   `cli:"h,help" usage:"show help"`
		Version string `cli:"v,version" usage:"show version" dft:"v0.0.0"`
	}

	sub1 := app.Register(&Command{
		Name: "sub1",
		Fn: func(ctx *Context) error {
			if ctx.Path() != "sub1" {
				t.Errorf("path: `%s` vs `%s`", ctx.Path(), "sub1")
			}
			argv := ctx.Argv().(*argT)
			if argv.Help != true || argv.Version != "v0.0.0" {
				t.Errorf("argv=%v", *argv)
			}
			if ctx.Command().Name != "sub1" {
				t.Errorf("command name want %s, got %s", "sub1", ctx.Command().Name)
			}
			return nil
		},
		Desc: "sub1 command describe",
		Argv: func() interface{} { return new(argT) },
	})

	sub1.Register(&Command{
		Name: "sub11",
		Fn: func(ctx *Context) error {
			if ctx.Path() != "sub1 sub11" {
				t.Errorf("path: `%s` vs `%s`", ctx.Path(), "sub1 sub11")
			}
			argv := ctx.Argv().(*argT)
			if argv.Help != false || argv.Version != "v1.0.0" {
				t.Errorf("argv=%v", *argv)
			}
			return nil
		},
		Desc: "sub11 desc",
		Text: "sub11 text",
		Argv: func() interface{} { return new(argT) },
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

	listWant := "sub11 sub11 desc\n"
	if listGot := sub1.ChildrenDescriptions("", " "); listGot != listWant {
		t.Errorf("ChildrenDescriptions want `%s`, got `%s`", listWant, listGot)
	}
}
