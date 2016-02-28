package main

import (
	"github.com/mkideal/cli"
)

type arg_t struct {
	Short         bool   `cli:"s" usage:"short flag"`
	ShortAndLong  string `cli:"S,long" usage:"short and long flags"`
	ShortsAndLong int    `cli:"x,y,abcd,omitof" usage:"many short and long flags"`
	Long          uint   `cli:"long-flag" usage:"long flag"`
	Required      int8   `cli:"*required" usage:"required flag, note the *"`
	Default       uint8  `cli:"dft" usage:"default value" dft:"102"`

	// Ignored field: spacify `cli` tag=-
	Ignored int16 `cli:"-" usage:"ignored field"`

	// No spacify `cli` tag value
	// so it's flag=`--UnName`
	UnName uint16 `usage:"unname field"`

	Int32   int32   `cli:"i32" usage:"type int32" dft:"123"`
	Uint32  uint32  `cli:"u32" usage:"type uint32"`
	Int64   int64   `cli:"i64" usage:"type int64"`
	Uint64  int64   `cli:"u64" usage:"type uint64"`
	Float32 float32 `cli:"f32" usage:"type float32"`
	Float64 float64 `cli:"f64" usage:"type float64"`
}

func main() {
	cli.Run(new(arg_t), func(ctx *cli.Context) error {
		// Show usage information
		ctx.String(ctx.Usage())

		// Get argv
		argv := ctx.Argv().(*arg_t)

		// Show json object
		ctx.JSON(argv).String("\n")
		ctx.JSONIndent(argv, "", "    ").String("\n")

		// Get show native args
		ctx.JSONln(ctx.Args())

		// Show the args as url.Values
		ctx.JSONIndentln(ctx.FormValues(), "", "    ")

		return nil
	})
}
