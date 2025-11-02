.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: install-tools
install-tools: ## Install development tools
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/bufbuild/buf/cmd/buf@v1.59.0

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy
	go mod verify

.PHONY: migrate-up
migrate-up: ## Run database migrations up
	migrate -path db/migrations -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down: ## Run database migrations down
	migrate -path db/migrations -database "$(DATABASE_URL)" down

.PHONY: migrate-create
migrate-create: ## Create new migration (usage: make migrate-create NAME=create_table)
	migrate create -ext sql -dir db/migrations -seq $(NAME)

.PHONY: sqlc-generate
sqlc-generate: ## Generate sqlc code
	sqlc generate

.PHONY: proto-gen
proto-gen: ## Generate protobuf code
	buf generate

.PHONY: proto-lint
proto-lint: ## Lint protobuf files
	buf lint

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t go-microservices-template:latest -f deployments/docker/Dockerfile .

.PHONY: docker-up
docker-up: ## Start all services with docker compose
	docker compose -f deployments/docker/docker-compose.yml up -d

.PHONY: docker-down
docker-down: ## Stop all services
	docker compose -f deployments/docker/docker-compose.yml down

.PHONY: docker-logs
docker-logs: ## Show logs from all services
	docker compose -f deployments/docker/docker-compose.yml logs -f

.PHONY: docker-clean
docker-clean: ## Clean all docker resources
	docker compose -f deployments/docker/docker-compose.yml down -v
	docker system prune -af

.PHONY: k8s-apply
k8s-apply: ## Apply all Kubernetes manifests
	kubectl apply -f deployments/k8s/namespace.yaml
	kubectl apply -f deployments/k8s/configmap.yaml
	kubectl apply -f deployments/k8s/secret.yaml
	kubectl apply -f deployments/k8s/postgres-deployment.yaml
	kubectl apply -f deployments/k8s/rabbitmq-deployment.yaml
	kubectl apply -f deployments/k8s/api-deployment.yaml
	kubectl apply -f deployments/k8s/hpa.yaml

.PHONY: k8s-delete
k8s-delete: ## Delete all Kubernetes resources
	kubectl delete -f deployments/k8s/ --recursive

.PHONY: k8s-status
k8s-status: ## Check status of all pods
	kubectl get all -n microservices
