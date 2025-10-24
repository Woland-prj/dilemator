SHELL := /bin/bash

LOCAL_BIN := $(CURDIR)/build/bin
BASE_STACK := docker compose -f docker-compose.yaml

.PHONY: default
default:
	@make help

.PHONY: help
help:
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-25s %s\n", $$1, $$2}'

.PHONY: compose-up
compose-up: ## Run docker compose (without backend and reverse proxy)
	$(BASE_STACK) up --build -d postgres minio createbuckets

.PHONY: compose-down
compose-down: ## Down docker compose
	$(BASE_STACK) down --remove-orphans

.PHONY: swag-v1
swag-v1: ## swag init
	docker run --rm -v $(CURDIR):/code ghcr.io/swaggo/swag:latest init -g internal/router/setup/setup.go

.PHONY: deps
deps: ## deps tidy + verify
	go mod tidy
	go mod verify
	npm --prefix web install

.PHONY: deps-audit
deps-audit: ## check dependencies vulnerabilities
	govulncheck ./...

.PHONY: format
format: ## Run code formatter
	gofumpt -l -w .
	gci write . --skip-generated -s standard -s default
	templ fmt internal/view

.PHONY: run
run: deps swag-v1 ## Run application
	go mod download
	CGO_ENABLED=0 env $$(grep -v '^#' .env.example | xargs) go run ./cmd/app

.PHONY: build
build: build-templ build-css build-js build-app ## build all stuff

.PHONY: build-app
build-app: ## build appliction
	go build -o $(LOCAL_BIN) ./cmd/app/main.go

.PHONY: build-templ
build-templ: ## build templates to go files
	templ generate

.PHONY: build-css
build-css: ## build css styles
	npm --prefix web run build:css --

.PHONY: build-js
build-js: ## build javascript bundle
	npm --prefix web run build:js -- --minify

.PHONY: watch
watch: ## watch rebuiding all stuff
	$(MAKE) -j5 watch-app watch-templ watch-css watch-js watch-assets

.PHONY: watch-app
watch-app: ## watch rebuiding application
	CGO_ENABLED=0 env $$(grep -v '^#' .env.example | xargs) \
	go run github.com/air-verse/air@latest \
	--build.cmd "$(MAKE) build-app" \
	--build.bin "$(LOCAL_BIN)" \
	--build.include_ext "go" \
	--build.exclude_dir "bin,web"

.PHONY: watch-templ
watch-templ: ## watch rebuiding templates
	templ generate \
	--watch \
	--proxy "http://localhost:8080" \
	--open-browser=false

.PHONY: watch-css
watch-css: ## watch rebuiding css
	npm --prefix web run build:css -- --watch=always

.PHONY: watch-js
watch-js: ## watch rebuiding js
	npm --prefix web run build:js -- --watch=forever

.PHONY: watch-assets
watch-assets: ## watch rebuiding static assets
	go run github.com/air-verse/air@latest \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.exclude_dir "" \
	--build.include_dir "web/public/assets" \
	--build.include_ext "css,js" \
	--build.delay "100"

.PHONY: docker-rm-volume
docker-rm-volume: ## remove docker volume
	docker volume rm dilemator_pg_data || true
	docker volume rm dilemator_minio_data || true

.PHONY: linter-golangci
linter-golangci: ## check by golangci linter
	golangci-lint run --fix

.PHONY: linter-hadolint
linter-hadolint: ## check by hadolint linter
	/bin/ls Dockerfile* | xargs hadolint

.PHONY: linter-dotenv
linter-dotenv: ## check by dotenv linter
	dotenv-linter check .env.example

.PHONY: test
test: ## run test
	go generate ./...
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...

.PHONY: migrate-local
migrate-local: ## Run Liquibase migrations against local PostgreSQL
	$(BASE_STACK) run --rm liquibase

.PHONY: pre-commit
pre-commit: ## run pre-commit checks
	$(MAKE) linter-hadolint
	$(MAKE) linter-dotenv
	$(MAKE) deps
	$(MAKE) deps-audit
	$(MAKE) swag-v1
	$(MAKE) format
	$(MAKE) linter-golangci
	#$(MAKE) test
	$(MAKE) build
