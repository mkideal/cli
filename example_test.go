package cli_test

import (
	"os"

	"github.com/mkideal/cli"
)

// This example demonstrates how to use short and long format flag
func ExampleParse_shortAndLongFlagName() {
	// argument object
	type argT struct {
		Port int `cli:"p,port" usage:"listening port"`
	}

	for _, args := range [][]string{
		[]string{"app", "-p", "8080"},
		[]string{"app", "-p8081"},
		[]string{"app", "-p=8082"},
		[]string{"app", "--port", "8083"},
		[]string{"app", "--port=8084"},
	} {
		cli.RunWithArgs(&argT{}, args, func(ctx *cli.Context) error {
			argv := ctx.Argv().(*argT)
			ctx.String("port=%d\n", argv.Port)
			return nil
		})
	}
	// Output:
	// port=8080
	// port=8081
	// port=8082
	// port=8083
	// port=8084
}

// This example demonstrates how to use default value
func ExampleParse_defaultValue() {
	type argT1 struct {
		Port int `cli:"p,port" usage:"listening port" dft:"8080"`
	}
	type argT2 struct {
		Port int `cli:"p,port" usage:"listening port" dft:"$CLI_TEST_HTTP_PORT"`
	}
	type argT3 struct {
		Port int `cli:"p,port" usage:"listening port" dft:"$CLI_TEST_HTTP_PORT+800"`
	}
	type argT4 struct {
		DevDir string `cli:"dir" usage:"develope directory" dft:"$CLI_TEST_DEV_PARENT_DIR/dev"`
	}

	os.Setenv("CLI_TEST_DEV_PARENT_DIR", "/home")
	os.Setenv("CLI_TEST_HTTP_PORT", "8000")

	for _, tt := range []struct {
		argv interface{}
		args []string
	}{
		{new(argT1), []string{"app"}},
		{new(argT2), []string{"app"}},
		{new(argT3), []string{"app"}},
		{new(argT4), []string{"app"}},
		{new(argT4), []string{"app", "--dir=/dev"}},
	} {
		cli.RunWithArgs(tt.argv, tt.args, func(ctx *cli.Context) error {
			ctx.String("argv=%v\n", ctx.Argv())
			return nil
		})
	}
	// Output:
	// argv=&{8080}
	// argv=&{8000}
	// argv=&{8800}
	// argv=&{/home/dev}
	// argv=&{/dev}
}

// This example demonstrates to use Slice and Map
func ExampleParse_sliceAndMap() {
	type argT1 struct {
		Slice []uint32 `cli:"U,u32-slice" usage:"uint32 slice"`
	}
	type argT2 struct {
		Slice []string `cli:"S,str-slice" usage:"string slice"`
	}
	type argT3 struct {
		Slice []bool `cli:"B,bool-slice" usage:"boolean slice"`
	}
	type argT4 struct {
		MapA map[string]int  `cli:"A" usage:"string => int"`
		MapB map[int]int     `cli:"B" usage:"int => int"`
		MapC map[int]string  `cli:"C" usage:"int => string"`
		MapD map[string]bool `cli:"D" usage:"string => bool"`
	}

	for _, tt := range []struct {
		argv interface{}
		args []string
	}{
		{new(argT1), []string{"app", "-U1", "-U2"}},
		{new(argT1), []string{"app", "-U", "1", "-U", "2"}},
		{new(argT1), []string{"app", "--u32-slice", "1", "--u32-slice", "2"}},
		{new(argT2), []string{"app", "-Shello", "-Sworld"}},
		{new(argT2), []string{"app", "-S", "hello", "-S", "world"}},
		{new(argT2), []string{"app", "--str-slice", "hello", "--str-slice", "world"}},
		{new(argT3), []string{"app", "-Btrue", "-Bfalse"}},
		{new(argT3), []string{"app", "-B", "true", "-B", "false"}},
		{new(argT3), []string{"app", "--bool-slice", "true", "--bool-slice", "false"}},

		{new(argT4), []string{"app",
			"-Ax=1",
			"-B", "1=2",
			"-C1=a",
			"-Dx",
		}},
	} {
		cli.RunWithArgs(tt.argv, tt.args, func(ctx *cli.Context) error {
			ctx.String("argv=%v\n", ctx.Argv())
			return nil
		})
	}
	// Output:
	// argv=&{[1 2]}
	// argv=&{[1 2]}
	// argv=&{[1 2]}
	// argv=&{[hello world]}
	// argv=&{[hello world]}
	// argv=&{[hello world]}
	// argv=&{[true false]}
	// argv=&{[true false]}
	// argv=&{[true false]}
	// argv=&{map[x:1] map[1:2] map[1:a] map[x:true]}
}
