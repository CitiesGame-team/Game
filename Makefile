fmt:
	find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -s -w
	find . -name '*.go' -not -path "./vendor/*" | xargs goimports -w

run_main:
	go run main.go

build_main:
	go build -o main main.go

run: fmt run_main

build: fmt build_main
