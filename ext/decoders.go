package ext

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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

// CSV reads one csv record
type CSVRecord struct {
	raw []string
}

func (d *CSVRecord) Decode(s string) error {
	reader := csv.NewReader(strings.NewReader(s))
	record, err := reader.Read()
	if err != nil {
		return err
	}
	d.raw = record
	return nil
}

func (d CSVRecord) Strings() []string {
	return d.raw
}

func (d CSVRecord) Bools() ([]bool, error) {
	ret := make([]bool, len(d.raw))
	for _, s := range d.raw {
		s = strings.ToLower(s)
		if s == "y" || s == "yes" || s == "true" {
			ret = append(ret, true)
		} else if s == "n" || s == "no" || s == "false" {
			ret = append(ret, false)
		} else {
			v, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("parse %s to bollean fail", s)
			}
			ret = append(ret, v != 0)
		}
	}
	return ret, nil
}

func (d CSVRecord) Ints() ([]int64, error) {
	ret := make([]int64, len(d.raw))
	for _, s := range d.raw {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func (d CSVRecord) Uints() ([]uint64, error) {
	ret := make([]uint64, len(d.raw))
	for _, s := range d.raw {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func (d CSVRecord) Floats() ([]float64, error) {
	ret := make([]float64, len(d.raw))
	for _, s := range d.raw {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}
