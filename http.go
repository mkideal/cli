package cli

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/labstack/gommon/color"
)

func (cmd *Command) RegisterHTTP(ctxs ...*Context) error {
	clr := color.Color{}
	clr.Disable()
	if len(ctxs) > 0 {
		clr = ctxs[0].color
	}
	if cmd.routersMap == nil {
		cmd.routersMap = make(map[string]string)
	}
	commands := []*Command{cmd}
	for len(commands) > 0 {
		c := commands[0]
		commands = commands[1:]
		if c.HTTPRouters != nil {
			for _, r := range c.HTTPRouters {
				if _, exists := c.routersMap[r]; exists {
					return throwRouterRepeat(clr.Yellow(r))
				}
				cmd.routersMap[r] = c.Path()
			}
		}
		if c.nochild() {
			continue
		}
		commands = append(commands, c.children...)
	}
	return nil
}

func (cmd *Command) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}

	var (
		path  = r.URL.Path
		found = false
	)
	if cmd.routersMap != nil {
		path, found = cmd.routersMap[path]
	}
	if !found {
		path = strings.TrimPrefix(r.URL.Path, "/")
		path = strings.TrimSuffix(path, "/")
	}

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
	Debugf("path: %s", path)
	Debugf("args: %q", args)

	buf := new(bytes.Buffer)
	statusCode := http.StatusOK
	if err := cmd.RunWith(args, buf, r.Method); err != nil {
		buf.Write([]byte(err.Error()))
		nativeError := err
		if werr, ok := err.(wrapError); ok {
			nativeError = werr.err
		}
		Debugf("error type: %s", TypeName(nativeError))
		switch nativeError.(type) {
		case commandNotFoundError:
			statusCode = http.StatusNotFound

		case methodNotAllowedError:
			statusCode = http.StatusMethodNotAllowed

		default:
			statusCode = http.StatusInternalServerError
		}
	}
	Debugf("resp: %s", buf.String())
	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}
