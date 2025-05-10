# Makefile for logwatch-llm project

.PHONY: all build run clean test fmt lint vet tidy mod vendor

all: build

build:
	@mkdir -p ./bin
	cd src/cmd/logwatch-llm && go build -o ../../../bin

run: build
	./bin

clean:
	rm -rf ./bin

fmt:
	gofmt -w src

lint:
	golint ./src/...

vet:
	go vet ./src/...

test:
	go test ./src/...

tidy:
	cd src && go mod tidy

mod:
	cd src && go mod download

vendor:
	cd src && go mod vendor
