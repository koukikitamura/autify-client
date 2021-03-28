package internal

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
)

const (
	ExitCodeOk    int = 0
	ExitCodeError int = 1
)

const (
	ScenarioCommandName = "scenario"
	ResultCommandName   = "result"
	RunCommandName      = "run"
)

type RunCommand struct{}

func (r *RunCommand) Help() string {
	return "Run test plan"
}

func (r *RunCommand) Run(args []string) int {
	var projectId, planId, interval, timeout int
	var debug bool

	flags := flag.NewFlagSet(RunCommandName, flag.ContinueOnError)
	flags.IntVar(&projectId, "project-id", -1, "Specify project id")
	flags.IntVar(&planId, "plan-id", -1, "Specify plan id")
	flags.IntVar(&interval, "interval", 3, "Specify interval, unit is second")
	flags.IntVar(&timeout, "timeout", 3, "Specify interval, unit is minute")
	flags.BoolVar(&debug, "debug", false, "Print excution logs")

	if err := flags.Parse(args); err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	autify := NewAutfiy(GetAccessToken())

	runResult, err := autify.RunTestPlan(planId)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	s := spinner.New(spinner.CharSets[12], 100*time.Millisecond)
	s.Prefix = "Waiting for test plan to finish. "
	s.Start()
	defer s.Stop()

	var testResult *TestPlanResult

	for {
		select {
		case <-ticker.C:
			testResult, err = autify.FetchResult(projectId, runResult.Attributes.Id)

			if err != nil {
				fmt.Println(err.Error())
				return ExitCodeError
			}
			if testResult.Status != TestPlanStatuWaiting &&
				testResult.Status != TestPlanStatusQueing &&
				testResult.Status != TestPlanStatusRunning {
				jsonStr, err := json.Marshal(*testResult)
				if err != nil {
					fmt.Println(err.Error())
					return ExitCodeError
				}

				fmt.Print("\r\033[K")
				fmt.Println(string(jsonStr))
				return ExitCodeOk
			}

		case <-time.After(time.Duration(timeout) * time.Minute):
			return ExitCodeOk
		}
	}
}

func (r *RunCommand) Synopsis() string {
	return "Run test plan"
}

type ScenarioCommand struct{}

func (s *ScenarioCommand) Help() string {
	return "Get scenario"
}

func (s *ScenarioCommand) Run(args []string) int {
	var projectId, scenarioId int
	var debug bool

	flags := flag.NewFlagSet(ScenarioCommandName, flag.ContinueOnError)
	flags.IntVar(&projectId, "project-id", -1, "Specify project id")
	flags.IntVar(&scenarioId, "scenario-id", -1, "Specify project id")
	flags.BoolVar(&debug, "debug", false, "Print excution logs")

	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	if debug {
		log.SetLevel(log.DebugLevel)
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

type ResultCommand struct{}

func (r *ResultCommand) Help() string {
	return "Get result"
}

func (r *ResultCommand) Run(args []string) int {
	var projectId, resultId int
	var debug bool

	flags := flag.NewFlagSet(ResultCommandName, flag.ContinueOnError)
	flags.IntVar(&projectId, "project-id", -1, "Specify project id")
	flags.IntVar(&resultId, "result-id", -1, "Specify project id")
	flags.BoolVar(&debug, "debug", false, "Print excution logs")

	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	if debug {
		log.SetLevel(log.DebugLevel)
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
