package cli_test

import (
	"reflect"

	"github.com/mkideal/cli"
)

type myParser struct {
	ptr interface{}
}

func newMyParser(ptr interface{}) cli.FlagParser {
	return &myParser{ptr}
}

// Parse implements FlagParser.Parse interface
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

type config3 struct {
	A int
	B string
}

// This example demonstrates how to use custom parser
func ExampleRegisterFlagParser() {
	// register parser factory function
	cli.RegisterFlagParser("myparser", newMyParser)

	type argT struct {
		Cfg3 config3 `cli:"cfg3" parser:"myparser"`
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
