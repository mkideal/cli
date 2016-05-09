package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/now"
)

// Time wrap time.Time
type Time time.Time

func init() {
	now.TimeFormats = append(now.TimeFormats, time.ANSIC)
	now.TimeFormats = append(now.TimeFormats, time.UnixDate)
	now.TimeFormats = append(now.TimeFormats, time.RubyDate)
	now.TimeFormats = append(now.TimeFormats, time.RFC822)
	now.TimeFormats = append(now.TimeFormats, time.RFC822Z)
	now.TimeFormats = append(now.TimeFormats, time.RFC850)
	now.TimeFormats = append(now.TimeFormats, time.RFC1123)
	now.TimeFormats = append(now.TimeFormats, time.RFC1123Z)
	now.TimeFormats = append(now.TimeFormats, time.RFC3339)
	now.TimeFormats = append(now.TimeFormats, time.RFC3339Nano)
	now.TimeFormats = append(now.TimeFormats, time.Kitchen)
	now.TimeFormats = append(now.TimeFormats, time.Stamp)
	now.TimeFormats = append(now.TimeFormats, time.StampMilli)
	now.TimeFormats = append(now.TimeFormats, time.StampMicro)
	now.TimeFormats = append(now.TimeFormats, time.StampNano)
}

func (t *Time) Decode(s string) error {
	v, err := now.Parse(s)
	if err != nil {
		return err
	}
	*t = Time(v)
	return nil
}

func (t Time) Encode() string {
	return time.Time(t).Format(time.RFC3339Nano)
}

// Duration wrap time.Duration
type Duration time.Duration

func (d *Duration) Decode(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}

func (d Duration) Encode() string {
	return time.Duration(d).String()
}

// File reads data from file or stdin(if filename is empty)
type File struct {
	filename string
	data     []byte
}

func (rf File) Data() []byte {
	return rf.data
}

func (rf File) String() string {
	return string(rf.data)
}

func (rf *File) Decode(s string) error {
	var (
		data []byte
		err  error
	)
	if s == "" {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		data, err = ioutil.ReadFile(s)
	}
	if err != nil {
		return err
	}
	rf.data = data
	rf.filename = s
	return nil
}

func (rf File) Encode() string {
	return rf.filename
}

// PidFile
type PidFile struct {
	filename string
}

func (pid PidFile) String() string {
	return pid.filename
}

func (pid *PidFile) Decode(s string) error {
	pid.filename = s
	dir, _ := filepath.Split(pid.filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
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
