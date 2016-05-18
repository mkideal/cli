package ext

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/mkideal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeTime(t *testing.T) {
	type argT struct {
		When Time `cli:"w" dft:"2016-01-02"`
	}
	local := time.Local
	for _, tt := range []struct {
		src string
		t   time.Time
	}{
		{"-w2016-01-02T15:04:05+00:00", time.Date(2016, 1, 2, 15, 4, 5, 0, time.UTC)},
		{"-w2016-01-02", time.Date(2016, 1, 2, 0, 0, 0, 0, local)},
		{"", time.Date(2016, 1, 2, 0, 0, 0, 0, local)},
	} {
		args := []string{}
		if tt.src != "" {
			args = append(args, tt.src)
		}
		argv := new(argT)
		assert.Nil(t, cli.Parse(args, argv))
		assert.Equal(t, tt.t.Unix(), argv.When.Unix())
	}
}

func TestTypeDuration(t *testing.T) {
	type argT struct {
		Long Duration `cli:"d" dft:"1s"`
	}
	for _, tt := range []struct {
		args []string
		long time.Duration
	}{
		{[]string{}, time.Second},
		{[]string{"-d10s"}, time.Second * 10},
		{[]string{"-d10ms"}, time.Millisecond * 10},
	} {
		argv := new(argT)
		assert.Nil(t, cli.Parse(tt.args, argv))
		assert.Equal(t, tt.long, argv.Long.Duration)
	}
}

func TestTypeFile(t *testing.T) {
	type argT struct {
		File File `cli:"f"`
	}
	filename := "yXLLBhNHkv9VdAarIF87"
	content := "hello,world"
	require.Nil(t, ioutil.WriteFile(filename, []byte(content), 0644))
	defer os.Remove(filename)
	argv := new(argT)
	assert.Nil(t, cli.Parse([]string{"-f", filename}, argv))
	assert.Equal(t, content, argv.File.String())
}
