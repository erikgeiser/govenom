setup:
	go get -u github.com/markbates/pkger/cmd/pkger
.PHONY: setup

build:
	go build
.PHONY: build

standalone:
	pkger -include /payloads -include /go.mod -o ./cmd
	go build
	rm ./cmd/pkged.go
		
.PHONY: embed

.DEFAULT_GOAL := build
