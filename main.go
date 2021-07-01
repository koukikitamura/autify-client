package main

import (
	"os"

	"github.com/koukikitamura/autify-cli/internal"
	"github.com/mitchellh/cli"
	"github.com/sirupsen/logrus"
)

// These variables are set in build step
const (
	Version  = "unset"
	Revision = "unset"
)

func main() {
	c := cli.NewCLI("atf", Version)
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		internal.VersionCommandName: func() (cli.Command, error) {
			return &internal.VersionCommand{Version: Version, Revision: Revision}, nil
		},
		internal.RunCommandName: func() (cli.Command, error) {
			return &internal.RunCommand{}, nil
		},
		internal.ScenarioCommandName: func() (cli.Command, error) {
			return &internal.ScenarioCommand{}, nil
		},
		internal.ResultCommandName: func() (cli.Command, error) {
			return &internal.ResultCommand{}, nil
		},
	}

	exitCode, err := c.Run()
	if err != nil {
		logrus.Errorf("%+v", err)
	}

	os.Exit(exitCode)
}
