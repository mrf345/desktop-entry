test:
	go test -count=2 ./...
lint:
	golangci-lint run
docs:
	godoc -http :8080
