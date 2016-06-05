package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

//TODO:func TestContextGetArgvList(t *testing.T) {
//}

func TestContextMisc(t *testing.T) {
	type argT struct {
		Hello string `cli:"hello"`
		Age   int    `cli:"age" dft:"10"`
	}
	root := &Command{Name: "root"}
	parent := &Command{Name: "parent", Fn: donothing}
	cmd := &Command{
		Name: "cmd",
		Argv: func() interface{} { return new(argT) },
		Fn: func(ctx *Context) error {
			argv := ctx.Argv().(*argT)
			assert.Equal(t, ctx.Args(), []string{"a", "b", "c"})
			assert.Equal(t, argv.Hello, "world")
			assert.Equal(t, argv.Age, 10)
			assert.Equal(t, ctx.NativeArgs(), []string{"--hello=world", "a", "b", "c"})
			assert.Equal(t, ctx.Command().Name, "cmd")
			assert.Equal(t, ctx.IsSet("--hello"), true)
			assert.Equal(t, ctx.IsSet("--age"), false)
			assert.Equal(t, ctx.NArg(), 3)
			assert.Equal(t, ctx.Path(), "parent cmd")
			assert.Equal(t, ctx.Router(), []string{"parent", "cmd"})
			return nil
		},
	}
	root.Register(parent)
	parent.Register(cmd)
	assert.Nil(t, root.RunWith([]string{"parent", "cmd", "--hello=world", "a", "b", "c"}, nil, nil))
}
