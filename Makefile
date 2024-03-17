build:
	go build -o bin/node cmd/main.go

run:
	bin/node

run_tests:
	go build -o bin/node_tests test/test.go test/market.go
	bin/node_tests

test: run_tests

all: build run