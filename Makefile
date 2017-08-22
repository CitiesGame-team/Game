fmt:
	find . -name '*.go' | xargs gofmt -s -w

run_main:
	go run main.go

build_main:
	go build -o main main.go

run: fmt run_main

build: fmt build_main
