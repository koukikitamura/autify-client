# autify-cli
This package provides a unified command line interface to Autify.

## Installation
### macOS
```
brew install koukikitamura/autify-cli/autify-cli
```

### Linux
Download TAR archive from github release page.
```
curl -LSfs https://raw.githubusercontent.com/koukikitamura/autify-cli/main/scripts/install.sh | \
  sh -s -- \
    --git koukikitamura/autify-cli \
    --target autify-cli_linux_x86_64 \
    --to /usr/local/bin
```

## Configuration
Before using the autify-cli, you need to configure your credentials. You can use environment variable.

```
export AUTIFY_PERSONAL_ACCESS_TOKEN=<access token>
```

## Basic Commands
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
