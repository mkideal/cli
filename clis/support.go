////////////////////////////////////////////////////////////////////////////
// Program: support.go
// Purpose: cli boilerplate runtime support functions
// Authors: Tong Sun (c) 2015-2018, All rights reserved
////////////////////////////////////////////////////////////////////////////

package clis

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/gommon/color"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The OptsT type defines all the configurable options from cli.
type OptsT struct {
	Verbose int
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var (
	progname = "wireframe"
	Opts     OptsT
)

func Setup(p string, v int) {
	progname = p
	Opts.Verbose = v
}

// SUPPORT-FUNCTIONS
//==========================================================================
// support functions

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Basename returns the file name without extension.
func Basename(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}
	return s
}

// IsExist checks if the given file exist
func IsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func Warning(m string) {
	fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", progname, color.Yellow("Warning"), m)
}

func WarnOn(errCase string, e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "[%s] %s, %s: %v\n",
			color.White(progname), color.Yellow("Error"), errCase, e)
	}
}

// abortOn will quit on anticipated errors gracefully without stack trace
func AbortOn(errCase string, e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "[%s] %s, %s: %v\n",
			color.White(progname), color.Red("Error"), errCase, e)
		os.Exit(1)
	}
}

// verbose will print info to stderr according to the verbose level setting
func Verbose(levelSet int, format string, args ...interface{}) {
	if Opts.Verbose >= levelSet {
		fmt.Fprintf(os.Stderr, "[%s] ", color.White(progname))
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}
