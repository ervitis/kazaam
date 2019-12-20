test:
	go test -race -v ./...

lint:
	golangci-lint run

check: lint test

cover:
	go test -race -cover -covermode=count -coverprofile=cover.out ./...