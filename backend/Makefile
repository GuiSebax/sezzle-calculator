.PHONY: run test coverage

run:
	go run ./cmd/server

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
