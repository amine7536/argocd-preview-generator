GO ?= go

.PHONY: all test lint fmt fmt-fix fmt-check tidy

all: fmt-check test lint

test:
	$(GO) test ./...

lint:
	golangci-lint run

fmt: fmt-fix

fmt-fix:
	@gofmt -w $$(find . -type f -name '*.go' -not -path './vendor/*')

fmt-check:
	@unformatted=$$(gofmt -l $$(find . -type f -name '*.go' -not -path './vendor/*')); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt needed on:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

tidy:
	$(GO) mod tidy
