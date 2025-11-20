SHELL := /bin/bash
export GO111MODULE=on

-include .env

TF_ENVS := dev preprod prod
GAR_REPOSITORY ?= us-central1-docker.pkg.dev/example-project/template-go
export GCP_PROJECT ?= example-project
export GCP_REGION ?= us-central1
export SERVICE_NAME ?= todo-app
export SKAFFOLD_DEFAULT_REPO ?= $(GAR_REPOSITORY)

.PHONY: help all bootstrap fmt lint test dev run integration-test load-test security-scan minikube-up plan preview-deploy

help: ## Print all documented targets (no variables required)
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-22s %s\n", $$1, $$2}'

all: fmt lint test integration-test load-test security-scan plan ## Run the full local CI/CD pipeline (uses GAR_REPOSITORY + TF_ENVS overrides)

bootstrap: ## Install buf, skaffold, minikube, and helm (uses GOPATH for buf install)
	@echo "Installing tooling"
	command -v buf >/dev/null || go install github.com/bufbuild/buf/cmd/buf@v1.29.0 || true
	command -v skaffold >/dev/null || curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64 && install skaffold /usr/local/bin/skaffold && rm skaffold || true
	command -v minikube >/dev/null || curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && install minikube /usr/local/bin/minikube && rm minikube || true
	command -v helm >/dev/null || curl -fsSL https://get.helm.sh/helm-v3.14.4-linux-amd64.tar.gz | tar -xz --strip-components=1 linux-amd64/helm && install helm /usr/local/bin/helm && rm helm

fmt: ## Run gofmt across the repo (no variables required)
	gofmt -w $(shell find . -name '*.go' -not -path './gen/*' -not -path './third_party/*')

lint: ## Run Buf, GolangCI-Lint, Terraform linting, and shell/docker linters (requires Docker for tflint)
	buf lint
	golangci-lint run ./...
	shellcheck scripts/*.sh || true
	docker run --rm -v $$PWD:/workspace ghcr.io/terraform-linters/tflint

test: ## Execute Go unit tests + container structure tests (requires GAR_REPOSITORY for image lookup)
	go test ./...
	container-structure-test test --image $(GAR_REPOSITORY)/todo-app:latest --config structure-test.yaml || true

run: ## Run the API locally without Kubernetes (requires DATABASE_URL/HTTP_ADDR/METRICS_ADDR)
	go run ./cmd/server

dev: minikube-up ## Start Skaffold dev loop against minikube (honors MINIKUBE_PROFILE + KUBECONFIG)
	skaffold dev -p minikube

minikube-up: ## Ensure a minikube cluster exists (uses MINIKUBE_PROFILE if provided)
	minikube start --kubernetes-version=v1.28.0

integration-test: ## Deploy to minikube and execute the Helm integration hook (requires MINIKUBE_PROFILE + KUBECONFIG)
	skaffold run -p minikube --tail=false
	helm test todo-app -n todo --filter integration --timeout 5m

load-test: ## Execute the Grafana k6 Helm hook (requires MINIKUBE_PROFILE + KUBECONFIG)
	helm test todo-app -n todo --filter load --timeout 10m

TF_COMPONENT ?= app

plan: ## terraform init/plan (backend disabled) for every environment (uses TF_ENVS + component tfvars)
	@for env in $(TF_ENVS); do \
		pushd terraform/components/$(TF_COMPONENT) >/dev/null; \
		VAR_FILE=$$(test -f $$env.tfvars && echo $$env.tfvars || echo $$env.tfvars.example); \
		terraform init -backend=false >/dev/null; \
		terraform plan -lock=false -input=false -refresh=false -var-file=$$VAR_FILE; \
		popd >/dev/null; \
	done

security-scan: ## Run osv-scanner, syft/grype, and gitleaks (requires internet access for vulnerability feeds)
	osv-scanner ./...
	syft . -o json > sbom.json
	grype sbom:sbom.json || true
	gitleaks detect --source .

preview-deploy: ## Deploy preview Cloud Run revision via Skaffold
	BRANCH="$(BRANCH)" scripts/preview-deploy.sh
