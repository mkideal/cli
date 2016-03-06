package cli

import (
	"net/http"
	"strings"
)

const (
	// Command run error
	StatusRunError = 1
)

func (cmd *Command) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	debugf("path: %s", path)
	debugf("args: %q", args)

	if err := cmd.RunWithWriter(args, w); err != nil {
		w.WriteHeader(StatusRunError)
		w.Write([]byte(err.Error()))
	}
}
