package main

import (
	"os"

	"github.com/koukikitamura/autify-cli/internal"
	"github.com/mitchellh/cli"
	"github.com/sirupsen/logrus"
)

func main() {
	c := cli.NewCLI("atf", "1.0.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
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
