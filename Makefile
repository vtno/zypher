build:
	go build -o build/zypher cmd/zypher.go
.PHONY: build

gen:
	go generate ./...
.PHONY: generate

test:
	go test ./...
.PHONY: test

fmt:
	go fmt ./...
.PHONY: fmt
