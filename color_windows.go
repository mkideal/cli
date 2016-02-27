package cli

import (
	"fmt"
)

func gray(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func red(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
