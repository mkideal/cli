package cli

func ExampleHelper() {
	type argT struct {
		Helper
	}
	RunWithArgs(new(argT), []string{"app", "-h"}, func(ctx *Context) error {
		return nil
	})
	// Output:
	// Options:
	//
	//   -h, --help   display help information
}
