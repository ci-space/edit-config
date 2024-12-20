.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: verify
verify: test lint
