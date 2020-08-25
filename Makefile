
# Install all the build and lint dependencies
setup:
	go get -u github.com/markbates/pkger/cmd/pkger
.PHONY: setup

build: embed
	go build
.PHONY: build

embed:
	@if [ ! -f ./cmd/pkged.checksums ] || [ -n "`find ./payloads ./go.mod -type f -print0 | xargs -0 md5sum | sort -k 2 | diff ./cmd/pkged.checksums -`" ] ; then \
		echo "Embedding updated payloads... " && \
		pkger -include /payloads -include /go.mod -o ./cmd && \
		find ./payloads ./go.mod -type f -print0 | xargs -0 md5sum | sort -k 2 > ./cmd/pkged.checksums ; \
		echo "Done!" ; \
	else \
		echo "Embeddings already up-to-date" ; \
	fi
	
.PHONY: embed

.DEFAULT_GOAL := build
