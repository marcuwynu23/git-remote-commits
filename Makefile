APP_NAME := commitview
BUILD_DIR := bin
RELEASE_DIR := dist

.PHONY: build release clean

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) .

release:
	@echo "Building release artifacts..."
	@mkdir -p $(RELEASE_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(RELEASE_DIR)/$(APP_NAME)-linux-amd64 .
	@GOOS=darwin GOARCH=amd64 go build -o $(RELEASE_DIR)/$(APP_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -o $(RELEASE_DIR)/$(APP_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 go build -o $(RELEASE_DIR)/$(APP_NAME)-windows-amd64.exe .

clean:
	@echo "Cleaning build and release directories..."
	@rm -rf $(BUILD_DIR) $(RELEASE_DIR)
