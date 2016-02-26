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
		flags    = []*Flag{}
	)
	for i := 0; i < fieldNum; i++ {
		tfield := tm.Field(i)
		vfield := vm.Field(i)
		flag, err := newFlag(tfield, vfield)
		if flagSet.Error = err; err != nil {
			return
		}
		// Ignored flag
		if flag == nil {
			continue
		}
		flags = append(flags, flag)
		value := ""
		if flag.assigned {
			value = fmt.Sprintf("%v", vfield.Interface())
		}

		names := append(flag.tag.shortNames, flag.tag.longNames...)
		for _, name := range names {
			if _, ok := flagSet.flags[name]; ok {
				flagSet.Error = fmt.Errorf("flag name `%s` repeat", name)
				return
			}
			flagSet.flags[name] = flag
			if flag.assigned {
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
			if strings.HasPrefix(arg, dashOne) {
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
		flag, ok := flagSet.flags[arg]
		if !ok {
			// If has prefix `--`
			if strings.HasPrefix(arg, dashTwo) {
				flagSet.Error = fmt.Errorf("unknown flag `%s`", arg)
				return
			}
			// Else find arg char by char
			chars := []byte(strings.TrimPrefix(arg, dashOne))
			for _, c := range chars {
				tmp := string([]byte{c})
				if flag, ok := flagSet.flags[tmp]; !ok {
					flagSet.Error = fmt.Errorf("unknown flag `%s`", tmp)
					return
				} else {
					if flagSet.Error = flag.set(""); flagSet.Error != nil {
						return
					}
					flagSet.Values[tmp] = []string{fmt.Sprintf("%v", flag.v.Interface())}
				}
			}
		}

		values = append(strs[1:], values...)
		if len(values) == 0 {
			flagSet.Error = flag.set("")
		} else if len(values) == 1 {
			flagSet.Error = flag.set(values[0])
		} else {
			flagSet.Error = fmt.Errorf("too many(%d) value for flag `%s`", len(values), arg)
		}
		if flagSet.Error != nil {
			return
		}
		flagSet.Values[arg] = []string{fmt.Sprintf("%v", flag.v.Interface())}
	}

	for _, flag := range flagSet.flags {
		if !flag.assigned && flag.tag.required {
			flagSet.Error = fmt.Errorf("required argument `%s` missing", flag.name())
			return
		}
		if flag.assigned && flag.invalid {
			flagSet.Error = fmt.Errorf("assign argument `%s` invalid", flag.name())
			return
		}
	}
}
