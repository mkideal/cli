package main

import (
	"fmt"
	"github.com/mkideal/cli"
	"os"
)

type arg_t struct {
	OnlySingle     bool    `cli:"v" usage:"only single char"`
	ManySingle     string  `cli:"X,Y,abcd" usage:"many single char"`
	SingleAndMulti int     `cli:"s,single-and-multi" usage:"single and multi"`
	OnlyMulti      uint    `cli:"only-multi, " usage:"only multi"`
	Required       int8    `cli:"*required" usage:"required value"`
	Default        uint8   `cli:"id" usage:"default value" dft:"102"`
	Ignored        int16   `cli:"-" usage:"ignored field"`
	UnName         uint16  `usage:"unname field"`
	Int32          int32   `cli:"i32" usage:"type int32" dft:"123"`
	Uint32         uint32  `cli:"u32" usage:"type uint32"`
	Int64          int64   `cli:"i64" usage:"type int64"`
	Uint64         int64   `cli:"u64" usage:"type uint64"`
	Float32        float32 `cli:"f32" usage:"type float32"`
	Float64        float64 `cli:"f364" usage:"type float64"`
}

func main() {
	//fmt.Printf("Usage of `%s': \n%s", os.Args[0], cli.Usage(new(arg_t)))
	t := new(arg_t)
	flagSet := cli.Parse(os.Args, t)
	if flagSet.Error != nil {
		fmt.Printf("%v\n", flagSet.Error)
	}
	fmt.Printf("Usage of `%s': \n%s", os.Args[0], flagSet.Usage)
	fmt.Printf("url.Values: %v\n", flagSet.Values)
	fmt.Printf("T: %v\n", *t)
}
