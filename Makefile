BINARY_NAME=./bin/client
GO_MAIN=./src/cmd/main.go

build:
	@go build -o ${BINARY_NAME} ${GO_MAIN}

run:
	@go run ./src/cmd/main.go

clean:
	@rm -rf ./bin

.PHONY: build run
