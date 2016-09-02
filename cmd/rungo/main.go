package main

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mkideal/cli"
	clix "github.com/mkideal/cli/ext"
	"github.com/mkideal/log"
	"github.com/mkideal/log/logger"
)

type argT struct {
	cli.Helper
	Level logger.Level `cli:"v" usage:"verbose level" dft:"info"`
	File  clix.File    `cli:"f,file" usage:"go script file"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		log.SetLevel(argv.Level)
		log.Trace("file content: \n%v", string(argv.File.Data()))

		filename := filepath.Join(os.TempDir(), randFilename())
		ioutil.WriteFile(filename, argv.File.Data(), 0644)

		cmd := exec.Command("go", "run", filename)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	})
}

const chars = "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890"

func randFilename() string {
	buf := new(bytes.Buffer)
	for i := 0; i < 32; i++ {
		b := chars[rand.Intn(len(chars))]
		buf.WriteByte(b)
	}
	return buf.String() + ".go"
}
