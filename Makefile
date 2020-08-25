
# Install all the build and lint dependencies
setup:
	go get -u github.com/markbates/pkger/cmd/pkger
.PHONY: setup

build: embed
	go build
.PHONY: build

embed:
	pkger -include /payloads -include /go.mod -o ./cmd
.PHONY: embed

.DEFAULT_GOAL := build
