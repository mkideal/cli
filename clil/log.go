package clil

import (
	"log"
	"os"

	"github.com/comail/colog"
)

// NewLog creates new logger
func NewLog(level string) (*log.Logger, error) {

	lvl, err := colog.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	cl := colog.NewCoLog(os.Stderr, "", log.Lshortfile|log.Ldate|log.Ltime)
	cl.SetMinLevel(lvl)
	cl.SetDefaultLevel(lvl)
	lg := cl.NewLogger()
	return lg, nil
}
