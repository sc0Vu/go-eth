GOPATH= $(shell go env GOPATH)

.PHONY: default
default: lint test clean

.PHONY: lint
lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.37.0
	$(GOPATH)/bin/golangci-lint run -e gosec ./...
	go fmt ./...

.PHONY: test
test: 
	go test --race ./...


.PHONY: clean
clean:
	go clean -cache ./...
	go mod tidy
