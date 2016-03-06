package cli

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/gommon/color"
)

type flagSet struct {
	err    error
	values url.Values

	flagMap map[string]*flag
	flags   []*flag

	dontValidate bool
}

func newFlagSet() *flagSet {
	return &flagSet{
		flagMap: make(map[string]*flag),
		flags:   []*flag{},
		values:  url.Values(make(map[string][]string)),
	}
}

type flag struct {
	t reflect.StructField
	v reflect.Value

	assigned bool
	err      error
	tag      fieldTag

	actual string
}

func newFlag(t reflect.StructField, v reflect.Value, tag *fieldTag, clr color.Color) (fl *flag, err error) {
	fl = &flag{t: t, v: v}
	if !fl.v.CanSet() {
		return nil, fmt.Errorf("field %s can not set", clr.Bold(fl.t.Name))
	}
	fl.tag = *tag
	err = fl.init(clr)
	return
}

func (fl *flag) init(clr color.Color) error {
	dft := fl.tag.defaultValue
	if strings.HasPrefix(dft, "$") {
		dft = dft[1:]
		if !strings.HasPrefix(dft, "$") {
			dft = os.Getenv(dft)
		}
	}
	if dft != "" {
		return fl.set("", dft, clr)
	}
	return nil
}

func (fl *flag) name() string {
	if fl.actual != "" {
		return fl.actual
	}
	if len(fl.tag.longNames) > 0 {
		return fl.tag.longNames[0]
	}
	if len(fl.tag.shortNames) > 0 {
		return fl.tag.shortNames[0]
	}
	return ""
}

func (fl *flag) getBool() bool {
	if fl.t.Type.Kind() != reflect.Bool {
		return false
	}
	return fl.v.Bool()
}

func (fl *flag) set(actual, s string, clr color.Color) error {
	kind := fl.t.Type.Kind()
	fl.assigned = true
	fl.actual = actual
	switch kind {
	case reflect.Bool:
		if v, err := getBool(s, clr); err == nil {
			fl.v.SetBool(v)
		} else {
			fl.err = err
		}

	case reflect.String:
		fl.v.SetString(s)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, err := getInt(s, clr); err == nil {
			if minmaxIntCheck(kind, v) {
				fl.v.SetInt(v)
			} else {
				fl.err = errors.New(clr.Red("value overflow"))
			}
		} else {
			fl.err = err
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, err := getUint(s, clr); err == nil {
			if minmaxUintCheck(kind, v) {
				fl.v.SetUint(uint64(v))
			} else {
				fl.err = errors.New(clr.Red("value overflow"))
			}
		} else {
			fl.err = err
		}

	case reflect.Float32, reflect.Float64:
		if v, err := getFloat(s, clr); err == nil {
			if minmaxFloatCheck(kind, v) {
				fl.v.SetFloat(float64(v))
			} else {
				fl.err = errors.New(clr.Red("value overflow"))
			}
		} else {
			fl.err = err
		}
	default:
		return fmt.Errorf("invalid field type: %s", kind.String())
	}
	return nil
}

func minmaxIntCheck(kind reflect.Kind, v int64) bool {
	switch kind {
	case reflect.Int, reflect.Int64:
		return v >= int64(math.MinInt64) && v <= int64(math.MaxInt64)
	case reflect.Int8:
		return v >= int64(math.MinInt8) && v <= int64(math.MaxInt8)
	case reflect.Int16:
		return v >= int64(math.MinInt16) && v <= int64(math.MaxInt16)
	case reflect.Int32:
		return v >= int64(math.MinInt32) && v <= int64(math.MaxInt32)
	}
	return true
}

func minmaxUintCheck(kind reflect.Kind, v uint64) bool {
	switch kind {
	case reflect.Uint, reflect.Uint64:
		return v <= math.MaxUint64
	case reflect.Uint8:
		return v <= math.MaxUint8
	case reflect.Uint16:
		return v <= math.MaxUint16
	case reflect.Uint32:
		return v <= math.MaxUint32
	}
	return true
}

func minmaxFloatCheck(kind reflect.Kind, v float64) bool {
	switch kind {
	case reflect.Float32:
		return v >= -float64(math.MaxFloat32) && v <= float64(math.MaxFloat32)
	case reflect.Float64:
		return v >= -float64(math.MaxFloat64) && v <= float64(math.MaxFloat64)
	}
	return true
}

func getBool(s string, clr color.Color) (bool, error) {
	if s == "true" || s == "" {
		return true, nil
	}
	if s == "false" || s == "none" || s == "no" || s == "not" {
		return false, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return false, fmt.Errorf("`%s` couldn't convert to a %s value", s, clr.Bold("bool"))
	}
	return i != 0, nil
}

func getInt(s string, clr color.Color) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("`%s` couldn't convert to an %s value", s, clr.Bold("int"))
	}
	return i, nil
}

func getUint(s string, clr color.Color) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("`%s` couldn't convert to an %s value", s, clr.Bold("uint"))
	}
	return i, nil
}

func getFloat(s string, clr color.Color) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("`%s` couldn't convert to a %s value", s, clr.Bold("float"))
	}
	return f, nil
}

type flagSlice []*flag

func (fs flagSlice) String(clr color.Color) string {
	var (
		lenShort          = 0
		lenLong           = 0
		lenDefaultAndLong = 0
		lenSep            = len(sepName)
		sepSpaces         = strings.Repeat(" ", lenSep)
	)
	for _, fl := range fs {
		tag := fl.tag
		l := 0
		for _, shortName := range tag.shortNames {
			l += len(shortName) + lenSep
		}
		if l > lenShort {
			lenShort = l
		}
		l = 0
		for _, longName := range tag.longNames {
			l += len(longName) + lenSep
		}
		if l > lenLong {
			lenLong = l
		}
		lenDft := 0
		if tag.defaultValue != "" {
			lenDft = len(tag.defaultValue) + 3 // 3=len("[=]")
		}
		if l+lenDft > lenDefaultAndLong {
			lenDefaultAndLong = l + lenDft
		}
	}

	buff := bytes.NewBufferString("")
	for _, fl := range fs {
		var (
			tag         = fl.tag
			shortStr    = strings.Join(tag.shortNames, sepName)
			longStr     = strings.Join(tag.longNames, sepName)
			format      = ""
			defaultStr  = ""
			usagePrefix = " "
		)
		if tag.defaultValue != "" {
			defaultStr = fmt.Sprintf("[=%s]", tag.defaultValue)
		}
		if tag.required {
			usagePrefix = clr.Red("*")
		}
		usage := usagePrefix + tag.usage

		spaceSize := lenSep + lenDefaultAndLong - len(defaultStr) - len(longStr)
		if defaultStr != "" {
			defaultStr = clr.Grey(defaultStr)
		}
		if longStr == "" {
			format = fmt.Sprintf("%%%ds%%s%s%%s", lenShort, sepSpaces)
			fillStr := fillSpaces(defaultStr, spaceSize)
			buff.WriteString(fmt.Sprintf(format+"\n", shortStr, fillStr, usage))
		} else {
			if shortStr == "" {
				format = fmt.Sprintf("%%%ds%%s%%s", lenShort+lenSep)
			} else {
				format = fmt.Sprintf("%%%ds%s%%s%%s", lenShort, sepName)
			}
			fillStr := fillSpaces(longStr+defaultStr, spaceSize)
			buff.WriteString(fmt.Sprintf(format+"\n", shortStr, fillStr, usage))
		}
	}
	return buff.String()
}

func fillSpaces(s string, spaceSize int) string {
	return s + strings.Repeat(" ", spaceSize)
}
