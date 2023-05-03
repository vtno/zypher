build:
	go build -o build cmd/zypher.go
.PHONY: build

generate:
	go generate ./...
.PHONY: generate

test:
	go test ./...
.PHONY: test
