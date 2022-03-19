# Run CI tasks
ci: lint build test
.PHONY: ci

# Format all files
fmt:
	@echo "==> Formatting source"
	@gofmt -s -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")
	@echo "==> Done"
.PHONY: fmt

# Tidy the go.mod file
tidy:
	@echo "==> Cleaning go.mod"
	@go mod tidy
	@echo "==> Done"
.PHONY: tidy

# Build the binary
build:
	@find ./cmd/* -maxdepth 1 -type d -exec go build {} \;
.PHONY: build

# Build the binary
build-linux:
	@GOOS=linux GOARCH=amd64 find ./cmd/* -maxdepth 1 -type d -exec go build {} \;
.PHONY: build-linux

# Run all tests
test:
	@go test -cover -race ./...
.PHONY: test

# Lint the project
lint:
	@golangci-lint run --go 1.18 ./...
.PHONY: lint

# Build the docker image
docker: build-linux
	@docker build -t ren .
.PHONY: docker
