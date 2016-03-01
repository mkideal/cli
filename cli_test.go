package cli

import (
	"fmt"
	"math"
	"os"
	"testing"
)

type argT struct {
	Short         bool   `cli:"s" usage:"short flag"`
	Short2        bool   `cli:"2" usage:"another short flag"`
	ShortAndLong  string `cli:"S,long" usage:"short and long flags"`
	ShortsAndLong int    `cli:"x,y,abcd,omitof" usage:"many short and long flags"`
	Long          uint   `cli:"long-flag" usage:"long flag"`
	Required      int8   `cli:"*required" usage:"required flag, note the *"`
	Default       uint8  `cli:"dft,default" usage:"default value" dft:"102"`
	UnName        uint16 `usage:"unname field"`

	Int8    int8    `cli:"i8" usage:"type int8"`
	Uint8   uint8   `cli:"u8" usage:"type uint8"`
	Int16   int16   `cli:"i16" usage:"type int16"`
	Uint16  uint16  `cli:"u16" usage:"type uint16"`
	Int32   int32   `cli:"i32" usage:"type int32"`
	Uint32  uint32  `cli:"u32" usage:"type uint32"`
	Int64   int64   `cli:"i64" usage:"type int64"`
	Uint64  uint64  `cli:"u64" usage:"type uint64"`
	Float32 float32 `cli:"f32" usage:"type float32"`
	Float64 float64 `cli:"f64" usage:"type float64"`
}

