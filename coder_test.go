package cli

import (
	"encoding/json"
	"testing"

	"github.com/labstack/gommon/color"
	"github.com/stretchr/testify/assert"
)

type customT struct {
	K1 string
	K2 int
}

func (t *customT) Decode(s string) error {
	return json.Unmarshal([]byte(s), t)
}

func TestDecoder(t *testing.T) {
	type argT struct {
		D customT `cli:"d"`
	}
	v := new(argT)
	clr := color.Color{}
	flagSet := parseArgv([]string{`-d`, `{"k1": "string", "k2": 2}`}, v, clr)
	assert.Nil(t, flagSet.err)
	assert.Equal(t, v.D, customT{K1: "string", K2: 2})
}
