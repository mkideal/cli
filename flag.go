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
	args   []string

	flagMap map[string]*flag
	flags   []*flag

	dontValidate bool
}

func newFlagSet() *flagSet {
	return &flagSet{
		flagMap: make(map[string]*flag),
		flags:   []*flag{},
		values:  url.Values(make(map[string][]string)),
		args:    make([]string, 0),
	}
}

type flag struct {
	t reflect.StructField
	v reflect.Value

	assigned bool
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
		zero := reflect.Zero(fl.t.Type)
		if reflect.DeepEqual(zero.Interface(), fl.v.Interface()) {
			return fl.set("", dft, clr)
		}
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

func (fl *flag) isBoolean() bool {
	return fl.t.Type.Kind() == reflect.Bool
}

func (fl *flag) getBool() bool {
	if !fl.isBoolean() {
		return false
	}
	return fl.v.Bool()
}

func (fl *flag) set(actual, s string, clr color.Color) error {
	fl.assigned = true
	fl.actual = actual
	return setWithProperType(fl.t.Type, fl.v, s, clr, false)
}

func setWithProperType(typ reflect.Type, val reflect.Value, s string, clr color.Color, isSubField bool) error {
	kind := typ.Kind()
	switch kind {
	case reflect.Bool:
		if v, err := getBool(s, clr); err == nil {
			val.SetBool(v)
		} else {
			return err
		}

	case reflect.String:
		val.SetString(s)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, err := getInt(s, clr); err == nil {
			if minmaxIntCheck(kind, v) {
				val.SetInt(v)
			} else {
				return errors.New(clr.Red("value overflow"))
			}
		} else {
			return err
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, err := getUint(s, clr); err == nil {
			if minmaxUintCheck(kind, v) {
				val.SetUint(uint64(v))
			} else {
				return errors.New(clr.Red("value overflow"))
			}
		} else {
			return err
		}

	case reflect.Float32, reflect.Float64:
		if v, err := getFloat(s, clr); err == nil {
			if minmaxFloatCheck(kind, v) {
				val.SetFloat(float64(v))
			} else {
				return errors.New(clr.Red("value overflow"))
			}
		} else {
			return err
		}

	case reflect.Slice:
		if isSubField {
			return fmt.Errorf("unsupported type %s as sub field type", kind.String())
		}
		sliceOf := typ.Elem()
		if val.IsNil() {
			slice := reflect.MakeSlice(typ, 0, 4)
			val.Set(slice)
		}
		index := val.Len()
		sliceCap := val.Cap()
		if index+1 <= sliceCap {
			val.SetLen(index + 1)
		} else {
			slice := reflect.MakeSlice(typ, index+1, index+sliceCap/2+1)
			for k := 0; k < index; k++ {
				slice.Index(k).Set(val.Index(k))
			}
			val.Set(slice)
		}
		return setWithProperType(sliceOf, val.Index(index), s, clr, true)

	case reflect.Map:
		if isSubField {
			return fmt.Errorf("unsupported type %s as sub field type", kind.String())
		}
		ks, vs, err := splitKeyVal(s)
		if err != nil {
			return err
		}
		kt := typ.Key()
		vt := typ.Elem()
		if val.IsNil() {
			val.Set(reflect.MakeMap(typ))
		}
		mk, mv := reflect.New(kt), reflect.New(vt)
		if err := setWithProperType(kt, mk.Elem(), ks, clr, true); err != nil {
			return err
		}
		if err := setWithProperType(vt, mv.Elem(), vs, clr, true); err != nil {
			return err
		}
		val.SetMapIndex(mk.Elem(), mv.Elem())

	default:
		return fmt.Errorf("unsupported type of field: %s", kind.String())
	}
	return nil
}

func splitKeyVal(s string) (key, val string, err error) {
	if s == "" {
		err = fmt.Errorf("empty key,val pair")
		return
	}
	index := strings.Index(s, "=")
	if index == -1 {
		return s, "", nil
	}
	return s[:index], s[index+1:], nil
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
	if s == "true" || s == "yes" || s == "y" || s == "" {
		return true, nil
	}
	if s == "false" || s == "none" || s == "no" || s == "not" || s == "n" {
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

type UsageStyle int32

const (
	NormalStyle UsageStyle = iota
	ManualStyle
)

var defaultStyle = NormalStyle

func GetUsageStyle() UsageStyle {
	return defaultStyle
}

func SetUsageStyle(style UsageStyle) {
	defaultStyle = style
}

type flagSlice []*flag

func (fs flagSlice) String(clr color.Color) string {
	var (
		lenShort                 = 0
		lenLong                  = 0
		lenNameAndDefaultAndLong = 0
		lenSep                   = len(sepName)
		sepSpaces                = strings.Repeat(" ", lenSep)
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
		l += lenDft
		if tag.name != "" {
			l += len(tag.name) + 1 // 1=len("=")
		}
		if l > lenNameAndDefaultAndLong {
			lenNameAndDefaultAndLong = l
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
			nameStr     = ""
			usagePrefix = " "
		)
		if tag.defaultValue != "" {
			defaultStr = fmt.Sprintf("[=%s]", tag.defaultValue)
		}
		if tag.name != "" {
			nameStr = "=" + tag.name
		}
		if tag.required {
			usagePrefix = clr.Red("*")
		}
		usage := usagePrefix + tag.usage

		spaceSize := lenSep + lenNameAndDefaultAndLong
		spaceSize -= len(nameStr) + len(defaultStr) + len(longStr)

		if defaultStr != "" {
			defaultStr = clr.Grey(defaultStr)
		}
		if nameStr != "" {
			nameStr = "=" + clr.Bold(tag.name)
		}

		if longStr == "" {
			format = fmt.Sprintf("%%%ds%%s%s%%s", lenShort, sepSpaces)
			fillStr := fillSpaces(nameStr+defaultStr, spaceSize)
			fmt.Fprintf(buff, format+"\n", shortStr, fillStr, usage)
		} else {
			if shortStr == "" {
				format = fmt.Sprintf("%%%ds%%s%%s", lenShort+lenSep)
			} else {
				format = fmt.Sprintf("%%%ds%s%%s%%s", lenShort, sepName)
			}
			fillStr := fillSpaces(longStr+nameStr+defaultStr, spaceSize)
			fmt.Fprintf(buff, format+"\n", shortStr, fillStr, usage)
		}
	}
	return buff.String()
}

func fillSpaces(s string, spaceSize int) string {
	return s + strings.Repeat(" ", spaceSize)
}

func (fs flagSlice) StringWithStyle(clr color.Color, style UsageStyle) string {
	if style != ManualStyle {
		return fs.String(clr)
	}

	buf := bytes.NewBufferString("")
	linePrefix := "  "
	for i, fl := range fs {
		if i != 0 {
			buf.WriteString("\n")
		}
		names := strings.Join(append(fl.tag.shortNames, fl.tag.longNames...), sepName)
		buf.WriteString(linePrefix)
		buf.WriteString(clr.Bold(names))
		if fl.tag.name != "" {
			buf.WriteString("=" + clr.Bold(fl.tag.name))
		}
		if fl.tag.defaultValue != "" {
			buf.WriteString(clr.Grey(fmt.Sprintf("[=%s]", fl.tag.defaultValue)))
		}
		buf.WriteString("\n")
		buf.WriteString(linePrefix)
		buf.WriteString("    ")
		if fl.tag.required {
			buf.WriteString(clr.Red("*"))
		}
		buf.WriteString(fl.tag.usage)
		buf.WriteString("\n")
	}
	return buf.String()
}
