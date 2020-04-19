PKGS := $(shell go list ./... | grep -v /vendor)
BINARY := pstree-json
VERSION ?= latest

.PHONY: test
test:
	go test -v $(PKGS)

.PHONY: release
release:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o release/$(BINARY)-$(VERSION)-linux-amd64

.PHONY: clean
clean:
	rm release/*

default: test
