.PHONY: help setup build build-docker run run-docker test test-api clean fmt lint docker-up docker-down docker-logs db-init db-reset

# Variables
APP_NAME=backend-ping-pong-app
BINARY_NAME=backend-ping-pong-app
GO=go
DOCKER=docker
DOCKER_COMPOSE=docker-compose

help:
	@echo "╔════════════════════════════════════════════════════════╗"
	@echo "║   Ping Pong Backend - Makefile                         ║"
	@echo "╚════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "📋 Local Development Commands:"
	@echo "  make setup         - Download dependencies"
	@echo "  make build         - Build executable locally"
	@echo "  make run           - Run server locally (requires PostgreSQL)"
	@echo "  make fmt           - Format code with gofmt"
	@echo "  make lint          - Run golangci-lint (if installed)"
	@echo "  make test          - Run unit tests"
	@echo "  make test-api      - Test API endpoints"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start services with Docker Compose"
	@echo "  make docker-down   - Stop all services"
	@echo "  make docker-logs   - View server logs"
	@echo "  make docker-clean  - Remove containers and volumes"
	@echo ""
	@echo "🗄️  Database Commands:"
	@echo "  make db-init       - Initialize database with docker"
	@echo "  make db-reset      - Drop and recreate database"
	@echo ""
	@echo "🧹 Cleanup:"
	@echo "  make clean         - Remove build artifacts"
	@echo ""

# ==================== Local Development ====================

setup:
	@echo "📦 Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "✅ Dependencies downloaded"

build:
	@echo "🔨 Building executable..."
	$(GO) build -o ./bin/$(BINARY_NAME) .
	@echo "✅ Built: ./bin/$(BINARY_NAME)"

run: build
	@echo "🚀 Starting server (localhost:8080)..."
	./bin/$(BINARY_NAME)

fmt:
	@echo "📝 Formatting code..."
	$(GO) fmt ./...
	@echo "✅ Code formatted"

lint:
	@echo "🔍 Running linter..."
	golangci-lint run ./...
	@echo "✅ Lint completed"

test:
	@echo "🧪 Running tests..."
	$(GO) test -v ./...
	@echo "✅ Tests completed"

test-api:
	@echo "🌐 Testing API endpoints..."
	@echo "Checking health endpoint..."
	curl -s http://localhost:8080/api/v1/health | jq .
	@echo "✅ API is responding"

clean:
	@echo "🧹 Cleaning up..."
	rm -rf bin/
	$(GO) clean
	@echo "✅ Cleaned"

# ==================== Docker ====================

docker-build:
	@echo "🐳 Building Docker image..."
	$(DOCKER) build -t $(APP_NAME):latest .
	@echo "✅ Docker image built: $(APP_NAME):latest"

docker-up:
	@echo "🚀 Starting Docker services..."
	$(DOCKER_COMPOSE) up -d
	@echo "⏳ Waiting for services to be ready..."
	@sleep 5
	@echo "✅ Services started!"
	@echo "   Backend: http://localhost:8080"
	@echo "   Database: localhost:5432"

docker-down:
	@echo "🛑 Stopping Docker services..."
	$(DOCKER_COMPOSE) down
	@echo "✅ Services stopped"

docker-logs:
	@echo "📋 Backend logs:"
	$(DOCKER_COMPOSE) logs -f backend

docker-logs-db:
	@echo "📋 Database logs:"
	$(DOCKER_COMPOSE) logs -f postgres

docker-clean:
	@echo "🧹 Removing Docker containers and volumes..."
	$(DOCKER_COMPOSE) down -v
	@echo "✅ Cleaned"

docker-restart:
	@echo "🔄 Restarting services..."
	$(DOCKER_COMPOSE) restart
	@echo "✅ Restarted"

# ==================== Database ====================

db-init:
	@echo "🗄️  Initializing database..."
	$(DOCKER_COMPOSE) up -d postgres
	@sleep 5
	@echo "⏳ Database initializing..."
	$(DOCKER_COMPOSE) exec postgres psql -U pingpong_user -d pingpong -f /docker-entrypoint-initdb.d/init.sql
	@echo "✅ Database initialized"

db-reset:
	@echo "⚠️  Resetting database..."
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) up -d postgres
	@sleep 5
	@echo "✅ Database reset"

# ==================== Combined ====================

dev: docker-up
	@echo "✅ Development environment ready!"
	@echo "Run 'make test-api' to test the API"

prod:
	@echo "📦 Building for production..."
	$(DOCKER) build -t $(APP_NAME):prod .
	@echo "✅ Production build ready"

all: setup build
	@echo "✅ All done!"

run: setup
	@echo "🚀 Starting server..."
	go run main.go handlers.go

fmt:
	@echo "📝 Formatting code..."
	go fmt ./...
	@echo "✅ Done"

test:
	@echo "🩺 Testing API health..."
	@curl -s http://localhost:8080/health | head -20 || echo "API not running. Start with: make run"

clean:
	@echo "🧹 Cleaning..."
	rm -f avatar-api
	@echo "✅ Done"

docker-up:
	@echo "🐳 Starting Docker containers..."
	docker-compose up -d
	@sleep 3
	@echo "✅ Services started!"
	@echo "📍 API: http://localhost:8080"
	@echo "📊 PostgreSQL: localhost:5432"
	@make test

docker-down:
	@echo "🛑 Stopping Docker containers..."
	docker-compose down
	@echo "✅ Stopped"

docker-logs:
	@echo "📋 API Logs (Ctrl+C to exit):"
	docker-compose logs -f api

docker-clean:
	@echo "🧹 Removing containers and volumes..."
	docker-compose down -v
	rm -rf uploads/
	@echo "✅ Cleaned"
