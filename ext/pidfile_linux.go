package ext

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	dir, _ := filepath.Split(pid.filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("PidFile: %v", err)
		}
	}
	if content, err := ioutil.ReadFile(pid.filename); err == nil {
		pidStr := strings.TrimSpace(string(content))
		if pidStr == "" {
			return nil
		}
		if _, err := os.Stat(filepath.Join("/proc", pidStr)); err == nil {
			return fmt.Errorf("pid file found, ensoure %s is not running", os.Args[0])
		}
	}
	if err := ioutil.WriteFile(pid.filename, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		return err
	}
	return nil
}

func (pid PidFile) Remove() error {
	if pid.filename != "" {
		return os.Remove(pid.filename)
	}
	return nil
}
