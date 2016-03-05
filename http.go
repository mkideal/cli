package cli

import (
	"bytes"
	"net/http"
	"strings"
)

const (
	// Command run error
	StatusRunError = 1
)

func (ctx *Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	path = strings.TrimSuffix(path, "/")
	router := strings.Split(path, "/")
	args := make([]string, 0, len(r.Form)*2+len(router))
	for _, r := range router {
		args = append(args, r)
	}
	for key, values := range r.Form {
		if len(key) == 0 || len(values) == 0 {
			continue
		}
		if !strings.HasPrefix(key, dashOne) {
			if len(key) == 1 {
				key = dashOne + key
			} else {
				key = dashTwo + key
			}
		}
		args = append(args, key, values[len(values)-1])
	}
	buff := bytes.NewBufferString("")
	if err := ctx.Command().Root().RunWithWriter(args, buff); err != nil {
		w.WriteHeader(StatusRunError)
		buff.WriteString(err.Error())
	}
	debugf("path: %s", path)
	debugf("args: %q", args)
	debugf("resp: %s", buff.String())
	w.Write(buff.Bytes())
}
