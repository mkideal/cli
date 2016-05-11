package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/Bowery/prompt"
	"github.com/labstack/gommon/color"
	"github.com/mkideal/pkg/expr"
)

type flagSet struct {
	err    error
	values url.Values
	args   []string

	flagMap map[string]*flag
	flags   []*flag

	hasForce bool
}

func newFlagSet() *flagSet {
	return &flagSet{
		flagMap: make(map[string]*flag),
		flags:   []*flag{},
		values:  url.Values(make(map[string][]string)),
		args:    make([]string, 0),
	}
}

func (fs *flagSet) readPrompt(w io.Writer, clr color.Color) {
	for _, fl := range fs.flags {
		if fl.isAssigned || fl.tag.prompt == "" {
			continue
		}
		// read ...
		prefix := fl.tag.prompt + ": "
		var (
			data string
			yes  bool
		)
		if fl.tag.isPassword {
			data, fs.err = prompt.Password(prefix)
			if fs.err == nil && data != "" {
				fl.set(data, data, clr)
			}
		} else if fl.isBoolean() {
			yes, fs.err = prompt.Ask(prefix)
			if fs.err == nil {
				fl.value.SetBool(yes)
			}
		} else if fl.tag.defaultValue != "" {
			data, fs.err = prompt.BasicDefault(prefix, fl.tag.defaultValue)
			if fs.err == nil {
				fl.set(data, data, clr)
			}
		} else {
			data, fs.err = prompt.Basic(prefix, fl.tag.required)
			if fs.err == nil {
				fl.set(data, data, clr)
			}
		}
		if fs.err != nil {
			return
		}
	}
}

type flag struct {
	field reflect.StructField
	value reflect.Value

	// isAssigned indicates wether the flag is set(contains default value)
	isAssigned bool

	// isAssigned indicates wether the flag is set
	isSet bool

	// tag properties
	tag tagProperty

	// actual flag name
	actualFlagName string

	isNeedDelaySet bool

	// last value for need delay set
	// flag maybe assigned too many, like:
	// -f xx -f yy -f zz
	// `zz` is the last value
	lastValue string
}

func newFlag(field reflect.StructField, value reflect.Value, tag *tagProperty, clr color.Color, dontSetValue bool) (fl *flag, err error) {
	fl = &flag{field: field, value: value}
	if !fl.value.CanSet() {
		return nil, fmt.Errorf("field %s can not set", clr.Bold(fl.field.Name))
	}
	fl.tag = *tag
	fl.isNeedDelaySet = fl.tag.parserCreator != nil ||
		(fl.field.Type.Kind() != reflect.Slice &&
			fl.field.Type.Kind() != reflect.Map)
	err = fl.init(clr, dontSetValue)
	return
}

func (fl *flag) init(clr color.Color, dontSetValue bool) error {
	isNumber := fl.isInteger() || fl.isFloat()
	dft, err := parseExpression(fl.tag.defaultValue, isNumber)
	if err != nil {
		return err
	}
	if isNumber {
		v, err := expr.Eval(dft, nil, nil)
		if err != nil {
			return err
		}
		if fl.isInteger() {
			dft = fmt.Sprintf("%d", int64(v))
		} else if fl.isFloat() {
			dft = fmt.Sprintf("%f", float64(v))
		}
	}
	if !dontSetValue && fl.tag.defaultValue != "" {
		zero := reflect.Zero(fl.field.Type)
		if reflect.DeepEqual(zero.Interface(), fl.value.Interface()) {
			return fl.setDefault(dft, clr)
		}
	}
	return nil
}

func isWordByte(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') ||
		b == '_'
}

