BIN=agentpane

.PHONY: build install test fmt

build:
	go build -o $(BIN) ./cmd/agentpane

install:
	go install ./cmd/agentpane

test:
	go test ./...

fmt:
	gofmt -w cmd internal
