
.PHONY: test-integration
test-integration:
	@go test ./test/integration/...

.PHONY: test-unit
test-unit:
	@go test ./internal/...

.PHONY: test
test: test-unit test-integration

.PHONY: run
run:
	@go run ./...

.PHONY: build
build:
	@go build

