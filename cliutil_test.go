package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpCommand(t *testing.T) {
	w := bytes.NewBufferString("")
	root := &Command{Name: "root"}
	help := HelpCommand("help command")
	root.Register(help)
	assert.Nil(t, root.RunWith([]string{"help"}, w, nil))
	assert.Equal(t, w.String(), "Commands:\n\n  help   help command\n")
	assert.Error(t, root.RunWith([]string{"help", "not-found"}, nil, nil))
}
