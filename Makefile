# Makefile for simple-multiplayer-service

# Variables
APP_NAME := simple-multiplayer-service
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG := latest
DOCKER_FULL_NAME := $(DOCKER_IMAGE):$(DOCKER_TAG)
SERVER_PORT := 8080
SESSION_LIMIT := 10

# Go related variables
GO_CMD := go
GO_BUILD := $(GO_CMD) build
GO_TEST := $(GO_CMD) test
GO_RUN := $(GO_CMD) run
GO_CLEAN := $(GO_CMD) clean
GO_SERVER_PATH := ./cmd/server

# Docker related variables
DOCKER_CMD := docker
DOCKER_COMPOSE_CMD := docker-compose

# Default target
.PHONY: all
all: help

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build-docker    - Build Docker image"
	@echo "  run-local       - Run the application locally (outside Docker)"
	@echo "  run-docker      - Run the application inside Docker"
	@echo "  test-local      - Run tests locally (outside Docker)"
	@echo "  test-docker     - Run tests inside Docker"
	@echo "  clean           - Clean build artifacts"
	@echo "  help            - Show this help message"

# Build Docker image
.PHONY: build-docker
build-docker:
	@echo "Building Docker image $(DOCKER_FULL_NAME)..."
	$(DOCKER_CMD) build -t $(DOCKER_FULL_NAME) .

# Run the application locally
.PHONY: run-local
run-local:
	@echo "Running application locally on port $(SERVER_PORT)..."
	SESSION_LIMIT=$(SESSION_LIMIT) $(GO_RUN) $(GO_SERVER_PATH)

# Run the application inside Docker
.PHONY: run-docker
run-docker: build-docker
	@echo "Running application in Docker on port $(SERVER_PORT)..."
	$(DOCKER_CMD) run -p $(SERVER_PORT):$(SERVER_PORT) -e SESSION_LIMIT=$(SESSION_LIMIT) $(DOCKER_FULL_NAME)

# Alternative: Run using docker-compose
.PHONY: run-docker-compose
run-docker-compose:
	@echo "Running application using docker-compose..."
	$(DOCKER_COMPOSE_CMD) up --build

# Run tests locally
.PHONY: test-local
test-local:
	@echo "Running tests locally..."
	$(GO_TEST) ./...

# Run tests inside Docker
.PHONY: test-docker
test-docker:
	@echo "Running tests inside Docker..."
	$(DOCKER_COMPOSE_CMD) -f docker-compose-test.yml up --abort-on-container-exit

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GO_CLEAN)
	rm -f $(APP_NAME)
	@echo "Removing Docker image..."
	$(DOCKER_CMD) rmi $(DOCKER_FULL_NAME) || true
	@echo "Removing Docker test containers..."
	$(DOCKER_COMPOSE_CMD) -f docker-compose-test.yml down --rmi all 2>/dev/null || true
