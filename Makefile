c ?= 2

test:
	go test -count=$(c) ./...
lint:
	golangci-lint run
docs:
	godoc -http :8080
