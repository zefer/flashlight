.PHONY: build build-arm test clean deploy status help

# Configuration
BINARY_NAME := flashlight
SERVER_HOST ?= livingroom
SERVER_USER ?= joe
INSTALL_PATH := /usr/bin/$(BINARY_NAME)
SERVICE_NAME := flashlight

# Build for local architecture (for development/testing compilation)
build:
	@echo "Building for local architecture..."
	go build -o $(BINARY_NAME)

# Cross-compile for Raspberry Pi 2/Zero (ARMv7)
build-arm:
	@echo "Cross-compiling for Raspberry Pi (ARMv7)..."
	GOOS=linux GOARCH=arm GOARM=7 go build -o $(BINARY_NAME)
	@echo "Build complete: $(BINARY_NAME)"
	@file $(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)

# Deploy to Raspberry Pi
deploy: build-arm
	@echo "Deploying to $(SERVER_USER)@$(SERVER_HOST)..."
	@if [ ! -f "$(BINARY_NAME)" ]; then \
		echo "Error: Binary not found. Run 'make build-arm' first."; \
		exit 1; \
	fi
	@echo "Uploading binary..."
	scp $(BINARY_NAME) $(SERVER_USER)@$(SERVER_HOST):/home/$(SERVER_USER)/
	@echo "Installing and restarting service..."
	ssh $(SERVER_USER)@$(SERVER_HOST) -t '\
		sudo systemctl stop $(SERVICE_NAME) && \
		sleep 1 && \
		sudo mv /home/$(SERVER_USER)/$(BINARY_NAME) $(INSTALL_PATH) && \
		sudo systemctl start $(SERVICE_NAME) && \
		sleep 1 && \
		sudo systemctl status $(SERVICE_NAME)'
	@echo "Deployment complete!"

# Check service status on the Pi
status:
	@echo "Checking service status on $(SERVER_USER)@$(SERVER_HOST)..."
	ssh $(SERVER_USER)@$(SERVER_HOST) -t 'sudo systemctl status $(SERVICE_NAME)'

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build for local architecture"
	@echo "  build-arm  - Cross-compile for Raspberry Pi (ARMv7)"
	@echo "  test       - Run tests"
	@echo "  clean      - Remove build artifacts"
	@echo "  deploy     - Build and deploy to Raspberry Pi (default: $(SERVER_HOST))"
	@echo "  status     - Check service status on Raspberry Pi"
	@echo "  help       - Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  SERVER_HOST - Deployment target (default: music)"
	@echo "  SERVER_USER - SSH user (default: joe)"
	@echo ""
	@echo "Examples:"
	@echo "  make build-arm              - Cross-compile for Pi"
	@echo "  make deploy                 - Deploy to default server"
	@echo "  SERVER_HOST=pi make deploy  - Deploy to custom server"
