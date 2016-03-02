package cli

import (
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
	if goroot := os.Getenv("GOROOT"); goroot != "" {
		paths = append(paths, filepath.Join(goroot, "src")+"/")
	}
	return paths
}

func debugf(format string, args ...interface{}) {
	if !enableDebug {
		return
	}
	fileline := getFileLine(2)
	fmt.Printf(Bold(fileline)+" "+format+"\n", args...)
}

func panicf(format string, args ...interface{}) {
	fileline := getFileLine(2)
	fmt.Fprintf(os.Stderr, Bold(fileline)+" "+Red(format, args...)+"\n")
	os.Exit(999)
}

func getFileLine(skip int) string {
	_, file, line, _ := runtime.Caller(skip)

	for _, path := range paths {
		if strings.HasPrefix(file, path) {
			file = strings.TrimPrefix(strings.TrimPrefix(file, path), "/")
			break
		}
	}
	return fmt.Sprintf("[%s:%d]", file, line)
}
