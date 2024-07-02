.PHONY: test lint bench lint-skip-fix migrate proto build build-docker install vendor deploy rollback

help: ## show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

APP_NAME?=stock-crawler
VERSION?=v2.0.1

SHELL = /bin/bash
SOURCE_LIST = $$(go list ./... | grep -v /third_party/ | grep -v /internal/app/pb)

###########
# install #
###########
## install: Install go dependencies
install:
	go mod tidy
	go mod download
	go get ./...

# vendor: Vendor go modules
vendor:
	go mod vendor

########
# test #
########

test: test-race test-leak test-coverage-ci ## launch all tests

test-race: ## launch all tests with race detection
	go test $(SOURCE_LIST)  -cover -race

test-leak: ## launch all tests with leak detection (if possible)
	go test $(SOURCE_LIST)  -leak

test-coverage-ci:
	go test -v $(SOURCE_LIST) -cover -race -covermode=atomic -coverprofile=coverage.out

test-coverage-report:
	go test -v $(SOURCE_LIST) -cover -race -covermode=atomic -coverprofile=coverage.out
	go tool cover -html=coverage.out

########
# lint #
########

lint: lint-check-deps ## lints the entire codebase
	@golangci-lint run ./... --config=./.golangci.yaml --timeout=15m && \
	if [ $$(gofumpt -e -l --extra cmd/ | wc -l) = "0" ] && \
		[ $$(gofumpt -e -l --extra internal/ | wc -l) = "0" ] && \
		[ $$(gofumpt -e -l --extra configs/ | wc -l) = "0" ] ; \
		then exit 0; \
	else \
		echo "these files needs to be gofumpt-ed"; \
		gofumpt -e -l --extra cmd/; \
		gofumpt -e -l --extra internal/; \
		gofumpt -e -l --extra configs/; \
	fi

lint-check-deps:
	@if [ -z `which golangci-lint` ]; then \
		echo "[go get] installing golangci-lint";\
		GO111MODULE=on go get -u github.com/golangci/golangci-lint/cmd/golangci-lint;\
	fi

lint-skip-fix: ## skip linting the system generate files
	@git checkout head internal/app/pb
	@git checkout head third_party/

#############
# benchmark #
#############

bench: ## launch benches
	go test $(SOURCE_LIST) -bench=. -benchmem | tee ./bench.txt

bench-compare: ## compare benches results
	benchstat ./bench.txt

#######
# sec #
#######

sec-scan: trivy-scan vuln-scan ## scan for security and vulnerabilities

trivy-scan: ## scan for sec issues with trivy (trivy binary needed)
	trivy fs --exit-code 1 --no-progress --severity CRITICAL ./

vuln-scan: ## scan for vuln issues with trivy (trivy binary needed)
	govulncheck ./...

###########
#  mock   #
###########

mock-gen: ## generate mocks
	go generate $(SOURCE_LIST)

############
# upgrades #
############

upgrade: ## upgrade dependencies (beware, it can break everything)
	go mod tidy && \
	go get -t -u ./... && \
	go mod tidy

##############
#   build    #
##############

build:
	@echo "[go build] build executable binary for development"
	@go build -o stock-crawler cmd/main.go

docker-build: lint test bench sec-scan docker-m1 ## build docker image in M1 device
	@printf "\nyou can now deploy to your env of choice:\nENV=dev make deploy\n"

docker-m1:
	@echo "[docker build] build local docker image on Mac M1"
	@docker build \
		-t samwang0723/$(APP_NAME):$(VERSION) \
		--build-arg LAST_MAIN_COMMIT_HASH=$(LAST_MAIN_COMMIT_HASH) \
		--build-arg LAST_MAIN_COMMIT_TIME=$(LAST_MAIN_COMMIT_TIME) \
		-f build/docker/app/Dockerfile.local .

docker-amd64-deps:
	@echo "[docker buildx] install buildx depedency"
	@docker buildx create --name m1-builder
	@docker buildx use m1-builder
	@docker buildx inspect --bootstrap

docker-amd64:
	@echo "[docker buildx] build amd64 version docker image for Ubuntu AWS EC2 instance"
	@docker buildx use m1-builder
	@docker buildx build \
		--load --platform=linux/amd64 \
		-t samwang0723/$(APP_NAME):$(VERSION) \
		--build-arg LAST_MAIN_COMMIT_HASH=$(LAST_MAIN_COMMIT_HASH) \
		--build-arg LAST_MAIN_COMMIT_TIME=$(LAST_MAIN_COMMIT_TIME) \
		-f build/docker/app/Dockerfile .

##################
# k8s Deployment #
##################
deploy:
	@kubectl apply -f deployments/helm/stock-crawler/deployment.yaml
	@kubectl rollout status deployment/stock-crawler

rollback:
	@kubectl rollout undo deployment/stock-crawler

#############
# changelog #
#############

MOD_VERSION = $(shell git describe --abbrev=0 --tags `git rev-list --tags --max-count=1`)

MESSAGE_CHANGELOG_COMMIT="chore(changelog): update CHANGELOG.md for $(MOD_VERSION)"

changelog-gen: ## generates the changelog in CHANGELOG.md
	@git cliff -o ./CHANGELOG.md && \
	printf "\nchangelog generated!\n"
	git add CHANGELOG.md

changelog-commit:
	git commit -m $(MESSAGE_CHANGELOG_COMMIT) ./CHANGELOG.md
