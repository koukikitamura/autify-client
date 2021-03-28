package internal

import (
	"encoding/json"
	"flag"
	"fmt"
)

const (
	ExitCodeOk    int = 0
	ExitCodeError int = 1
)

const (
	ScenarioCommandName = "scenario"
	ResultCommandName   = "result"
)

// atf senario --project-id=id --scenario-id=id
type ScenarioCommand struct{}

func (s *ScenarioCommand) Help() string {
	return "Get scenario"
}

func (s *ScenarioCommand) Run(args []string) int {
	var projectId, scenarioId int

	flags := flag.NewFlagSet(ScenarioCommandName, flag.ContinueOnError)
	flags.IntVar(&projectId, "project-id", -1, "Specify project id")
	flags.IntVar(&scenarioId, "scenario-id", -1, "Specify project id")

	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	autify := NewAutfiy(GetAccessToken())

	scenario, err := autify.FetchScenario(projectId, scenarioId)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	jsonStr, err := json.Marshal(scenario)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	fmt.Println(string(jsonStr))
	return ExitCodeOk
}

func (s *ScenarioCommand) Synopsis() string {
	return "Get scenario"
}

// atf result --project-id=id --result-id=id
type ResultCommand struct{}

func (r *ResultCommand) Help() string {
	return "Get result"
}

func (r *ResultCommand) Run(args []string) int {
	var projectId, resultId int

	flags := flag.NewFlagSet(ResultCommandName, flag.ContinueOnError)
	flags.IntVar(&projectId, "project-id", -1, "Specify project id")
	flags.IntVar(&resultId, "result-id", -1, "Specify project id")

	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	autify := NewAutfiy(GetAccessToken())

	result, err := autify.FetchResult(projectId, resultId)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	jsonStr, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	fmt.Println(string(jsonStr))
	return ExitCodeOk
}

func (r *ResultCommand) Synopsis() string {
	return "Get result"
}
