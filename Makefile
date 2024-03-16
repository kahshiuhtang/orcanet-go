build:
	go build -o bin/node cmd/main.go

run:
	bin/node

all: build run