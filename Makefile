setup:
	go get -u github.com/markbates/pkger/cmd/pkger
.PHONY: setup

build:
	go build
.PHONY: build

standalone:
	pkger -include /payloads -o ./cmd
	go build
	rm ./cmd/pkged.go
.PHONY: standalone

dist:
	goreleaser --snapshot --skip-publish --rm-dist
	rm -f ./cmd/pkged.go
.PHONY: dist

.DEFAULT_GOAL := build
