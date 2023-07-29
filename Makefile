build:
	go build -o build/zypher cmd/zypher/zypher.go
.PHONY: build

build/signer:
	go build -o build/signer cmd/signer/signer.go
.PHONY: build/signer

gen:
	go generate ./...
.PHONY: generate

test:
	go test ./...
.PHONY: test

fmt:
	go fmt ./...
.PHONY: fmt

xbuild:
	gox -os="linux darwin windows" -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" github.com/vtno/zypher/cmd/zypher
.PHONY: xbuild

lint:
	golangci-lint run
.PHONY: lint
