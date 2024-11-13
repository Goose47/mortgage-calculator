.PHONY: lint, test,build, run, stoprm

lint:
	@golangci-lint run -v

test:
	@go test -coverprofile="coverage.out" ./... && \
	go tool cover -func="coverage.out" && \
	rm "coverage.out"

build:
	@docker build -t mortgage_calculator . && \
	docker image ls | head -n 1 && \
	docker image ls | grep mortgage_calculator

run:
	@docker run -p "8080:8080" -d --name "mortgage_calculator" mortgage_calculator && \
	docker ps | head -n 1 && \
	docker ps | grep mortgage_calculator

stoprm:
	@docker stop mortgage_calculator && \
	docker container rm mortgage_calculator