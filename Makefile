.PHONY: deps
deps:
	@go mod download

.PHONY: test
test:
	@go test ./...
