test:
	go test -race -v ./...

lint:
	golangci-lint run

check: lint test

cover:
	go test -race -cover -coverprofile=cover.out ./...
	go tool cover -html=cover.out
	cat cover.out >> coverage.txt