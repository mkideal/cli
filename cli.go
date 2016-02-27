package cli

import (
	"fmt"
	"reflect"
	"strings"
)

func Parse(args []string, v interface{}) *FlagSet {
	var (
		typ     = reflect.TypeOf(v)
		val     = reflect.ValueOf(v)
		flagSet = newFlagSet()
	)
	switch typ.Kind() {
	case reflect.Ptr:
		if reflect.Indirect(val).Type().Kind() != reflect.Struct {
			flagSet.Error = fmt.Errorf("object pointer does not indirect a struct")
			return flagSet
		}
		parse(args, typ, val, flagSet)
		return flagSet
	default:
		flagSet.Error = fmt.Errorf("type of object is not a pointer")
		return flagSet
	}
}

func Usage(v interface{}) string {
	var (
		typ     = reflect.TypeOf(v)
		val     = reflect.ValueOf(v)
		flagSet = newFlagSet()
	)
	if typ.Kind() == reflect.Ptr {
		if reflect.Indirect(val).Type().Kind() == reflect.Struct {
			initFlagSet(typ, val, flagSet)
			return flagSet.Usage
		}
	}
	return ""
}

func initFlagSet(typ reflect.Type, val reflect.Value, flagSet *FlagSet) {
	var (
		tm       = typ.Elem()
		vm       = val.Elem()
		fieldNum = vm.NumField()
		flags    = []*flag{}
	)
	for i := 0; i < fieldNum; i++ {
		tfield := tm.Field(i)
		vfield := vm.Field(i)
		fl, err := newFlag(tfield, vfield)
		if flagSet.Error = err; err != nil {
			return
		}
		// Ignored flag
		if fl == nil {
			continue
		}
		flags = append(flags, fl)
		value := ""
		if fl.assigned {
			value = fmt.Sprintf("%v", vfield.Interface())
		}

		names := append(fl.tag.shortNames, fl.tag.longNames...)
		for _, name := range names {
			if _, ok := flagSet.flags[name]; ok {
				flagSet.Error = fmt.Errorf("flag `%s` repeat", name)
				return
			}
			flagSet.flags[name] = fl
			if fl.assigned {
				flagSet.Values[name] = []string{value}
			}
		}
	}
	flagSet.Usage = flagSlice(flags).String()
}

func parse(args []string, typ reflect.Type, val reflect.Value, flagSet *FlagSet) {
	initFlagSet(typ, val, flagSet)
	if flagSet.Error != nil {
		return
	}

	size := len(args)
	for i := 0; i < size; i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, dashOne) {
			continue
		}
		values := []string{}
		for j := i + 1; j < size; j++ {
			if strings.HasPrefix(args[j], dashOne) {
				break
			}
			values = append(values, args[j])
		}
		i += len(values)

		strs := strings.Split(arg, "=")
		if strs == nil || len(strs) == 0 {
			continue
		}

		arg = strs[0]
		fl, ok := flagSet.flags[arg]
		if !ok {
			// If has prefix `--`
			if strings.HasPrefix(arg, dashTwo) {
				flagSet.Error = fmt.Errorf("unknown flag `%s`", arg)
				return
			}
			// Else find arg char by char
			chars := []byte(strings.TrimPrefix(arg, dashOne))
			for _, c := range chars {
				tmp := dashOne + string([]byte{c})
				if fl, ok := flagSet.flags[tmp]; !ok {
					flagSet.Error = fmt.Errorf("unknown flag `%s`", tmp)
					return
				} else {
					if flagSet.Error = fl.set(""); flagSet.Error != nil {
						return
					}
					if !fl.invalid {
						flagSet.Values[tmp] = []string{fmt.Sprintf("%v", fl.v.Interface())}
					}
				}
			}
			continue
		}

		values = append(strs[1:], values...)
		if len(values) == 0 {
			flagSet.Error = fl.set("")
		} else if len(values) == 1 {
			flagSet.Error = fl.set(values[0])
		} else {
			flagSet.Error = fmt.Errorf("too many(%d) value for flag `%s`", len(values), arg)
		}
		if flagSet.Error != nil {
			return
		}
		if !fl.invalid {
			flagSet.Values[arg] = []string{fmt.Sprintf("%v", fl.v.Interface())}
		}
	}

	for _, fl := range flagSet.flags {
		if !fl.assigned && fl.tag.required {
			flagSet.Error = fmt.Errorf("required argument `%s` missing", fl.name())
			return
		}
		if fl.assigned && fl.invalid {
			flagSet.Error = fmt.Errorf("assigned argument `%s` invalid", fl.name())
			return
		}
	}
}
