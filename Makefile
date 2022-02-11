REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := "-X main.Revision=$(REVISION)"

build:
	go build -ldflags $(LDFLAGS) -o atf

test:
	go test ./...

test-scenario:
	@./scripts/atf_scenario.sh

test-result:
	@./scripts/atf_result.sh

test-run:
	@./scripts/atf_run.sh
