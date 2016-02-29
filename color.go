// +build !windows

package cli

import (
	"fmt"
)

func gray(format string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[90m"+format+"\x1b[0m", args...)
}

func red(format string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[31m"+format+"\x1b[0m", args...)
}

func yellow(format string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[33m"+format+"\x1b[0m", args...)
}

func bold(format string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[;1m"+format+"\x1b[0m", args...)
}
