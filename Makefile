.PHONY: list date

lint:
	@golangci-lint run -v

test:
	go list ./... | \
	grep -v '/mocks/' | \
 	xargs go test -tags test -coverprofile="coverage.out" ./... && \
 	go tool cover -func="coverage.out" && \
 	rm "coverage.out"

build:
	@echo not implemented

run:
	@echo not implemented

stoprm:
	@echo not implemented