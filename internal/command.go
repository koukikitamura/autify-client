package internal

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/koukikitamura/autify-client/pkg/client"
	"github.com/sirupsen/logrus"
)

const (
	ExitCodeOk    int = 0
	ExitCodeError int = 1
)

const (
	VersionCommandName  = "version"
	ScenarioCommandName = "scenario"
	ResultCommandName   = "result"
	RunCommandName      = "run"
)

func RequireCredential() bool {
	if ok := client.CheckAccessToken(); !ok {
		fmt.Printf("Require %s environment variable\n", client.AccessTokenEnvName)
		return false
	}

	return true
}

type VersionCommand struct {
	Version  string
	Revision string
}

func (v *VersionCommand) Help() string {
	return "Print cli version & revision"
}

func (v *VersionCommand) Run(args []string) int {
	fmt.Printf("Version: %s\nRevision: %s\n", v.Version, v.Revision)
	return ExitCodeOk
}

func (r *VersionCommand) Synopsis() string {
	return "Print cli version"
}

type RunCommand struct{}

func (r *RunCommand) Help() string {
	return "Run test plan"
}

func (r *RunCommand) Run(args []string) int {
	if ok := RequireCredential(); !ok {
		return ExitCodeError
	}

	var projectId, planId, interval, timeout int
	var debug, showSpinner bool

	flags := flag.NewFlagSet(RunCommandName, flag.ContinueOnError)
	flags.IntVar(&projectId, "project-id", -1, "Specify project id")
	flags.IntVar(&planId, "plan-id", -1, "Specify plan id")
	flags.IntVar(&interval, "interval", 3, "Specify interval, unit is second")
	flags.IntVar(&timeout, "timeout", 3, "Specify interval, unit is minute")
	flags.BoolVar(&debug, "debug", false, "Print excution logs")
	flags.BoolVar(&showSpinner, "spinner", true, "Show spinner for waiting to finnish test")

	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	if projectId < 0 || planId < 0 {
		logrus.Error("project-id and plan-id is greater than or equal to zero.")
		return ExitCodeError
	}

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	autify := client.NewAutfiy(client.GetAccessToken())

	runResult, err := autify.RunTestPlan(planId)
	if err != nil {
		logrus.Errorf("%+v", err)
		return ExitCodeError
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	if showSpinner {
		s := spinner.New(spinner.CharSets[12], 100*time.Millisecond)
		s.Prefix = "Waiting for test plan to finish. "
		s.Start()
		defer s.Stop()
	}

	var testResult *client.TestPlanResult

	for {
		select {
		case <-ticker.C:
			testResult, err = autify.FetchResult(projectId, runResult.Attributes.Id)

			if err != nil {
				fmt.Print("\r\033[K")
				logrus.Errorf("%+v", err)
				return ExitCodeError
			}
			if testResult.Status != client.TestPlanStatuWaiting &&
				testResult.Status != client.TestPlanStatusQueuing &&
				testResult.Status != client.TestPlanStatusRunning {
				jsonStr, err := json.Marshal(*testResult)

				fmt.Print("\r\033[K")
				if err != nil {
					logrus.Errorf("%+v", err)
					return ExitCodeError
				}

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
	if ok := RequireCredential(); !ok {
		return ExitCodeError
	}

	var projectId, scenarioId int
	var debug bool

	flags := flag.NewFlagSet(ScenarioCommandName, flag.ContinueOnError)
	flags.IntVar(&projectId, "project-id", -1, "Specify project id")
	flags.IntVar(&scenarioId, "scenario-id", -1, "Specify project id")
	flags.BoolVar(&debug, "debug", false, "Print excution logs")

	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	if projectId < 0 || scenarioId < 0 {
		logrus.Error("project-id and scenario-id is greater than or equal to zero.")
		return ExitCodeError
	}

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	autify := client.NewAutfiy(client.GetAccessToken())

	scenario, err := autify.FetchScenario(projectId, scenarioId)
	if err != nil {
		logrus.Errorf("%+v", err)
		return ExitCodeError
	}

	jsonStr, err := json.Marshal(scenario)
	if err != nil {
		logrus.Errorf("%+v", err)
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
	if ok := RequireCredential(); !ok {
		return ExitCodeError
	}

	var projectId, resultId int
	var debug bool

	flags := flag.NewFlagSet(ResultCommandName, flag.ContinueOnError)
	flags.IntVar(&projectId, "project-id", -1, "Specify project id")
	flags.IntVar(&resultId, "result-id", -1, "Specify project id")
	flags.BoolVar(&debug, "debug", false, "Print excution logs")

	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	if projectId < 0 || resultId < 0 {
		logrus.Error("project-id and result-id is greater than or equal to zero.")
		return ExitCodeError
	}

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	autify := client.NewAutfiy(client.GetAccessToken())

	result, err := autify.FetchResult(projectId, resultId)
	if err != nil {
		logrus.Errorf("%+v", err)
		return ExitCodeError
	}

	jsonStr, err := json.Marshal(result)
	if err != nil {
		logrus.Errorf("%+v", err)
		return ExitCodeError
	}

	fmt.Println(string(jsonStr))
	return ExitCodeOk
}

func (r *ResultCommand) Synopsis() string {
	return "Get result"
}
