package cli_test

import (
	"io/ioutil"
	"os"
	"reflect"

	"github.com/mkideal/cli"
)

type Config1 struct {
	A string
	B int
}

type Config2 struct {
	C string
	D bool
}

// This example demonstrates how to use builtin praser(json,jsonfile)
func Example_parser() {
	type argT struct {
		Cfg1 Config1 `cli:"cfg1" parser:"json"`
		Cfg2 Config2 `cli:"cfg2" parser:"jsonfile"`
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

type myParser struct {
	ptr interface{}
}

func newMyParser(ptr interface{}) cli.FlagParser {
	return &myParser{ptr}
}

func (parser *myParser) Parse(s string) error {
	typ := reflect.TypeOf(parser.ptr)
	val := reflect.ValueOf(parser.ptr)
	if typ.Kind() == reflect.Ptr {
		kind := reflect.Indirect(val).Type().Kind()
		if kind == reflect.Struct {
			typElem, valElem := typ.Elem(), val.Elem()
			numField := valElem.NumField()
			for i := 0; i < numField; i++ {
				_, valField := typElem.Field(i), valElem.Field(i)
				if valField.Kind() == reflect.Int && valField.CanSet() {
					valField.SetInt(2)
				}
				if valField.Kind() == reflect.String && valField.CanSet() {
					valField.SetString("B")
				}
			}
		}
	}
	return nil
}

type Config3 struct {
	A int
	B string
}

// This example demonstrates how to use custom parser
func Example_customParser() {
	// register parser factory function
	cli.RegisterFlagParser("myparser", newMyParser)

	type argT struct {
		Cfg3 Config3 `cli:"cfg3" parser:"myparser"`
	}

	args := []string{"app",
		`--cfg3`, `hello`,
	}

	cli.RunWithArgs(new(argT), args, func(ctx *cli.Context) error {
		ctx.JSON(ctx.Argv())
		return nil
	})
	// Output:
	// {"Cfg3":{"A":2,"B":"B"}}
}
