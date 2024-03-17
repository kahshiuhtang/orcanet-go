build:
	go build -o bin/node cmd/main.go

run:
	bin/node

build_tests:
	go build -o bin/tests test/test.go test/market.go test/blockchain.go
	
run_tests:	
	bin/tests

test: build_tests run_tests

all: build run