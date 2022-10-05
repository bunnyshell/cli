.DEFAULT_GOAL := build-local

.PHONY: build-local

build-local:
	goreleaser release --snapshot --rm-dist
