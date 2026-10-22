.PHONY: run test coverage

run:
	go run ./cmd/server

test:
	go test ./...

# Writes coverage/coverage.out and an HTML report at coverage/index.html.
coverage:
	mkdir -p coverage
	go test -coverprofile=coverage/coverage.out ./...
	go tool cover -func=coverage/coverage.out
	go tool cover -html=coverage/coverage.out -o coverage/index.html
