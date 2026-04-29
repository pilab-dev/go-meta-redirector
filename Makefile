.PHONY: help build run test clean docker-build docker-run docker-push k8s-deploy k8s-delete helm-package helm-install

BINARY_NAME=go-meta-redirector
DOCKER_IMAGE=ghcr.io/pilab-dev/go-meta-redirector
VERSION?=latest

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go binary
	go build -o $(BINARY_NAME) .

run: ## Run the server locally
	go run . :8080

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	go clean

# Docker targets
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

docker-run: ## Run Docker container locally
	docker run -d -p 8080:8080 --name $(BINARY_NAME) \
		$(DOCKER_IMAGE):$(VERSION)

docker-stop: ## Stop Docker container
	docker stop $(BINARY_NAME) || true
	docker rm $(BINARY_NAME) || true

docker-push: ## Push Docker image to GHCR
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

docker-test: docker-stop ## Test Docker container
	$(MAKE) docker-run
	sleep 2
	curl -H "Host: go.pilab.hu" "http://localhost:8080/cloud/log?go-get=1"
	$(MAKE) docker-stop

# Kubernetes targets
k8s-deploy: ## Deploy to Kubernetes
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/deployment.yaml
	kubectl apply -f k8s/service.yaml
	kubectl apply -f k8s/ingress.yaml

k8s-delete: ## Delete from Kubernetes
	kubectl delete -f k8s/

k8s-status: ## Check Kubernetes deployment status
	kubectl get all -n go-meta-redirector

# Helm targets
helm-lint: ## Lint Helm chart
	helm lint helm/

helm-package: ## Package Helm chart
	helm package helm/ -d dist/

helm-install: ## Install Helm chart
	helm install go-meta-redirector helm/

helm-upgrade: ## Upgrade Helm release
	helm upgrade go-meta-redirector helm/

helm-uninstall: ## Uninstall Helm release
	helm uninstall go-meta-redirector

# Release targets
release: ## Create a new release (VERSION required)
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "latest" ]; then \
		echo "Error: VERSION is required. Usage: make release VERSION=v1.0.1"; \
		exit 1; \
	fi
	git tag $(VERSION)
	git push origin $(VERSION)
	gh release create $(VERSION) --title "Release $(VERSION)" \
		--notes "See CHANGELOG.md for details" \
		--repo pilab-dev/go-meta-redirector

# Development helpers
dev: ## Start development server with hot reload (requires air)
	@which air > /dev/null 2>&1 || go install github.com/cosmtrek/air@latest
	air

fmt: ## Format code
	gofmt -w .

lint: ## Lint code
	golangci-lint run ./... || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"

tidy: ## Tidy Go modules
	go mod tidy

# CI targets
ci: lint test build ## Run CI checks locally

.DEFAULT_GOAL := help