func parseExpression(s string, isNumber bool) (string, error) {
	src := []byte(s)
	var expr bytes.Buffer

	escaping := false
	const escapeByte = '$'
	var envvar bytes.Buffer
	writeEnv := func(envName string) error {
		if envName == "" {
			return fmt.Errorf("unexpected end after %v", escapeByte)
		}
		env := os.Getenv(envName)
		if env == "" && isNumber {
			env = "0"
		}
		expr.WriteString(env)
		return nil
	}
	for i, b := range src {
		if b == escapeByte {
			if escaping && envvar.Len() == 0 {
				expr.WriteByte(b)
				escaping = false
			} else {
				escaping = true
				if i+1 == len(src) {
					return "", fmt.Errorf("unexpected end after %v", escapeByte)
				}
				envvar.Reset()
			}
			continue
		}
		if escaping {
			if isWordByte(b) {
				envvar.WriteByte(b)
				if i+1 == len(src) {
					if err := writeEnv(envvar.String()); err != nil {
						return "", err
					}
				}
			} else {
				if err := writeEnv(envvar.String()); err != nil {
					return "", err
				}
				expr.WriteByte(b)
				envvar.Reset()
				escaping = false
			}
		} else {
			expr.WriteByte(b)
		}
	}
	return expr.String(), nil
}

func (fl *flag) name() string {
	if fl.actualFlagName != "" {
		return fl.actualFlagName
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
	return fl.field.Type.Kind() == reflect.Bool
}

func (fl *flag) isInteger() bool {
	switch fl.field.Type.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return true
	}
	return false
}

func (fl *flag) isFloat() bool {
	kind := fl.field.Type.Kind()
	return kind == reflect.Float32 || kind == reflect.Float64
}

func (fl *flag) isString() bool {
	return fl.field.Type.Kind() == reflect.String
}

func (fl *flag) getBool() bool {
	if !fl.isBoolean() {
		return false
	}
	return fl.value.Bool()
}

func (fl *flag) setDefault(s string, clr color.Color) error {
	fl.isAssigned = true
	if fl.isNeedDelaySet {
		fl.lastValue = s
		return nil
	}
	return setWithProperType(fl, fl.field.Type, fl.value, s, clr, false)
}

func (fl *flag) set(actualFlagName, s string, clr color.Color) error {
	fl.isSet = true
	fl.isAssigned = true
	fl.actualFlagName = actualFlagName
	if fl.isNeedDelaySet {
		fl.lastValue = s
		return nil
	}
	return setWithProperType(fl, fl.field.Type, fl.value, s, clr, false)
}

func setWithProperType(fl *flag, typ reflect.Type, val reflect.Value, s string, clr color.Color, isSubField bool) error {
	kind := typ.Kind()

	// try parser first of all
	if fl.tag.parserCreator != nil && val.CanInterface() {
		if kind != reflect.Ptr && val.CanAddr() {
			val = val.Addr()
		}
		return fl.tag.parserCreator(val.Interface()).Parse(s)
	}

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
			return fmt.Errorf("unsupported type %s as sub field", kind.String())
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
		return setWithProperType(fl, sliceOf, val.Index(index), s, clr, true)

	case reflect.Map:
		if isSubField {
			return fmt.Errorf("unsupported type %s as sub field", kind.String())
		}
		keyString, valString, err := splitKeyVal(s)
		if err != nil {
			return err
		}
		keyType := typ.Key()
		valType := typ.Elem()
		if val.IsNil() {
			val.Set(reflect.MakeMap(typ))
		}
		k, v := reflect.New(keyType), reflect.New(valType)
		if err := setWithProperType(fl, keyType, k.Elem(), keyString, clr, true); err != nil {
			return err
		}
		if err := setWithProperType(fl, valType, v.Elem(), valString, clr, true); err != nil {
			return err
		}
		val.SetMapIndex(k.Elem(), v.Elem())

	default:
		if val.CanInterface() {
			if kind != reflect.Ptr && val.CanAddr() {
				val = val.Addr()
			}
			// try Decoder
			if i := val.Interface(); i != nil {
				if decoder, ok := i.(Decoder); ok {
					return decoder.Decode(s)
				}
			}
		}
		return fmt.Errorf("unsupported type: %s", kind.String())
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

// UsageStyle is style of usage
type UsageStyle int32

const (
	// NormalStyle : left-right
	NormalStyle UsageStyle = iota
	// ManualStyle : up-down
	ManualStyle
)

var defaultStyle = NormalStyle

// GetUsageStyle gets default style
func GetUsageStyle() UsageStyle {
	return defaultStyle
}

// SetUsageStyle sets default style
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
