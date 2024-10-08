# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# Brew Installation
#
# 	Install Homebrew:
#	$ /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# ==============================================================================
# Install Tooling and Dependencies
#
#	This project uses Docker and it is expected to be installed. Please provide Docker at least 4 CPUs.
#
#	Run these commands to install everything needed.
#	$ make dev-brew
#	$ make dev-docker
#	$ make dev-gotooling
#
# Add the following to your shell profile (e.g. ~/.bashrc, ~/.zshrc, ~/.profile, or ~/.bash_profile):
# export PATH=$GOPATH/bin:$PATH

# ==============================================================================
# Running Test
#
#	$ make test

# ==============================================================================
# Running The Project
#
#	$ make compose-build-up
#	$ make token
#	$ export TOKEN=<token>
#	$ make users

# ==============================================================================
# Define dependencies

GOLANG          := golang:1.22
ALPINE          := alpine:3.20
POSTGRES        := postgres:16.3
GRAFANA         := grafana/grafana:10.4.0
PROMETHEUS      := prom/prometheus:v2.52.0
TEMPO           := grafana/tempo:2.5.0
LOKI            := grafana/loki:2.9.0
PROMTAIL        := grafana/promtail:2.9.0

BASE_IMAGE_NAME := local/nhannguyenacademy
ECOMMERCE_APP   := ecommerce
VERSION         := 1.0.0
ECOMMERCE_IMAGE := $(BASE_IMAGE_NAME)/$(ECOMMERCE_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)

# VERSION       := "1.0.0-$(shell git rev-parse --short HEAD)"

# ==============================================================================
# Install dependencies

dev-gotooling:
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-brew:
	brew update
	brew list golang-migrate || brew install golang-migrate
	brew list mockery || brew install mockery # https://vektra.github.io/mockery/latest/notes/#internal-error-package-without-types-was-imported

dev-docker:
	docker pull $(POSTGRES) & \
	docker pull $(GOLANG) & \
	docker pull $(ALPINE) & \
#	docker pull $(GRAFANA) & \
#	docker pull $(PROMETHEUS) & \
#	docker pull $(TEMPO) & \
#	docker pull $(LOKI) & \
#	docker pull $(PROMTAIL) & \
	wait;

// ==============================================================================

run:
	export ECOMMERCE_DB_HOST=localhost ECOMMERCE_SERVER_HOST=0.0.0.0:8081; go run cmd/ecommerce/main.go

# ==============================================================================
# Building containers


build: ecommerce

ecommerce:
	docker build \
		-f build/ecommerce/ecommerce.dockerfile \
		-t $(ECOMMERCE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ==============================================================================
# Docker Compose

compose-up:
	cd ./deployments/compose/ && docker compose -f docker_compose.yaml -p compose up -d

compose-build-up: build compose-up

compose-down:
	cd ./deployments/compose/ && docker compose -f docker_compose.yaml down -v

compose-logs:
	cd ./deployments/compose/ && docker compose -f docker_compose.yaml logs --follow

# ==============================================================================
# Run admin commands from local

# example: make create-migration name=create_table_users
create-migration:
	migrate create -ext sql -dir internal/sdk/sdkbus/migrate/migrations -seq $(name)

migrate:
	export ECOMMERCE_DB_HOST=localhost; go run tools/admin/main.go migrate

migrate-down:
	export ECOMMERCE_DB_HOST=localhost; go run tools/admin/main.go migrate-down

seed:
	export ECOMMERCE_DB_HOST=localhost; go run tools/admin/main.go seed

admin-users:
	export ECOMMERCE_DB_HOST=localhost; go run tools/admin/main.go users

admin-gen-token:
	export ECOMMERCE_DB_HOST=localhost; go run tools/admin/main.go gentoken 97ee07e2-ebbb-4c69-a681-d5fe165c2cb9 54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

# ==============================================================================
# Metrics and Tracing

grafana:
	open -a "Google Chrome" http://localhost:3100/

# ==============================================================================
# Unit-tests, linting, and security checks

test: test-r lint vuln-check

test-r:
	CGO_ENABLED=1 go test -race -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

mock:
	mockery --config=./configs/mockery/mockery.yaml

# ==============================================================================
# Hitting endpoints

token:
	curl -il \
	--user "admin@email.com:abc123" http://localhost:8080/api/v1/auth/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

create-user:
	curl -il -X POST \
	-H "Authorization: Bearer ${TOKEN}" \
	-H 'Content-Type: application/json' \
	-d '{"name":"Nhan Nguyen","email":"nhannguyen@email.com","roles":["ADMIN"],"password":"test123","passwordConfirm":"test123"}' \
	http://localhost:8080/api/v1/users

users:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:8080/api/v1/users?page=1&rows=2"

load:
	hey -m GET -c 100 -n 1000 \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:8080/api/v1/users?page=1&rows=2"

otel-test:
	curl -il \
	-H "Traceparent: 00-918dd5ecf264712262b68cf2ef8b5239-896d90f23f69f006-01" \
	--user "admin@email.com:abc123" http://localhost:8080/api/v1/users/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

