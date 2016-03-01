// +build !windows

package cli

import (
	"fmt"
)

func Gray(format string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[90m"+format+"\x1b[0m", args...)
}

func Red(format string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[31m"+format+"\x1b[0m", args...)
}

func Yellow(format string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[33m"+format+"\x1b[0m", args...)
}

func Bold(format string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[;1m"+format+"\x1b[0m", args...)
}
