package ext

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
type Time struct {
	time.Time
}

var timeFormats = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

func (t *Time) Decode(s string) error {
	for _, format := range timeFormats {
		v, err := time.Parse(format, s)
		if err == nil {
			t.Time = v
			return nil
		}
	}
	v, err := now.Parse(s)
	if err != nil {
		return err
	}
	t.Time = v
	return nil
}

func (t Time) Encode() string {
	return t.Format(time.RFC3339Nano)
}

// Duration wrap time.Duration
type Duration struct {
	time.Duration
}

func (d *Duration) Decode(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = v
	return nil
}

func (d Duration) Encode() string {
	return d.Duration.String()
}

// File reads data from file or stdin(if filename is empty)
type File struct {
	filename string
	data     []byte
}

func (f File) Data() []byte {
	return f.data
}

func (f File) String() string {
	return string(f.data)
}

func (f *File) Decode(s string) error {
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
	f.data = data
	f.filename = s
	return nil
}

func (f File) Encode() string {
	return f.filename
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
	return nil
}

func (pid *PidFile) New() error {
	dir, _ := filepath.Split(pid.filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("PidFile: %v", err)
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
