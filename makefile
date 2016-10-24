.PHONY: build

full: deps default build run

default: lint test

build:
	go build -o build/scraper *.go

lint:
	@golint ./...
	@go vet ./...

test:
	@go test -timeout 10s ./...

run:
	./build/scraper "http://wiprodigital.com/" 3 4

deps:
	@go get -v golang.org/x/tools/cmd/vet
	@go get -v github.com/golang/lint/golint
