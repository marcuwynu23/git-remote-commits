APP_NAME := commitview
BUILD_DIR := bin
RELEASE_DIR := dist

ifeq ($(OS),Windows_NT)
	CURRENT_OS := windows
	EXE_EXT := .exe
	RM_DIR := powershell -NoProfile -Command "if (Test-Path '$(1)') { Remove-Item -Recurse -Force '$(1)' }"
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		CURRENT_OS := darwin
	else ifeq ($(UNAME_S),Linux)
		CURRENT_OS := linux
	else
		CURRENT_OS := unknown
	endif
	EXE_EXT :=
	RM_DIR := rm -rf $(1)
endif

.PHONY: build release clean

build:
	@echo "Building $(APP_NAME) for $(CURRENT_OS)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME)$(EXE_EXT) .

release:
	@echo "Building release artifacts..."
	@mkdir -p $(RELEASE_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(RELEASE_DIR)/$(APP_NAME)-linux-amd64 .
	@GOOS=darwin GOARCH=amd64 go build -o $(RELEASE_DIR)/$(APP_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -o $(RELEASE_DIR)/$(APP_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 go build -o $(RELEASE_DIR)/$(APP_NAME)-windows-amd64.exe .

clean:
	@echo "Cleaning build and release directories..."
	@$(call RM_DIR,$(BUILD_DIR))
	@$(call RM_DIR,$(RELEASE_DIR))