func toStr(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

func TestParse(t *testing.T) {
	for i, tab := range []struct {
		args  []string
		want  argT
		isErr bool
	}{
		//Case: ENV
		{
			args: []string{"--required=0"},
			want: argT{Default: 102},
		},
		//Case: missing required
		{
			args:  []string{},
			isErr: true,
		},
		//Case: undefined flag
		{
			args:  []string{"--required=0", "-Q"},
			isErr: true,
		},
		//Case: undefined flag
		{
			args:  []string{"--required=0", "--KdjiiejdfwkHJH"},
			isErr: true,
		},
		//Case: short flag group
		{
			args: []string{"--required=0", "-s2"},
			want: argT{Default: 102, Short: true, Short2: true},
		},
		//Case: check default
		{
			args: []string{"--required=0"},
			want: argT{Default: 102},
		},
		//Case: modify default
		{
			args: []string{"--required=0", "--dft", "55"},
			want: argT{Default: 55},
		},
		//Case: modify default
		{
			args: []string{"--required=0", "--default", "55"},
			want: argT{Default: 55},
		},
		//Case: UnName
		{
			args: []string{"--required=0", "--UnName", "64"},
			want: argT{Default: 102, UnName: 64},
		},
		//Case: not a bool
		{
			args:  []string{"--required=0", "-s", "not-a-bool"},
			isErr: true,
		},
		//Case: "" -> bool
		{
			args: []string{"--required=0", "-s"},
			want: argT{Default: 102, Short: true},
		},
		//Case: "true" -> bool
		{
			args: []string{"--required=0", "-s", "true"},
			want: argT{Default: 102, Short: true},
		},
		//Case: non-zero integer -> bool
		{
			args: []string{"--required=0", "-s", "1"},
			want: argT{Default: 102, Short: true},
		},
		//Case: zero -> bool
		{
			args: []string{"--required=0", "-s", "0"},
			want: argT{Default: 102},
		},
		//Case: no -> bool
		{
			args: []string{"--required=0", "-s", "no"},
			want: argT{Default: 102},
		},
		//Case: not -> bool
		{
			args: []string{"--required=0", "-s", "not"},
			want: argT{Default: 102},
		},
		//Case: none -> bool
		{
			args: []string{"--required=0", "-s", "none"},
			want: argT{Default: 102},
		},
		//Case: false -> bool
		{
			args: []string{"--required=0", "-s", "false"},
			want: argT{Default: 102},
		},
		//Case: int64
		{
			args: []string{"--required=0", "--i64", toStr(12)},
			want: argT{Default: 102, Int64: 12},
		},
		//Case: int64 overflow
		{
			args:  []string{"--required=0", "--i64", toStr(uint64(math.MaxUint64))},
			isErr: true,
		},
		//Case: uint64
		{
			args: []string{"--required=0", "--u64", toStr(12)},
			want: argT{Default: 102, Uint64: 12},
		},
		//Case: max uint64
		{
			args: []string{"--required=0", "--u64", toStr(uint64(math.MaxUint64))},
			want: argT{Default: 102, Uint64: uint64(math.MaxUint64)},
		},
		//Case: negative -> uint64
		{
			args:  []string{"--required=0", "--u64", "-1"},
			isErr: true,
		},
		//Case: int32
		{
			args: []string{"--required=0", "--i32", toStr(12)},
			want: argT{Default: 102, Int32: 12},
		},
		//Case: int32 overflow
		{
			args:  []string{"--required=0", "--i32", toStr(uint32(math.MaxUint32))},
			isErr: true,
		},
		//Case: uint32
		{
			args: []string{"--required=0", "--u32", toStr(12)},
			want: argT{Default: 102, Uint32: 12},
		},
		//Case: max uint32
		{
			args: []string{"--required=0", "--u32", toStr(uint32(math.MaxUint32))},
			want: argT{Default: 102, Uint32: uint32(math.MaxUint32)},
		},
		//Case: negative -> uint32
		{
			args:  []string{"--required=0", "--u32", "-1"},
			isErr: true,
		},
		//Case: int16
		{
			args: []string{"--required=0", "--i16", toStr(12)},
			want: argT{Default: 102, Int16: 12},
		},
		//Case: int16 overflow
		{
			args:  []string{"--required=0", "--i16", toStr(uint16(math.MaxUint16))},
			isErr: true,
		},
		//Case: uint16
		{
			args: []string{"--required=0", "--u16", toStr(12)},
			want: argT{Default: 102, Uint16: 12},
		},
		//Case: max uint16
		{
			args: []string{"--required=0", "--u16", toStr(uint16(math.MaxUint16))},
			want: argT{Default: 102, Uint16: uint16(math.MaxUint16)},
		},
		//Case: negative -> uint16
		{
			args:  []string{"--required=0", "--u16", "-1"},
			isErr: true,
		},
		//Case: int8
		{
			args: []string{"--required=0", "--i8", toStr(12)},
			want: argT{Default: 102, Int8: 12},
		},
		//Case: int8 overflow
		{
			args:  []string{"--required=0", "--i8", toStr(uint8(math.MaxUint8))},
			isErr: true,
		},
		//Case: uint8
		{
			args: []string{"--required=0", "--u8", toStr(12)},
			want: argT{Default: 102, Uint8: 12},
		},
		//Case: max uint8
		{
			args: []string{"--required=0", "--u8", toStr(uint8(math.MaxUint8))},
			want: argT{Default: 102, Uint8: uint8(math.MaxUint8)},
		},
		//Case: negative -> uint8
		{
			args:  []string{"--required=0", "--u8", "-1"},
			isErr: true,
		},
		//Case: many invalid value
		{
			args:  []string{"--required=0", "--u8", "-1", "--u8", "256"},
			isErr: true,
		},
		//Case: too many value
		{
			args:  []string{"--required=0=1"},
			isErr: true,
		},
		//Case: float32
		{
			args: []string{"--required=0", "--f32", "12.34"},
			want: argT{Default: 102, Float32: 12.34},
		},
		//Case: not a float32
		{
			args:  []string{"--required=0", "--f32", "not-a-float32"},
			isErr: true,
		},
		//Case: float32 overflow
		{
			args:  []string{"--required=0", "--f32", "123456789123456789123456789123456789123456789"},
			isErr: true,
		},
		//Case: float32 overflow
		{
			args:  []string{"--required=0", "--f32", "-123456789123456789123456789123456789123456789"},
			isErr: true,
		},
		//Case: float64
		{
			args: []string{"--required=0", "--f64=-1234.5678"},
			want: argT{Default: 102, Float64: -1234.5678},
		},
		//Case: not a float64
		{
			args:  []string{"--required=0", "--f64=not-a-float64"},
			isErr: true,
		},
	} {
		v := new(argT)
		flagSet := parseArgv(tab.args, v)
		if tab.isErr {
			if flagSet.err == nil {
				t.Errorf("[%2d] want error, got nil", i)
			}
			continue
		}
		if flagSet.err != nil {
			t.Errorf("[%2d] parseArgv error: %v", i, flagSet.err)
			continue
		}
		if *v != tab.want {
			t.Errorf("[%2d] want %v, got %v", i, tab.want, *v)
		}
	}

	//Case parse non-pointer object
	if flagSet := parseArgv([]string{}, argT{}); flagSet.err != errNotAPointer {
		t.Errorf("want %v, got %v", errNotAPointer, flagSet.err)
	}
	if usage(argT{}) != "" {
		t.Errorf("want usage empty, but not")
	}

	//Case parse pointer, but not indirect a struct
	tmp := 0
	ptrInt := &tmp
	if flagSet := parseArgv([]string{}, ptrInt); flagSet.err != errNotPointToStruct {
		t.Errorf("want %v, got %v", errNotPointToStruct, flagSet.err)
	}
	if usage(ptrInt) != "" {
		t.Errorf("want usage empty, but not")
	}

	//Case repeat tag
	type tmpT struct {
		A bool `cli:"a"`
		B bool `cli:"a"`
	}
	if flagSet := parseArgv([]string{}, new(tmpT)); flagSet.err == nil {
		t.Errorf("want error, got nil")
	}
	if usage(new(tmpT)) != "" {
		t.Errorf("want usage empty, but not")
	}

	type envT struct {
		DefaultEnv string `cli:"default-env" usage:"default value" dft:"$ENV_CLI_TEST"`
	}
	envV := new(envT)
	if flagSet := parseArgv([]string{}, envV); flagSet.err != nil {
		t.Errorf(flagSet.err.Error())
	} else {
		if want := os.Getenv("ENV_CLI_TEST"); want != envV.DefaultEnv {
			t.Errorf("ENV_CLI_TEST want `%s`, got `%s`", want, envV.DefaultEnv)
		}
	}
}

func TestUsage(t *testing.T) {
	usage := usage(new(argT))
	want := fmt.Sprintf(
		`      -s                             short flag
      -2                             another short flag
      -S, --long                     short and long flags
  -x, -y, --abcd, --omitof           many short and long flags
          --long-flag                long flag
          --required                %srequired flag, note the *
          --dft, --default%s     default value
          --UnName                   unname field
          --i8                       type int8
          --u8                       type uint8
          --i16                      type int16
          --u16                      type uint16
          --i32                      type int32
          --u32                      type uint32
          --i64                      type int64
          --u64                      type uint64
          --f32                      type float32
          --f64                      type float64
`, Red("*"), Gray("[=102]"))
	if usage != want {
		t.Errorf("usage want `%s`, got `%s`", want, usage)
	}
}
