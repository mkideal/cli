package cli

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/labstack/gommon/color"
)

// Run runs a single command app
func Run(argv interface{}, fn CommandFunc) {
	err := (&Command{
		Name:        os.Args[0],
		Argv:        func() interface{} { return argv },
		CanSubRoute: true,
		Fn:          fn,
	}).Run(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
}

// Root registers forest for root and return root
func Root(root *Command, forest ...*CommandTree) *Command {
	root.RegisterTree(forest...)
	return root
}

// Tree creates a CommandTree
func Tree(cmd *Command, forest ...*CommandTree) *CommandTree {
	return &CommandTree{
		command: cmd,
		forest:  forest,
	}
}

//------------------
// Implements parse
//------------------

func parseArgv(args []string, argv interface{}, clr color.Color) *flagSet {
	var (
		typ     = reflect.TypeOf(argv)
		val     = reflect.ValueOf(argv)
		flagSet = newFlagSet()
	)
	switch typ.Kind() {
	case reflect.Ptr:
		if reflect.Indirect(val).Type().Kind() != reflect.Struct {
			flagSet.err = errNotPointToStruct
			return flagSet
		}
		parse(args, typ, val, flagSet, clr)
		return flagSet
	default:
		flagSet.err = errNotAPointer
		return flagSet
	}
}

func usage(v interface{}, clr color.Color) string {
	var (
		typ     = reflect.TypeOf(v)
		val     = reflect.ValueOf(v)
		flagSet = newFlagSet()
	)
	if typ.Kind() == reflect.Ptr {
		if reflect.Indirect(val).Type().Kind() == reflect.Struct {
			initFlagSet(typ, val, flagSet, clr)
			if flagSet.err != nil {
				return ""
			}
			return flagSlice(flagSet.flags).String(clr)
		}
	}
	return ""
}

func initFlagSet(typ reflect.Type, val reflect.Value, flagSet *flagSet, clr color.Color) {
	var (
		tm       = typ.Elem()
		vm       = val.Elem()
		fieldNum = vm.NumField()
	)
	for i := 0; i < fieldNum; i++ {
		tfield := tm.Field(i)
		vfield := vm.Field(i)
		tag, isEmpty := parseTag(tfield.Name, tfield.Tag)
		if tag == nil {
			continue
		}
		// if `cli` tag is empty and the field is a struct
		if isEmpty && vfield.Kind() == reflect.Struct {
			subObj := vfield.Addr().Interface()
			initFlagSet(reflect.TypeOf(subObj), reflect.ValueOf(subObj), flagSet, clr)
			if flagSet.err != nil {
				return
			}
			continue
		}
		fl, err := newFlag(tfield, vfield, tag, clr)
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
		for i, name := range names {
			if _, ok := flagSet.flagMap[name]; ok {
				flagSet.err = fmt.Errorf("flag %s repeat", clr.Bold(name))
				return
			}
			flagSet.flagMap[name] = fl
			if fl.assigned && i == 0 {
				flagSet.values[name] = []string{value}
			}
		}
	}
}

func parse(args []string, typ reflect.Type, val reflect.Value, flagSet *flagSet, clr color.Color) {
	initFlagSet(typ, val, flagSet, clr)
	if flagSet.err != nil {
		return
	}

	size := len(args)
	for i := 0; i < size; i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, dashOne) {
			// append a freedom argument
			flagSet.args = append(flagSet.args, arg)
			continue
		}

		var (
			next   = ""
			offset = 0
		)
		if i+1 < size {
			if !strings.HasPrefix(args[i+1], dashOne) {
				next = args[i+1]
				offset = 1
			}
		}

		// terminate the flag parse while occur `--`
		if arg == "--" {
			flagSet.args = append(flagSet.args, args[i+1:]...)
			break
		}

		// split arg by "="(key=value)
		strs := strings.Split(arg, "=")
		if strs == nil || len(strs) == 0 {
			continue
		}

		arg = strs[0]
		fl, ok := flagSet.flagMap[arg]
		if ok {
			l := len(strs)
			if l == 1 {
				if fl.isBoolean() {
					fl.v.SetBool(true)
				} else {
					i += offset
					flagSet.err = fl.set(arg, next, clr)
				}
			} else if l == 2 {
				flagSet.err = fl.set(arg, strs[1], clr)
			} else {
				flagSet.err = fmt.Errorf("too many(%d) value", l, clr.Bold(arg))
			}
			if flagSet.err != nil {
				name := clr.Bold(fl.name())
				flagSet.err = fmt.Errorf("argument %s invalid: %v", name, flagSet.err)
				return
			}
			flagSet.values[arg] = []string{fmt.Sprintf("%v", fl.v.Interface())}
			continue
		}

		// if arg has prefix `--`, then it's an invalid flag
		if strings.HasPrefix(arg, dashTwo) {
			flagSet.err = fmt.Errorf("undefined flag %s", clr.Bold(arg))
			return
		}

		arg = strings.TrimPrefix(arg, dashOne)

		// try parse `-F<value>`
		// NOTE: fl must be not a boolean
		key, val := dashOne+arg[0:1], arg[1:]
		if fl, ok := flagSet.flagMap[key]; ok && !fl.isBoolean() {
			if flagSet.err = fl.set(key, val, clr); flagSet.err != nil {
				return
			}
			continue
		}

		// other cases, find flag char by char
		// NOTE: every fold flag should be boolean
		chars := []byte(arg)
		for _, c := range chars {
			tmp := dashOne + string([]byte{c})
			fl, ok := flagSet.flagMap[tmp]
			if !ok {
				flagSet.err = fmt.Errorf("undefined flag %s", clr.Bold(tmp))
				return
			}

			if !fl.isBoolean() {
				flagSet.err = fmt.Errorf("every fold flag should be boolean, but %s not", clr.Bold(tmp))
				return
			}

			fl.v.SetBool(true)
			flagSet.values[tmp] = []string{"true"}
		}
		continue
	}

	buff := bytes.NewBufferString("")
	for _, fl := range flagSet.flags {
		if !fl.assigned && fl.tag.required {
			if buff.Len() > 0 {
				buff.WriteByte('\n')
			}
			fmt.Fprintf(buff, "required argument %s missing", clr.Bold(fl.name()))
		}
	}
	for _, fl := range flagSet.flags {
		if fl.tag.isHelp && fl.getBool() {
			flagSet.dontValidate = true
			break
		}
	}
	if buff.Len() > 0 && !flagSet.dontValidate {
		flagSet.err = fmt.Errorf(buff.String())
	}
}
