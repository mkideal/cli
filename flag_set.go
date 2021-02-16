package cli

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/labstack/gommon/color"
)

type flagSet struct {
	err    error
	values url.Values
	args   []string

	flagMap   map[string]*flag
	flagSlice []*flag

	hasForce bool
}

func newFlagSet() *flagSet {
	return &flagSet{
		flagMap:   make(map[string]*flag),
		flagSlice: []*flag{},
		values:    url.Values(make(map[string][]string)),
		args:      make([]string, 0),
	}
}

func (fs *flagSet) readPrompt(w io.Writer, clr color.Color) {
	for _, fl := range fs.flagSlice {
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
			data, fs.err = password(prefix)
			if fs.err == nil && data != "" {
				fl.setWithNoDelay("", data, clr)
			}
		} else if fl.isBoolean() {
			yes, fs.err = ask(prefix, false)
			if fs.err == nil {
				fl.setWithNoDelay("", fmt.Sprintf("%v", yes), clr)
			}
		} else if fl.tag.dft != "" {
			data, fs.err = promptDefault(prefix, fl.tag.dft)
			if fs.err == nil {
				fl.setWithNoDelay("", data, clr)
			}
		} else {
			data, fs.err = prompt(prefix, fl.tag.isRequired)
			if fs.err == nil {
				fl.setWithNoDelay("", data, clr)
			}
		}
		if fs.err != nil {
			return
		}
	}
}

func (fs *flagSet) readEditor(clr color.Color) {
	editor, editorErr := getEditor()
	for _, fl := range fs.flagSlice {
		if fl.isAssigned || !fl.tag.isEdit {
			continue
		}
		if editorErr != nil {
			fs.err = editorErr
			return
		}
		filename := fl.tag.editFile
		if filename == "" {
			filename = randomFilename()
		}
		data, err := launchEditorWithFilename(editor, filename)
		if fs.err = err; err != nil {
			return
		}
		if fs.err = fl.setWithNoDelay("", string(data), clr); fs.err != nil {
			return
		}
	}
}

// UsageStyle is style of usage
type UsageStyle int32

const (
	// NormalStyle : left-right
	NormalStyle UsageStyle = iota
	// DenseNormalStyle : left-right, too
	DenseNormalStyle
	// ManualStyle : up-down
	ManualStyle
	// DenseManualStyle : up-down, too
	DenseManualStyle
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
		if defaultStyle == NormalStyle && tag.dft != "" {
			lenDft = len(tag.dft) + 3 // 3=len("[=]")
			l += lenDft
		}
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
		spaceSize, lenDft := lenNameAndDefaultAndLong, 0

		if tag.dft != "" {
			defaultStr = fmt.Sprintf("[=%s]", tag.dft)
			lenDft = len(defaultStr)
			defaultStr = clr.Grey(defaultStr)
		}
		if tag.name != "" {
			nameStr = "=" + tag.name
		}
		if tag.isRequired {
			usagePrefix = clr.Red("*")
		}
		usage := usagePrefix + tag.usage
		lastNotNewLineIndex := len(usage) - 1
		for i := len(usage) - 1; i >= 0; i-- {
			if usage[i] != '\n' {
				lastNotNewLineIndex = i
				break
			}
		}

		// move defaultStr to the end when in DenseNormalStyle
		if defaultStyle == DenseNormalStyle {
			usage = usage[:lastNotNewLineIndex+1] + " " + defaultStr + usage[lastNotNewLineIndex+1:]
			defaultStr = ""
			lenDft = 0
		}

		spaceSize -= len(nameStr) + lenDft + len(longStr)

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
	if style != ManualStyle && style != DenseManualStyle {
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
		if fl.tag.dft != "" {
			buf.WriteString(clr.Grey(fmt.Sprintf("[=%s]", fl.tag.dft)))
		}
		buf.WriteString("\n")
		buf.WriteString(linePrefix)
		buf.WriteString("    ")
		if fl.tag.isRequired {
			buf.WriteString(clr.Red("*"))
		}
		buf.WriteString(fl.tag.usage)
		if style != DenseManualStyle {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}
