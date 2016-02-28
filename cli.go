package cli

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Run runs a single command app
func Run(argv interface{}, fn CommandFunc) {
	err := Command{
		Name: os.Args[0],
		Argv: func() interface{} { return argv },
		Fn:   fn,
	}.Run(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
}

func parseArgv(args []string, argv interface{}) *flagSet {
	var (
		typ     = reflect.TypeOf(argv)
		val     = reflect.ValueOf(argv)
		flagSet = newFlagSet()
	)
	switch typ.Kind() {
	case reflect.Ptr:
		if reflect.Indirect(val).Type().Kind() != reflect.Struct {
			flagSet.err = fmt.Errorf("argv does not indirect a struct")
			return flagSet
		}
		parse(args, typ, val, flagSet)
		return flagSet
	default:
		flagSet.err = fmt.Errorf("argv is not a pointer")
		return flagSet
	}
}

func usage(v interface{}) string {
	var (
		typ     = reflect.TypeOf(v)
		val     = reflect.ValueOf(v)
		flagSet = newFlagSet()
	)
	if typ.Kind() == reflect.Ptr {
		if reflect.Indirect(val).Type().Kind() == reflect.Struct {
			initFlagSet(typ, val, flagSet)
			return flagSlice(flagSet.flags).String()
		}
	}
	return ""
}

func initFlagSet(typ reflect.Type, val reflect.Value, flagSet *flagSet) {
	var (
		tm       = typ.Elem()
		vm       = val.Elem()
		fieldNum = vm.NumField()
	)
	for i := 0; i < fieldNum; i++ {
		tfield := tm.Field(i)
		vfield := vm.Field(i)
		fl, err := newFlag(tfield, vfield)
		if flagSet.err = err; err != nil {
			return
		}
		// Ignored flag
		if fl == nil {
			continue
		}
		flagSet.flags = append(flagSet.flags, fl)
		value := ""
		if fl.assigned {
			value = fmt.Sprintf("%v", vfield.Interface())
		}

		names := append(fl.tag.shortNames, fl.tag.longNames...)
		for _, name := range names {
			if _, ok := flagSet.flagMap[name]; ok {
				flagSet.err = fmt.Errorf("flag `%s` repeat", name)
				return
			}
			flagSet.flagMap[name] = fl
			if fl.assigned {
				flagSet.values[name] = []string{value}
			}
		}
	}
}

func parse(args []string, typ reflect.Type, val reflect.Value, flagSet *flagSet) {
	initFlagSet(typ, val, flagSet)
	if flagSet.err != nil {
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
		fl, ok := flagSet.flagMap[arg]
		if !ok {
			// If has prefix `--`
			if strings.HasPrefix(arg, dashTwo) {
				flagSet.err = fmt.Errorf("undefined flag `%s`", arg)
				return
			}
			// Else find arg char by char
			chars := []byte(strings.TrimPrefix(arg, dashOne))
			for _, c := range chars {
				tmp := dashOne + string([]byte{c})
				fl, ok := flagSet.flagMap[tmp]
				if !ok {
					flagSet.err = fmt.Errorf("undefined flag `%s`", tmp)
					return
				}

				if flagSet.err = fl.set(""); flagSet.err != nil {
					return
				}
				if fl.err == nil {
					flagSet.values[tmp] = []string{fmt.Sprintf("%v", fl.v.Interface())}
				}

			}
			continue
		}

		values = append(strs[1:], values...)
		if len(values) == 0 {
			flagSet.err = fl.set("")
		} else if len(values) == 1 {
			flagSet.err = fl.set(values[0])
		} else {
			flagSet.err = fmt.Errorf("too many(%d) value for flag `%s`", len(values), arg)
		}
		if flagSet.err != nil {
			return
		}
		if fl.err == nil {
			flagSet.values[arg] = []string{fmt.Sprintf("%v", fl.v.Interface())}
		}
	}

	buff := bytes.NewBufferString("")
	for _, fl := range flagSet.flags {
		if !fl.assigned && fl.tag.required {
			if buff.Len() > 0 {
				buff.WriteByte('\n')
			}
			fmt.Fprintf(buff, "%s required argument `%s` missing", red("ERR!"), fl.name())
		}
		if fl.assigned && fl.err != nil {
			if buff.Len() > 0 {
				buff.WriteByte('\n')
			}
			fmt.Fprintf(buff, "%s assigned argument `%s` invalid: %v", red("ERR!"), fl.name(), fl.err)
		}
	}
	if buff.Len() > 0 {
		flagSet.err = fmt.Errorf(buff.String())
	}
}
