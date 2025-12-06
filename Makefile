.PHONY: format
format:
	gofmt -w -s .

.PHONY: test
test:
	go test -v -race -timeout 3m -coverpkg=./... -coverprofile=coverage.out ./...

.PHONY: test-ddn
test-ddn:
	./tests/scripts/test-ddn.sh

# Install golangci-lint tool to run lint locally
# https://golangci-lint.run/usage/install
.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

.PHONY: build-jsonschema
build-jsonschema:
	cd jsonschema && go run .

