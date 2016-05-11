package cli_test

import (
	"github.com/mkideal/cli"
)

func ExampleCommand() {
	root := &cli.Command{
		Name: "app",
	}

	type childT struct {
		S string `cli:"s" usage:"string flag"`
		B bool   `cli:"b" usage:"boolean flag"`
	}
	root.Register(&cli.Command{
		Name:        "child",
		Aliases:     []string{"sub"},
		Desc:        "child command",
		Text:        "detailed description for command",
		Argv:        func() interface{} { return new(childT) },
		CanSubRoute: true,
		NoHook:      true,
		NoHTTP:      true,
		NeedArgs:    true,
		HTTPRouters: []string{"/v1/child", "/v2/child"},
		HTTPMethods: []string{"GET", "POST"},

		OnRootPrepareError: func(err error) error {
			return err
		},
		OnBefore: func(ctx *cli.Context) error {
			ctx.String("OnBefore\n")
			return nil
		},
		OnAfter: func(ctx *cli.Context) error {
			ctx.String("OnAfter\n")
			return nil
		},
		OnRootBefore: func(ctx *cli.Context) error {
			ctx.String("OnRootBefore\n")
			return nil
		},
		OnRootAfter: func(ctx *cli.Context) error {
			ctx.String("OnRootAfter\n")
			return nil
		},

		Fn: func(ctx *cli.Context) error {
			return nil
		},
	})
}
