# Makefile

# Project name
PROJECT_NAME=task-manager-gallatin
# Source directory
SRC_DIR=./cmd
# Build directory
BUILD_DIR=./build
# Docker image name
IMAGE_NAME=task-manager-gallatin
# Docker compose file
DOCKER_COMPOSE_FILE=docker-compose.yaml

# Target binary paths
LINUX_BINARY=$(BUILD_DIR)/$(PROJECT_NAME)_linux_amd64
MAC_BINARY=$(BUILD_DIR)/$(PROJECT_NAME)_darwin_amd64
WINDOWS_BINARY=$(BUILD_DIR)/$(PROJECT_NAME)_windows_amd64.exe

env:
	export APP_ENV=docker

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(LINUX_BINARY) $(SRC_DIR)/main.go

# Build for macOS
build-mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(MAC_BINARY) $(SRC_DIR)/main.go

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $(WINDOWS_BINARY) $(SRC_DIR)/main.go

# Build for all platforms
build-all: env build-linux build-mac build-windows

docker-remove:
	docker rmi $(IMAGE_NAME)

# Docker build using local Linux binary
docker-build: env build-linux
	docker build -t $(IMAGE_NAME) .

# Run docker-compose using the built Docker image
docker-compose-up: docker-build
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Stop the services
docker-compose-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Clean up build directory
clean:
	rm -rf $(BUILD_DIR)

# Phony targets to avoid conflicts with file names
.PHONY: env build-linux build-mac build-windows build-all docker-build docker-compose-up docker-compose-down clean
