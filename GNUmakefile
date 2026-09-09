PKG_NAME            ?= internal
TESTARGS                ?= "-run=TestAcc"

default: testacc

# Please keep targets in alphabetical order

build: ## Build provider
	go install

gen:
	go generate

install-tools: ## Install required development tools
	go install $(shell go list -f '{{range .Imports}}{{.}} {{end}}' tools/tools.go)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

golangci-lint: ## Lint Go source (via golangci-lint)
	@echo "==> Checking source code with golangci-lint..."
	@golangci-lint run \
		--config .golangci.yml \
		./$(PKG_NAME)/...

test:
	go test ./...

# Run acceptance tests
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY:
	- build
	- generate
	- golangci-lint
	- install-tools
	- testacc
	- test