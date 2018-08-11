package cli

import (
	"net/url"
	"testing"

	"github.com/labstack/gommon/color"
)

func TestJSONParser(t *testing.T) {
	type T struct {
		A string
		B int
	}
	type argT struct {
		Value T `cli:"t" parser:"json"`
	}

	v := new(argT)
	clr := color.Color{}
	flagSet := parseArgv([]string{`-t`, `{"a": "string", "b": 2}`}, v, clr)
	if flagSet.err != nil {
		t.Errorf("error: %v", flagSet.err)
		return
	}
	want := T{A: "string", B: 2}
	if v.Value != want {
		t.Errorf("want %v, got %v", want, v.Value)
	}
}

func TestURLParser(t *testing.T) {
	type argT struct {
		Addr url.URL `cli:"u,url" parser:"url"`
	}

	v := new(argT)
	clr := color.Color{}
	flagSet := parseArgv([]string{`-u`, `https://www.google.com`}, v, clr)
	if flagSet.err != nil {
		t.Errorf("error: %v", flagSet.err)
		return
	}
	if v.Addr.Scheme != "https" {
		t.Errorf("schema want %v, got %v", "https", v.Addr.Scheme)
	}
	if v.Addr.Host != "www.google.com" {
		t.Errorf("host want %v, got %v", "www.google.com", v.Addr.Host)
	}
}
