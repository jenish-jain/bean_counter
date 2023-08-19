.PHONY: all
export GO111MODULE=on

APP=bean_counter
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "${GREEN}$$(echo $$l | cut -f 1 -d':')${RESET}$$(echo $$l | cut -f 2- -d'#')\n"; done

compile:
	go build -o $(APP_EXECUTABLE) cmd/$(APP)/*.go

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

lint:
	@for p in $(ALL_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done

build: # complies format vet and lint your code
build: compile fmt vet lint

run: # run's your go code from source
run:
	go run cmd/$(APP)/main.go

start: # start's your generated binary
start:
	${APP_EXECUTABLE} start

clean:
	rm -rf out/
