BINARY_NAME=./bin/client
GO_MAIN=./main.go

build:
	@go build 

run:
	@go build
	@go install
	@gochat -h sjdoo.zapto.org


dev:
	@go build
	@go install
	@gochat -h "localhost:1337" -u "test"

clean:
	@rm gochat

install:
	@go build
	@go install

.PHONY: build run install clean
