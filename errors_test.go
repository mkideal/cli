package cli

import (
	"testing"

	"github.com/labstack/gommon/color"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	assert.Equal(t, ExitError.Error(), "exit")
	assert.Equal(t, throwCommandNotFound("cmd").Error(), "command cmd not found")
	assert.Equal(t, throwMethodNotAllowed("POST").Error(), "method POST not allowed")
	assert.Equal(t, throwRouterRepeat("R").Error(), "router R repeat")
	clr := color.Color{}
	clr.Disable()
	assert.Equal(t, wrapErr(throwCommandNotFound("cmd"), "_end", clr).Error(), `ERR! command cmd not found_end`)

	assert.Equal(t, argvError{isEmpty: true}.Error(), "argv list is empty")
	assert.Equal(t, argvError{isOutOfRange: true}.Error(), "argv list out of range")
	assert.Equal(t, argvError{ith: 1, msg: "ERROR MSG"}.Error(), "1th argv: ERROR MSG")
}
