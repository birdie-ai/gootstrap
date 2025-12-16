SHELL = /usr/bin/env bash

# TOOLS VERSIONS
ALPINE_VERSION=3.20

# configuration/aliases
version=$(shell git rev-parse --short HEAD)
base_image=us-central1-docker.pkg.dev/birdie-org/birdie/{{.Group}}/{{.Name}}
image=$(base_image):$(version)
covreport ?= coverage.txt

all: lint test image


## Format go code
.PHONY: fmt
fmt:
	go tool golang.org/x/tools/cmd/goimports -w .

## Install commands
.PHONY: install/%
install/%:
	go install ./cmd/${*}

## Builds commands
.PHONY: build/%
build/%:
	go build -o ./cmd/${*}/${*} ./cmd/${*}

## runs the service locally
.PHONY: run
run: build/{{.Name}}
	./cmd/{{.Name}}/{{.Name}}

## tidy up go modules
.PHONY: mod
mod:
	go mod tidy

## run tests
.PHONY: test
test:
	go test ./... -timeout 30s -short -race -shuffle on

## generates coverage report
.PHONY: test/coverage
test/coverage:
	go test -count=1 -coverprofile=$(covreport) ./...

## generates coverage report and shows it on the browser locally
.PHONY: test/coverage/show
test/coverage/show: test/coverage
	go tool cover -html=$(covreport)

## Build a binary suitable to be copied to a container image
.PHONY: image-bin/%
image-bin/%:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./cmd/${*}/${*} ./cmd/${*}

## Build the service image
.PHONY: image
image: image-bin/{{.Name}}
	docker build .  --build-arg ALPINE_VERSION=$(ALPINE_VERSION) -t $(image)

## Publish the service image
.PHONY: image/publish
image/publish: image
	docker push $(image)

# deploy the current (local git) revision to specified GKE cluster
.PHONY: deploy
deploy: guard-project guard-cluster guard-zone guard-image
	gcloud config set project $(project)
	gcloud container clusters get-credentials $(cluster) --zone $(zone) --project $(project)
	kubectl set image -n {{.Group}} deployment/{{.Name}} {{.Name}}=$(image)

# deploy in staging
.PHONY: deploy/staging
deploy/staging: project=k8s-staging-359211
deploy/staging: cluster=k8s-staging
deploy/staging: zone=us-central1-c
deploy/staging: deploy

# deploy in production
.PHONY: deploy/production
deploy/production: project=k8s-production-360315
deploy/production: cluster=k8s-cluster-production
deploy/production: zone=us-central1
deploy/production: deploy

## deadcode detects dead code on the project
.PHONY: deadcode
deadcode:
	go tool golang.org/x/tools/cmd/deadcode -test ./...

.PHONY: lint/mod
lint/mod:
	go mod tidy -diff

## checks code for vulnerabilities
.PHONY: vulncheck
vulncheck:
	go tool golang.org/x/vuln/cmd/govulncheck ./...

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

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "parameter '$*' not set"; \
		exit 1; \
	fi
