.PHONY: help setup start stop test-backend test-bff test-frontend

help: ## Mostra esta ajuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Inicializa infraestrutura (MongoDB, Valkey)
	docker-compose up -d mongodb valkey
	@echo "Aguardando serviços iniciarem..."
	@sleep 5
	@echo "MongoDB e Valkey prontos!"

start: ## Inicia todos os serviços
	docker-compose up -d

stop: ## Para todos os serviços
	docker-compose down

logs: ## Mostra logs de todos os serviços
	docker-compose logs -f

# Backend commands
test-backend: ## Executa testes do backend
	cd backend && go test -v -race -coverprofile=coverage.out ./...

coverage-backend: ## Mostra cobertura de testes do backend
	cd backend && go tool cover -html=coverage.out

run-backend: ## Executa backend localmente
	cd backend && go run ./cmd/server

build-backend: ## Compila backend
	cd backend && go build -o bin/server ./cmd/server

lint-backend: ## Executa linter no backend
	cd backend && golangci-lint run

# BFF commands
test-bff: ## Executa testes do BFF
	cd bff && pytest -v --cov=src --cov-report=term-missing

run-bff: ## Executa BFF localmente
	cd bff && uvicorn src.main:app --reload --host 0.0.0.0 --port 8000

lint-bff: ## Executa linter no BFF
	cd bff && ruff check src
	cd bff && black --check src

format-bff: ## Formata código do BFF
	cd bff && black src
	cd bff && ruff check --fix src

# Frontend commands
test-frontend: ## Executa testes do frontend
	cd frontend && npm test

run-frontend: ## Executa frontend localmente
	cd frontend && npm run dev

build-frontend: ## Compila frontend
	cd frontend && npm run build

lint-frontend: ## Executa linter no frontend
	cd frontend && npm run lint

format-frontend: ## Formata código do frontend
	cd frontend && npm run format

# Database
db-shell: ## Abre shell do MongoDB
	docker-compose exec mongodb mongosh -u admin -p admin123 --authenticationDatabase admin agencia_viagem

valkey-cli: ## Abre CLI do Valkey
	docker-compose exec valkey valkey-cli
