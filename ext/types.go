package ext

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/jinzhu/now"
)

// Time wrap time.Time
type Time struct {
	time.Time
	isSet bool
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
	if s == "" {
		t.Time = time.Now()
		return nil
	}
	for _, format := range timeFormats {
		v, err := time.Parse(format, s)
		if err == nil {
			t.Time = v
			t.isSet = true
			return nil
		}
	}
	v, err := now.Parse(s)
	if err != nil {
		return err
	}
	t.Time = v
	t.isSet = true
	return nil
}

func (t Time) Encode() string {
	return t.Format(time.RFC3339Nano)
}

func (t Time) IsSet() bool {
	return t.isSet
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
	if f.data == nil {
		return ""
	}
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
