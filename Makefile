.PHONY: test lint

test:
	@echo "[go test] running tests and collecting coverage metrics"
	@go test -v -tags all_tests -race -coverprofile=coverage.txt -covermode=atomic $$(go list ./... | grep -v /third_party/)

lint: lint-check-deps
	@echo "[golangci-lint] linting sources"
	@golangci-lint run \
		-E misspell \
		-E golint \
		-E gofmt \
		-E unconvert \
		--exclude-use-default=false \
		./...

lint-check-deps:
	@if [ -z `which golangci-lint` ]; then \
		echo "[go get] installing golangci-lint";\
		GO111MODULE=on go get -u github.com/golangci/golangci-lint/cmd/golangci-lint;\
	fi

build:
	@echo "[go build] build executable binary for development"
	@go build -o stock-crawler cmd/main.go

docker-m1:
	@echo "[docker build] build local docker image on Mac M1"
	@docker build -t samwang0723/stock-crawler:m1 -f build/docker/app/Dockerfile.local .

docker-amd64-deps:
	@echo "[docker buildx] install buildx depedency"
	@docker buildx create --name m1-builder
	@docker buildx use m1-builder
	@docker buildx inspect --bootstrap

docker-amd64:
	@echo "[docker buildx] build amd64 version docker image for Ubuntu AWS EC2 instance"
	@docker buildx use m1-builder
	@docker buildx build --load --platform=linux/amd64 -t samwang0723/stock-crawler:latest -f build/docker/app/Dockerfile .
