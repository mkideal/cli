package cli

import (
	"fmt"
	"testing"
)

type arg_t struct {
	OnlySingle     bool    `cli:"v" usage:"only single char"`
	ManySingle     string  `cli:"X,Y" usage:"many single char"`
	SingleAndMulti int     `cli:"s,single-and-multi" usage:"single and multi"`
	OnlyMulti      uint    `cli:"only-multi, " usage:"only multi"`
	Required       int8    `cli:"*required" usage:"required value"`
	Default        uint8   `cli:"id" usage:"default value" dft:"1024"`
	Ignored        int16   `cli:"-" usage:"ignored field"`
	UnName         uint16  `usage:"unname field"`
	Int32          int32   `cli:"i32" usage:"type int32"`
	Uint32         uint32  `cli:"u32" usage:"type uint32"`
	Int64          int64   `cli:"i64" usage:"type int64"`
	Uint64         int64   `cli:"u64" usage:"type uint64"`
	Float32        float32 `cli:"f32" usage:"type float32"`
	Float64        float64 `cli:"f364" usage:"type float64"`
}

func TestUsage(t *testing.T) {
	args := []string{
		"--required=120",
	}
	v := new(arg_t)
	flagSet := Parse(args, v)
	if flagSet.Error != nil {
		t.Errorf("error: %v", flagSet.Error)
	}
	wantUsage := fmt.Sprintf(`      -v                         only single char
  -X, -Y                         many single char
      -s, --single-and-multi     single and multi
          --only-multi           only multi
          --required            %srequired value
          --id                   default value%s
          --UnName               unname field
          --i32                  type int32
          --u32                  type uint32
          --i64                  type int64
          --u64                  type uint64
          --f32                  type float32
          --f364                 type float64
`, red("*"), gray("[default=1024]"))
	if flagSet.Usage != wantUsage {
		t.Errorf("usage from `Parse` func want: `%s`, got `%s`", wantUsage, flagSet.Usage)
	}
	if u := Usage(v); u != wantUsage {
		t.Errorf("usage from `Usage` func want: `%s`, got `%s`", wantUsage, u)
	}
}
