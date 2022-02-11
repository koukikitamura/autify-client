# autify-client

The autify-client provides [Autify](https://autify.com/) client go package & CLI.

## CLI Installation

### Homebrew

```
brew tap koukikitamura/autify-cli
brew install autify-cli
```

### Download a binary

Download TAR archive from [Github release page](https://github.com/koukikitamura/autify-client/releases).

```
curl -LSfs https://raw.githubusercontent.com/koukikitamura/autify-client/main/scripts/install.sh | \
  sh -s -- \
    --git koukikitamura/autify-client \
    --target autify-cli_linux_x86_64 \
    --to /usr/local/bin
```

## CLI Configuration

Before using the autify-cli, you need to configure your credentials. You can use environment variable.

```
export AUTIFY_PERSONAL_ACCESS_TOKEN=<access token>
```

## CLI Basic Commands

An autify-cli command has the following structure:

```
$ atf <command> [options]
```

To run test plan and wait to finish, the command would be:

```
$ atf run --project-id=999 --plan-id=999
{"id":999,"status":"passed","duration":26251,"started_at":"2021-03-28T11:03:31.288Z","finished_at":"2021-03-28T11:03:57.54Z","created_at":"2021-03-28T11:03:04.716Z","updated_at":"2021-03-28T11:04:00.738Z","test_plan":{"id":999,"name":"main flow","created_at":"2021-03-26T08:25:12.987Z","updated_at":"2021-03-26T08:33:45.462Z"}}
```

To fetch scenario, the command would be:

```
$ atf scenario --project-id=999 --scenario-id=999
{"id":999,"name":"login","created_at":"2021-03-26T07:53:20.039Z","updated_at":"2021-03-26T08:20:51.86Z"}
```

To fetch test plan excution result, the command would be:

```
$ atf result --project-id=999 --result-id=999
{"id":999,"status":"waiting","duration":26621,"started_at":"2021-03-26T10:09:12.915Z","finished_at":"2021-03-26T10:09:39.537Z","created_at":"2021-03-26T10:08:54.769Z","updated_at":"2021-03-26T10:09:44.542Z","test_plan":{"id":999,"name":"main flow","created_at":"2021-03-26T08:25:12.987Z","updated_at":"2021-03-26T08:33:45.462Z"}}
```

## Terms

### project-id

The path of the Autify dashboard home is `/projects/[project-id]`. This path parameter is the project-id.

### plan-id

The path of the test plan's detail page is `/projects/[project-id]/test_plans/[plan-id]`. This path parameter is the plan-id.

### scenario-id

The path of the scenario's detail page is `/projects/[project-id]/scenarios/[scenario-id]`. This path parameter is the scenario-id.

### result-id

The path of the result's detail page is `/projects/[project-id]/results/[result-id]`. This path parameter is the result-id.

## Go package

The following is the code to run the test plan and poll its status.

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/koukikitamura/autify-client/pkg/client"
)

const (
	ExitCodeOk    int = 0
	ExitCodeError int = 1
)

func main() {
	var projectId = 999
	var planId = 999

	autify := client.NewAutfiy(client.GetAccessToken())

	runResult, err := autify.RunTestPlan(planId)
	if err != nil {
		fmt.Println("Error: Failed to run the test plan")
		os.Exit(ExitCodeError)
	}

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	defer ticker.Stop()

	var testResult *client.TestPlanResult
	for {
		select {
		case <-ticker.C:
			testResult, err = autify.FetchResult(projectId, runResult.Attributes.Id)
			if err != nil {
				fmt.Println("Error: Failed to fetch the result")
				os.Exit(ExitCodeError)
			}

			if testResult.Status != client.TestPlanStatuWaiting &&
				testResult.Status != client.TestPlanStatusQueuing &&
				testResult.Status != client.TestPlanStatusRunning {
				jsonStr, err := json.Marshal(*testResult)
				if err != nil {
					fmt.Println("Error: Failed to marshal the test result")
					os.Exit(ExitCodeError)
				}

				fmt.Println((string(jsonStr)))
				os.Exit(ExitCodeOk)
			}

		case <-time.After(time.Duration(5) * time.Minute):
			fmt.Println("Error: Timeout")
			os.Exit(1)
		}
	}
}
```
