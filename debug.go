package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

var paths = gopaths()

func gopaths() []string {
	paths := strings.Split(os.Getenv("GOPATH"), ";")
	for i, path := range paths {
		paths[i] = filepath.Join(path, "src") + "/"
	}
	if goroot := runtime.GOROOT(); goroot != "" {
		paths = append(paths, filepath.Join(goroot, "src")+"/")
	}
	return paths
}

func debugf(format string, args ...interface{}) {
	if !enableDebug {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	fileline := makeFileLine(file, line)
	fmt.Printf(Bold(fileline)+" "+format+"\n", args...)
}

func panicf(format string, args ...interface{}) {
	buff := bytes.NewBufferString("")
	buff.WriteString(Red(format, args...))
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
	fmt.Fprintf(os.Stderr, buff.String())
	os.Exit(999)
}

func makeFileLine(file string, line int) string {
	for _, path := range paths {
		if strings.HasPrefix(file, path) {
			file = strings.TrimPrefix(strings.TrimPrefix(file, path), "/")
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
