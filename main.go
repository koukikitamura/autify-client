package main

import (
	"fmt"
	"os"

	"github.com/koukikitamura/autify-cli/internal"
	"github.com/mitchellh/cli"
)

func main() {
	if ok := internal.CheckAccessToken(); !ok {
		fmt.Printf("Require %s environment variable\n", internal.AccessTokenEnvName)
		os.Exit(internal.ExitCodeError)
	}

	c := cli.NewCLI("atf", "1.0.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		internal.ScenarioCommandName: func() (cli.Command, error) {
			return &internal.ScenarioCommand{}, nil
		},
		internal.ResultCommandName: func() (cli.Command, error) {
			return &internal.ResultCommand{}, nil
		},
	}

	exitCode, err := c.Run()
	if err != nil {
		fmt.Printf("Failed to execute: %s\n", err.Error())
	}

	os.Exit(exitCode)
}
