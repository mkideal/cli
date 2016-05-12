package cli_test

import (
	"io/ioutil"
	"os"

	"github.com/mkideal/cli"
)

type config1 struct {
	A string
	B int
}

type config2 struct {
	C string
	D bool
}

// This example demonstrates how to use builtin praser(json,jsonfile)
func ExampleFlagParser() {
	type argT struct {
		Cfg1 config1 `cli:"cfg1" parser:"json"`
		Cfg2 config2 `cli:"cfg2" parser:"jsonfile"`
	}
	jsonfile := "1.json"
	args := []string{"app",
		`--cfg1`, `{"A": "hello", "B": 2}`,
		`--cfg2`, jsonfile,
	}
	ioutil.WriteFile(jsonfile, []byte(`{"C": "world", "D": true}`), 0644)
	defer os.Remove(jsonfile)

	cli.RunWithArgs(new(argT), args, func(ctx *cli.Context) error {
		ctx.JSON(ctx.Argv())
		return nil
	})
	// Output:
	// {"Cfg1":{"A":"hello","B":2},"Cfg2":{"C":"world","D":true}}
}
