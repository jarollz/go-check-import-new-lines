.PHONY: build
build:
	@go build -o ./bin/check-newline-in-imports ./cmd/main.go

.PHONY: run
run: build
	@./bin/check-newline-in-imports -n=$(MAX_NEW_LINE) -f=$(FILE_PATH)

.PHONY: init
init:
	@chmod +x ./.github/pre-commit.sh
	@chmod +x ./.dev/run-pre-commit.sh
	@cp -f ./.github/pre-commit.sh ./.git/hooks/pre-commit
	@chmod +x ./.git/hooks/pre-commit
	@echo "👍 GIT hooks configured"

.PHONY: test
test:
	@which gotest 2>/dev/null || go get -v github.com/rakyll/gotest
	@gotest -v -cover --race $$(go list ./... | grep -v vendor | grep -v cmd)

.PHONY: test-cover
test-cover:
	@echo "=== Compiling Test Coverage... Please Wait... ==="
	@which gotest 2>/dev/null || go get -v github.com/rakyll/gotest
	@gotest -gcflags="-l" -v -coverpkg=./... -coverprofile=coverage.out $$(go list ./... | grep -v vendor | grep -v cmd) > /dev/null
	@echo "=== Reporting Coverate Result ==="
	@go tool cover -func=coverage.out > coverage-report.out
	@cat coverage-report.out
	@echo "=== Open file 'coverage-report.out' to easily see full report ==="

.PHONY: lint
lint:
	@which golangci-lint 2>/dev/null || go get -v -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@golangci-lint run ./... -D errcheck -D goimports -E depguard -E gofmt -E nakedret -E goconst

.PHONY: check-buildable
check-buildable:
	@go build -o /dev/null ./...

.PHONY: check-dep
check-dep:
	@dep check

.PHONY: pre-commit
pre-commit:
	@.dev/run-pre-commit.sh