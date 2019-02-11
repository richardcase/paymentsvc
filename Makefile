BUILDCOMMIT := $(shell git describe --dirty --always)
BUILDDATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VER_FLAGS=-X main.commit=$(BUILDCOMMIT) -X main.date=$(BUILDDATE)

.DEFAULT_GOAL:=help

##@ Build

.PHONY: build
build: ## Build the lambda function in relase
	go get -u
	go mod vendor
	env GOOS=linux go build -ldflags="-s -w" -o bin/payments cmd/paymentsvc/main.go

.PHONY: build-debug
build-debug: ## Build the lambda function with debug symbols
	go get -u
	go mod vendor
	env GOOS=linux go build -ldflags="-N -l" -o bin/payments cmd/paymentsvc/main.go
	env GOARCH=amd64 GOOS=linux go build -o bin/dlv github.com/derekparker/delve/cmd/dlv

.PHONY: build-local
build-local: ## Build the function to run locally for dev/testing purposes
	go get -u
	go mod vendor
	go build -ldflags="-s -w" -o bin/payments cmd/paymentsvc/main.go

##@ Testing & CI

.PHONY: test
test: generate ## Run unit tests
	@git diff --exit-code ./pkg/repository/mocks > /dev/null || (git --no-pager diff ./pkg/repository/mocks; exit 1)
	@go test -v -covermode=count -coverprofile=coverage.out ./pkg/... ./cmd/...

.PHONY: integration-test-local
integration-test-local: ## Run integration tests locally
	echo "This assumes you are running the service locally (i.e with make start-local) and DynamoDb Emulator"
	@go test -v -tags integration  -timeout 5m ./test/integration/... 

.PHONY: lint
lint: ## Run linting over codebase
	golangci-lint run

.PHONY: ci
ci: test lint ## Target for CI system to invoke to run tests and linting


##@ Code/Config Generation

.PHONY: generate-sam
generate-sam: ## Generate SAM template from Serverless
	serverless sam export --output ./template.yml

.PHONY: generate ## Generate mocks
generate:
	@go generate ./pkg/repository/mocks


##@ Run locally

.PHONY: start-local-db
start-local-db: ## Start the DynamoDb emulator
	docker run --rm -p 8000:8000 --name dynamodb amazon/dynamodb-local -jar DynamoDBLocal.jar -inMemory  -sharedDb

.PHONY: populate-local-db
populate-local-db: ## Create table/data in DynamoDb
	go run tools/db-populate/main.go --endpoint http://127.0.0.1:8000

.PHONY: start-local-func
start-local-func: build-local ## Start the functions running locally
	env REGION=eu-west-2 DB_TABLE=payments DB_ENDPOINT_OVERRIDE=http://127.0.0.1:8000 RUN_LOCAL=true bin/payments

##@ Deployment

.PHONY: clean
clean: ## Cleanup the build/vendor folder
	rm -rf ./bin ./vendor

.PHONY: deploy
deploy: clean build ## Deploy the functions to AWS
	sls deploy --verbose

##@ Utility

.PHONY: fmt
fmt: ## Format all the source code using gofmt
	@gofmt -l -w $(SRC)

.PHONY: help
help:  ## Display this help. Thanks to https://suva.sh/posts/well-documented-makefiles/
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

