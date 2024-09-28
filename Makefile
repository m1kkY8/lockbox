BINARY_NAME=./bin/client
GO_MAIN=./main.go

build:
	@go build 

run:
	@gochat -ip sjdoo.zapto.org

clean:
	@rm gochat

install:
	@go build
	@go install

.PHONY: build run install clean
