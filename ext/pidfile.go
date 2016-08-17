package ext

import (
	pidpkg "github.com/mkideal/pkg/pid"
)

// PidFile
type PidFile struct {
	filename string
}

func (pid PidFile) String() string {
	return pid.filename
}

func (pid *PidFile) Decode(s string) error {
	pid.filename = s
	return nil
}

func (pid *PidFile) New() error {
	return pidpkg.New(pid.filename)
}

func (pid PidFile) Remove() error {
	return pidpkg.Remove(pid.filename)
}
