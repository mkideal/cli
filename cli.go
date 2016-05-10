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
func Run(argv interface{}, fn CommandFunc, descs ...string) {
	desc := ""
	if len(descs) > 0 {
		desc = strings.Join(descs, "\n")
	}
	err := (&Command{
		Name:        os.Args[0],
		Desc:        desc,
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

// Parse parses args to object argv
func Parse(args []string, argv interface{}) error {
	clr := color.Color{}
	fset := parseArgv(args, argv, clr)
	return fset.err
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
		parseWithTypeValue(args, typ, val, flagSet, clr)
		return flagSet
	default:
		flagSet.err = errNotAPointer
		return flagSet
	}
}

func usage(v interface{}, clr color.Color, style UsageStyle) string {
	var (
		typ     = reflect.TypeOf(v)
		val     = reflect.ValueOf(v)
		flagSet = newFlagSet()
	)
	if typ.Kind() == reflect.Ptr &&
		reflect.Indirect(val).Type().Kind() == reflect.Struct {
		// initialize flagSet
		initFlagSet(typ, val, flagSet, clr)
		if flagSet.err != nil {
			return ""
		}
		return flagSlice(flagSet.flags).StringWithStyle(clr, style)
	}
	return ""
}

func initFlagSet(typ reflect.Type, val reflect.Value, flagSet *flagSet, clr color.Color) {
	var (
		typElem  = typ.Elem()
		valElem  = val.Elem()
		numField = valElem.NumField()
	)
	for i := 0; i < numField; i++ {
		var (
			typField     = typElem.Field(i)
			valField     = valElem.Field(i)
			tag, isEmpty = parseTag(typField.Name, typField.Tag)
		)
		if tag == nil {
			continue
		}
		// if `cli` tag is empty and the field is a struct
		if isEmpty && valField.Kind() == reflect.Struct {
			var (
				subObj   = valField.Addr().Interface()
				subType  = reflect.TypeOf(subObj)
				subValue = reflect.ValueOf(subObj)
			)
			initFlagSet(subType, subValue, flagSet, clr)
			if flagSet.err != nil {
				return
			}
			continue
		}
		fl, err := newFlag(typField, valField, tag, clr)
		if flagSet.err = err; err != nil {
			return
		}
		// ignored flag
		if fl == nil {
			continue
		}
		flagSet.flags = append(flagSet.flags, fl)

		// encode flag value
		value := ""
		if fl.assigned {
			if !valField.CanInterface() {
				flagSet.err = fmt.Errorf("field %s cannot interface", typField.Name)
				return
			}
			intf := valField.Interface()
			if encoder, ok := intf.(Encoder); ok {
				value = encoder.Encode()
			} else {
				value = fmt.Sprintf("%v", intf)
			}
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

func parseWithTypeValue(args []string, typ reflect.Type, val reflect.Value, flagSet *flagSet, clr color.Color) {
	initFlagSet(typ, val, flagSet, clr)
	if flagSet.err != nil {
		return
	}

	size := len(args)
	for i := 0; i < size; i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, dashOne) {
			// append a free argument
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

		if arg == dashOne {
			flagSet.err = fmt.Errorf("unexpected single dash")
			return
		}

		// terminate the flag parse while occur `--`
		if arg == dashTwo {
			flagSet.args = append(flagSet.args, args[i+1:]...)
			break
		}

		// split arg by "="(key=value)
		strs := []string{arg}
		index := strings.Index(arg, "=")
		if index >= 0 {
			strs = []string{arg[:index], arg[index+1:]}
		}

		arg = strs[0]
		fl, ok := flagSet.flagMap[arg]

		// found in flagMap
		if ok {
			retOffset := parseToFoundFlag(flagSet, fl, strs, arg, next, offset, clr)
			if flagSet.err != nil {
				return
			}
			i += retOffset
			continue
		}

		// not found in flagMap
		// it's an invalid flag if arg has prefix `--`
		if strings.HasPrefix(arg, dashTwo) {
			flagSet.err = fmt.Errorf("undefined flag %s", clr.Bold(arg))
			return
		}

		// try parse `-F<value>`
		if parseSiameseFlag(flagSet, arg[0:2], args[i][2:], clr) {
			continue
		} else if flagSet.err != nil {
			return
		}

		// other cases, find flag char by char
		arg = strings.TrimPrefix(arg, dashOne)
		parseFlagCharByChar(flagSet, arg, clr)
		if flagSet.err != nil {
			return
		}
		continue
	}

	for _, fl := range flagSet.flags {
		if fl.tag.isHelp && fl.getBool() {
			flagSet.dontValidate = true
			break
		}
	}
	if !flagSet.dontValidate {
		flagSet.readPrompt(os.Stdout, clr)
		if flagSet.err != nil {
			return
		}
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
	if buff.Len() > 0 && !flagSet.dontValidate {
		flagSet.err = fmt.Errorf(buff.String())
	}
}

func parseToFoundFlag(flagSet *flagSet, fl *flag, strs []string, arg, next string, offset int, clr color.Color) int {
	retOffset := 0
	l := len(strs)
	if l == 1 {
		if fl.isBoolean() {
			fl.v.SetBool(true)
		} else {
			retOffset = offset
			flagSet.err = fl.set(arg, next, clr)
		}
	} else if l == 2 {
		flagSet.err = fl.set(arg, strs[1], clr)
	} else {
		flagSet.err = fmt.Errorf("too many(%d) value", l)
	}
	if flagSet.err != nil {
		name := clr.Bold(fl.name())
		flagSet.err = fmt.Errorf("argument %s invalid: %v", name, flagSet.err)
		return retOffset
	}
	flagSet.values[arg] = []string{fmt.Sprintf("%v", fl.v.Interface())}
	return retOffset
}

func parseFlagCharByChar(flagSet *flagSet, arg string, clr color.Color) {
	// NOTE: each fold flag should be boolean
	chars := []byte(arg)
	for _, c := range chars {
		tmp := dashOne + string([]byte{c})
		fl, ok := flagSet.flagMap[tmp]
		if !ok {
			flagSet.err = fmt.Errorf("undefined flag %s", clr.Bold(tmp))
			return
		}

		if !fl.isBoolean() {
			flagSet.err = fmt.Errorf("each fold flag should be boolean, but %s not", clr.Bold(tmp))
			return
		}

		fl.v.SetBool(true)
		flagSet.values[tmp] = []string{"true"}
	}
}

func parseSiameseFlag(flagSet *flagSet, firstHalf, latterHalf string, clr color.Color) bool {
	// NOTE: fl must be not a boolean
	key, val := firstHalf, latterHalf
	if fl, ok := flagSet.flagMap[key]; ok && !fl.isBoolean() {
		if flagSet.err = fl.set(key, val, clr); flagSet.err != nil {
			return false
		}
		return true
	}
	return false
}
