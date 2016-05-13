package cli

import (
	"encoding/json"
	"testing"

	"github.com/labstack/gommon/color"
)

type customT struct {
	K1 string
	K2 int
}

func (t *customT) Decode(s string) error {
	return json.Unmarshal([]byte(s), t)
}

func TestCustomType(t *testing.T) {
	type argT struct {
		D customT `cli:"d"`
	}
	v := new(argT)
	clr := color.Color{}
	flagSet := parseArgv([]string{`-d`, `{"k1": "string", "k2": 2}`}, v, clr)
	if flagSet.err != nil {
		t.Errorf("error: %v", flagSet.err)
		return
	}
	want := customT{K1: "string", K2: 2}
	if v.D != want {
		t.Errorf("want %v, got %v", want, v.D)
	}
}
