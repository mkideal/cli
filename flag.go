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

	flags map[string]*Flag
}

func newFlagSet() *FlagSet {
	return &FlagSet{
		flags:  make(map[string]*Flag),
		Values: url.Values(make(map[string][]string)),
	}
}

type Flag struct {
	t reflect.StructField
	v reflect.Value

	assigned bool
	invalid  bool
	tag      cliTag
	typ      string
}

func newFlag(t reflect.StructField, v reflect.Value) (flag *Flag, err error) {
	flag = &Flag{t: t, v: v}
	tag, err := parseTag(t.Name, t.Tag)
	if tag == nil {
		return nil, nil
	}
	flag.tag = *tag
	err = flag.init()
	flag.typ = t.Type.Kind().String()
	return
}

func (flag *Flag) init() error {
	if flag.tag.defaultValue != "" {
		return flag.set(flag.tag.defaultValue)
	}
	return nil
}

func (flag *Flag) name() string {
	if len(flag.tag.longNames) > 0 {
		return flag.tag.longNames[0]
	}
	if len(flag.tag.shortNames) > 0 {
		return flag.tag.shortNames[0]
	}
	return ""
}

func (flag *Flag) set(s string) error {
	kind := flag.t.Type.Kind()
	flag.assigned = true
	switch kind {
	case reflect.Bool:
		if v, ok := getBool(s); ok {
			flag.v.SetBool(v)
		} else {
			flag.invalid = true
		}

	case reflect.String:
		flag.v.SetString(s)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, ok := getInt(s); ok {
			flag.v.SetInt(v)
		} else {
			flag.invalid = true
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, ok := getInt(s); ok {
			flag.v.SetUint(uint64(v))
		} else {
			flag.invalid = true
		}

	case reflect.Float32, reflect.Float64:
		if v, ok := getFloat(s); ok {
			flag.v.SetFloat(float64(v))
		} else {
			flag.invalid = true
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

type flagSlice []*Flag

func (fs flagSlice) String() string {
	var (
		lenShort = 0
		lenLong  = 0
		lenType  = 0
		sepSpace = len(sepName)
	)
	for _, flag := range fs {
		tag := flag.tag
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
		l = len(flag.typ) + sepColSpace
		if l > lenType {
			lenType = l
		}
	}

	buff := bytes.NewBufferString("")
	for _, flag := range fs {
		tag := flag.tag
		shortStr := strings.Join(tag.shortNames, sepName)
		longStr := strings.Join(tag.longNames, sepName)
		typeStr := ""
		if tag.required {
			typeStr = fmt.Sprintf("(%s*)", flag.typ)
		} else {
			typeStr = fmt.Sprintf("(%s)", flag.typ)
		}
		format := ""
		l1, l2, l3 := lenShort+sepSpace, lenLong+sepSpace, lenType+sepColSpace
		if shortStr == "" {
			format = fmt.Sprintf("%%%ds%%-%ds%%-%ds%%s", l1, l2, l3)
		} else if longStr == "" {
			format = fmt.Sprintf("%%%ds%s%%-%ds%%-%ds%%s", lenShort, strings.Repeat(" ", sepSpace), l2, l3)
		} else {
			format = fmt.Sprintf("%%%ds%s%%-%ds%%-%ds%%s", lenShort, sepName, l2, l3)
		}
		//fmt.Printf("format: %v\n", format) // FIXME: remove me
		usage := tag.usage
		if tag.defaultValue != "" {
			usage = fmt.Sprintf("%s[default=%s]", usage, tag.defaultValue)
		}
		buff.WriteString(fmt.Sprintf(format+"\n", shortStr, longStr, typeStr, usage))
	}
	return buff.String()
}
