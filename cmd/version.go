package main

import (
	"fmt"

	"github.com/metal-stack/v"
	"github.com/urfave/cli/v2"
)

var versionCmd = &cli.Command{
	Name:  "version",
	Usage: "show current version",
	Action: func(*cli.Context) error {
		fmt.Printf("%s\n", v.V.String())
		return nil
	},
}
