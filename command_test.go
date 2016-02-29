package cli

import (
	"reflect"
	"sort"
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
	if listGot := sub1.ListChildren("", " "); listGot != listWant {
		t.Errorf("ListChildren want `%s`, got `%s`", listWant, listGot)
	}
}

func TestSuggestions(t *testing.T) {
	var root = &Command{}

	fn := func(ctx *Context) error { return nil }

	sub1 := root.Register(&Command{Name: "abc", Fn: fn})
	sub11 := sub1.Register(&Command{Name: "def", Fn: fn})
	sub12 := sub1.Register(&Command{Name: "deg", Fn: fn})
	sub2 := root.Register(&Command{Name: "abd", Fn: fn})
	_ = sub1
	_ = sub11
	_ = sub12
	_ = sub2

	for _, arg := range []struct {
		in   string
		want []string
	}{
		{"abc", []string{"abc"}},
		{"abe", []string{}},
		{"abc def", []string{"abc", "abc def", "abd"}},
	} {
		got := root.Suggestions(arg.in)
		sort.Strings(got)
		if !reflect.DeepEqual(got, arg.want) {
			t.Errorf("Suggestions `%s` want `%s`, got `%s`", arg.in, arg.want, got)
		}
	}
}

func compareStrings(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
