package cli

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type FlagSet struct {
	Error  error
	Usage  string
	Values url.Values

	flags map[string]*flag
}

func newFlagSet() *FlagSet {
	return &FlagSet{
		flags:  make(map[string]*flag),
		Values: url.Values(make(map[string][]string)),
	}
}

type flag struct {
	t reflect.StructField
	v reflect.Value

	assigned bool
	invalid  bool
	tag      cliTag
	typ      string
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
		if v, ok := getBool(s); ok {
			fl.v.SetBool(v)
		} else {
			fl.invalid = true
		}

	case reflect.String:
		fl.v.SetString(s)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, ok := getInt(s); ok {
			fl.v.SetInt(v)
		} else {
			fl.invalid = true
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, ok := getInt(s); ok {
			fl.v.SetUint(uint64(v))
		} else {
			fl.invalid = true
		}

	case reflect.Float32, reflect.Float64:
		if v, ok := getFloat(s); ok {
			fl.v.SetFloat(float64(v))
		} else {
			fl.invalid = true
		}
	default:
		return fmt.Errorf("invalid field type: %s", kind.String())
	}
	return nil
}

func getBool(s string) (bool, bool) {
	if s == "true" || s == "" {
		return true, true
	}
	if s == "false" || s == "none" {
		return false, true
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return false, false
	}
	return i != 0, true
}

func getInt(s string) (int64, bool) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}

func getFloat(s string) (float64, bool) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

type flagSlice []*flag

func (fs flagSlice) String() string {
	var (
		lenShort = 0
		lenLong  = 0
		sepSpace = len(sepName)
	)
	for _, fl := range fs {
		tag := fl.tag
		l := 0
		for _, shortName := range tag.shortNames {
			l += len(shortName) + sepSpace
		}
		if l > lenShort {
			lenShort = l
		}
		l = 0
		for _, longName := range tag.longNames {
			l += len(longName) + sepSpace
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
		l1, l2 := lenShort+sepSpace, lenLong+sepSpace
		if shortStr == "" {
			format = fmt.Sprintf("%%%ds%%-%ds%%s", l1, l2)
		} else if longStr == "" {
			format = fmt.Sprintf("%%%ds%s%%-%ds%%s", lenShort, strings.Repeat(" ", sepSpace), l2)
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
