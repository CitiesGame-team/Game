run: fmt run_main

fmt:
	find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -s -w
	find . -name '*.go' -not -path "./vendor/*" | xargs goimports -w

run_main:
	go run main.go

build_main:
	go build -o main main.go

build: fmt build_main

run_init:
	go run main.go -init

init: fmt run_init