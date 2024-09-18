BINARY_NAME=./bin/client
GO_MAIN=./main.go

build:
	@go build -o ${BINARY_NAME} ${GO_MAIN}

run:
	@go run ${GO_MAIN}

clean:
	@rm -rf ./bin

install:
	@go build
	@go install

.PHONY: build run install
