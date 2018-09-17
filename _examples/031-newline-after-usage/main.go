package main

import (
	"github.com/mkideal/cli"
	"os"
)

type argT struct {
	cli.Helper
	Self      *argT       `json:"-" cli:"c,config" usage:"config\n" parser:"jsonfile" dft:"$__EXEC_FILENAME.json"`
	Host      string      `cli:"H,host" usage:"host addr" dft:"$HOST"`
	Port      int         `cli:"p,port" usage:"listening port\n"`
	Daemonize bool        `cli:"D,daemonize" usage:"daemonize the service"`
	Output    string      `cli:"o,output" usage:"The output file" dft:"some file"`
	Verbose   cli.Counter `cli:"v,verbose" usage:"Verbose mode (Multiple -v options increase the verbosity.)"`
}

func main() {
	cli.SetUsageStyle(cli.DenseNormalStyle)
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.JSONln(argv)
		return nil
	}))
}
