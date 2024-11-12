.PHONY: lint, test,build, run, stoprm

lint:
	@golangci-lint run -v

test:
	@go test -coverprofile="coverage.out" ./... && \
	go tool cover -func="coverage.out" && \
	rm "coverage.out"

build:
	@echo not implemented

run:
	@echo not implemented

stoprm:
	@echo not implemented