REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := "-X main.Revision=$(REVISION)"

build-cli:
	go build -ldflags $(LDFLAGS) -o dist/atf ./cmd/cli/main.go

test:
	go test ./...

test-scenario:
	@./scripts/atf_scenario.sh

test-result:
	@./scripts/atf_result.sh

test-run:
	@./scripts/atf_run.sh
