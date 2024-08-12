# Makefile

# Project name
PROJECT_NAME=task-manager-gallatin
# Source directory
SRC_DIR=./cmd
# Build directory
BUILD_DIR=./build

# Target binary paths
LINUX_BINARY=$(BUILD_DIR)/$(PROJECT_NAME)_linux_amd64
MAC_BINARY=$(BUILD_DIR)/$(PROJECT_NAME)_darwin_amd64
WINDOWS_BINARY=$(BUILD_DIR)/$(PROJECT_NAME)_windows_amd64.exe

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
build-all: build-linux build-mac build-windows

# Docker build using local Linux binary
docker-build: build-linux
	docker build -t $(PROJECT_NAME)-service .

# Clean up build directory
clean:
	rm -rf $(BUILD_DIR)

# Phony targets to avoid conflicts with file names
.PHONY: build-linux build-mac build-windows build-all docker-build clean
