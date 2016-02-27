package cli

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var (
	errValueOverflow = errors.New(red("value overflow"))
)

type FlagSet struct {
	Error  error
	Usage  string
	Values url.Values

	flags map[string]*flag
	slice []*flag
}

func newFlagSet() *FlagSet {
	return &FlagSet{
		flags:  make(map[string]*flag),
		slice:  []*flag{},
		Values: url.Values(make(map[string][]string)),
	}
}

type flag struct {
	t reflect.StructField
	v reflect.Value

	assigned bool
	err      error
	/*invalid     bool
	invalidDesc string*/
	tag cliTag
	typ string
}

func newFlag(t reflect.StructField, v reflect.Value) (fl *flag, err error) {
	fl = &flag{t: t, v: v}
	tag, err := parseTag(t.Name, t.Tag)
	if tag == nil {
		return nil, nil
	}
	fl.tag = *tag
	err = fl.init()
	fl.typ = t.Type.Kind().String()
	return
}

func (fl *flag) init() error {
	if fl.tag.defaultValue != "" {
		return fl.set(fl.tag.defaultValue)
	}
	return nil
}

func (fl *flag) name() string {
	if len(fl.tag.longNames) > 0 {
		return fl.tag.longNames[0]
	}
	if len(fl.tag.shortNames) > 0 {
		return fl.tag.shortNames[0]
	}
	return ""
}

func (fl *flag) set(s string) error {
	kind := fl.t.Type.Kind()
	fl.assigned = true
	switch kind {
	case reflect.Bool:
		if v, err := getBool(s); err == nil {
			fl.v.SetBool(v)
		} else {
			fl.err = err
		}

	case reflect.String:
		fl.v.SetString(s)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, err := getInt(s); err == nil {
			if minmaxIntCheck(kind, v) {
				fl.v.SetInt(v)
			} else {
				fl.err = errValueOverflow
			}
		} else {
			fl.err = err
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, err := getUint(s); err == nil {
			if minmaxUintCheck(kind, v) {
				fl.v.SetUint(uint64(v))
			} else {
				fl.err = errValueOverflow
			}
		} else {
			fl.err = err
		}

	case reflect.Float32, reflect.Float64:
		if v, err := getFloat(s); err == nil {
			if minmaxFloatCheck(kind, v) {
				fl.v.SetFloat(float64(v))
			} else {
				fl.err = errValueOverflow
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

func getBool(s string) (bool, error) {
	if s == "true" || s == "" {
		return true, nil
	}
	if s == "false" || s == "none" || s == "no" || s == "not" {
		return false, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return false, fmt.Errorf("`%s` couldn't convert to a %s value", s, red("bool"))
	}
	return i != 0, nil
}

func getInt(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("`%s` couldn't convert to an %s value", s, red("int"))
	}
	return i, nil
}

func getUint(s string) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("`%s` couldn't convert to an %s value", s, red("uint"))
	}
	return i, nil
}

func getFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("`%s` couldn't convert to a %s value", s, red("float"))
	}
	return f, nil
}

type flagSlice []*flag

func (fs flagSlice) String() string {
	var (
		lenShort  = 0
		lenLong   = 0
		lenSep    = len(sepName)
		sepSpaces = strings.Repeat(" ", lenSep)
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
	}

	buff := bytes.NewBufferString("")
	for _, fl := range fs {
		tag := fl.tag
		shortStr := strings.Join(tag.shortNames, sepName)
		longStr := strings.Join(tag.longNames, sepName)
		format := ""
		l1, l2 := lenShort+lenSep, lenLong+lenSep
		if shortStr == "" {
			format = fmt.Sprintf("%%%ds%%-%ds%%s", l1, l2)
		} else if longStr == "" {
			format = fmt.Sprintf("%%%ds%s%%-%ds%%s", lenShort, sepSpaces, l2)
		} else {
			format = fmt.Sprintf("%%%ds%s%%-%ds%%s", lenShort, sepName, l2)
		}
		usagePrefix := " "
		if tag.required {
			usagePrefix = red("*")
		}
		usage := usagePrefix + tag.usage
		if tag.defaultValue != "" {
			usage += gray("[default=%s]", tag.defaultValue)
		}
		buff.WriteString(fmt.Sprintf(format+"\n", shortStr, longStr, usage))
	}
	return buff.String()
}
