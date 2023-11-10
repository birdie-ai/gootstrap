SHELL = /bin/bash

# TOOLS VERSIONS
GO_VERSION=1.21.4
GOLANGCI_LINT_VERSION=v1.54.0

# configuration/aliases
version=$(shell git rev-parse --short HEAD)
base_image=us-central1-docker.pkg.dev/birdie-org/birdie/{{.Group}}/{{.Name}}
image=$(base_image):$(version)
devimage={{.Name}}-dev
# To avoid downloading deps everytime it runs on containers
gopkg=$(devimage)-gopkg
gocache=$(devimage)-gocache
devrun=docker run $(devrunopts) --rm \
	-v `pwd`:/app \
	-v $(gopkg):/go/pkg \
	-v $(gocache):/root/.cache/go-build \
	$(devimage)

covreport ?= coverage.txt

all: lint test image

## run isolated tests
.PHONY: test
test:
	go test ./... -timeout 10s -race -shuffle on

## Format go code
.PHONY: fmt
fmt:
	goimports -w .

## builds the service
.PHONY: service
service:
	go build -o ./cmd/{{.Name}}/{{.Name}} ./cmd/{{.Name}}

## runs the service locally
.PHONY: run
run: service
	./cmd/{{.Name}}/{{.Name}}

## tidy up go modules
.PHONY: mod
mod:
	go mod tidy

## lint the whole project
.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run ./...

## Run vulncheck
.PHONY: vulncheck
vulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## generates coverage report
.PHONY: test/coverage
test/coverage:
	go test -count=1 -coverprofile=$(covreport) ./...

## generates coverage report and shows it on the browser locally
.PHONY: test/coverage/show
test/coverage/show: test/coverage
	go tool cover -html=$(covreport)

## Configure a proxy to the service running on GKE (environment depends on which cluster is configured on kubectl).
.PHONY: port-forward
port-forward:
	kubectl port-forward -n enrichment service/{{.Name}} 8080:80

## Build the service image
.PHONY: image
image:
	docker build . \
		--build-arg GO_VERSION=$(GO_VERSION) \
		-t $(image)

## Build a production ready container image and run it locally for testing.
.PHONY: image/run
image/run: image
	docker run --rm -ti \
		-v $(gopkg):/go/pkg \
		$(image)

## Publish the service image
.PHONY: image/publish
image/publish: image
	docker push $(image)

## Releases to production
.PHONY: release
release: release_version=release-$(shell date '+%Y-%m-%d')-$(version)
release: release_image=$(base_image):$(release_version)
release:
	@echo "releasing from image: $(image)"
	@echo "release image:        $(release_image)"
	@echo "git tag:              $(release_version)"
	docker pull $(image)
	docker image tag $(image) $(release_image)
	docker push $(release_image)
	git tag -a $(release_version) -m "release to production: $(release_image)"
	git push origin $(release_version)

## Create the dev container image
.PHONY: dev/image
dev/image:
	docker build \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION) \
		-t $(devimage) \
		-f Dockerfile.dev \
		.

## Create a shell inside the dev container
.PHONY: dev
dev: devrunopts=-ti
dev: dev/image
	$(devrun)

## run a make target inside the dev container.
dev/%: dev/image
	$(devrun) make ${*}

## Display help for all targets
.PHONY: help
help:
	@awk '/^.PHONY: / { \
		msg = match(lastLine, /^## /); \
			if (msg) { \
				cmd = substr($$0, 9, 100); \
				msg = substr(lastLine, 4, 1000); \
				printf "  ${GREEN}%-30s${RESET} %s\n", cmd, msg; \
			} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
