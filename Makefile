build:
	go build -o build
.PHONY: build

generate:
	go generate ./...
.PHONY: generate

test:
	go test ./...
.PHONY: test
