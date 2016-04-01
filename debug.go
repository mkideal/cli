package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/labstack/gommon/color"
	"github.com/mattn/go-colorable"
)

var enableDebug = false

// EnableDebug open debug mode
func EnableDebug() {
	enableDebug = true
}

// DisableDebug close debug mode
func DisableDebug() {
	enableDebug = false
}

var gopaths = func() []string {
	paths := strings.Split(os.Getenv("GOPATH"), ";")
	for i, path := range paths {
		paths[i] = filepath.Join(path, "src") + "/"
	}
	if goroot := runtime.GOROOT(); goroot != "" {
		paths = append(paths, filepath.Join(goroot, "src")+"/")
	}
	return paths
}()

var debugOut, debugColor = func() (io.Writer, color.Color) {
	clr := color.Color{}
	out := colorable.NewColorableStdout()
	colorSwitch(&clr, out, os.Stdout.Fd())
	return out, clr
}()

func Debugf(format string, args ...interface{}) {
	if !enableDebug {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	fileline := makeFileLine(file, line)
	fmt.Fprintf(debugOut, "[DEBUG]"+debugColor.Bold(fileline)+" "+format+"\n", args...)
}

func Panicf(format string, args ...interface{}) {
	clr := color.Color{}
	out := colorable.NewColorableStderr()
	colorSwitch(&clr, out, os.Stderr.Fd())
	buff := bytes.NewBufferString("")
	buff.WriteString(clr.Red(fmt.Sprintf(format, args...)))
	buff.WriteString("\n\n[stack]\n")
	skip := 1
	for {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		skip++
		buff.WriteString(makeFileLine(file, line))
		buff.WriteString("\n")
	}
	fmt.Fprintf(out, buff.String())
	os.Exit(999)
}

func TypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func makeFileLine(file string, line int) string {
	for _, path := range gopaths {
		if strings.HasPrefix(file, path) {
			file = strings.TrimPrefix(strings.TrimPrefix(file, path), "/")
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
