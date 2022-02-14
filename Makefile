test:
	go test -race ./...

vet:
	go vet ./...

format:
	gofmt -s -w .