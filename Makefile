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
	@chmod +x scripts/generate-proto.sh
	@./scripts/generate-proto.sh

.PHONY: proto-lint
proto-lint: ## Lint protobuf files
	buf lint
