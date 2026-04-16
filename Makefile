APP_NAME := git-remote-commits
BUILD_DIR := bin
RELEASE_DIR := dist
SHORTCUT_DIR_WIN := C:/Bin/git-remote-commits
SHORTCUT_EXE_WIN := $(SHORTCUT_DIR_WIN)/$(APP_NAME).exe
VERSION ?= dev
LDFLAGS := -X 'main.Version=$(VERSION)'

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

.PHONY: test build release shortcut clean

test:
	@echo "Running tests..."
	@go test ./...

build:
	@echo "Building $(APP_NAME) for $(CURRENT_OS)..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)$(EXE_EXT) .

release:
	@echo "Building release artifacts..."
	@mkdir -p $(RELEASE_DIR)
	@GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(RELEASE_DIR)/$(APP_NAME)-linux-amd64 .
	@GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(RELEASE_DIR)/$(APP_NAME)-linux-arm64 .
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(RELEASE_DIR)/$(APP_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(RELEASE_DIR)/$(APP_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(RELEASE_DIR)/$(APP_NAME)-windows-amd64.exe .
	@GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(RELEASE_DIR)/$(APP_NAME)-windows-arm64.exe .

shortcut:
ifeq ($(OS),Windows_NT)
	@echo "Building Windows shortcut binary at $(SHORTCUT_EXE_WIN)..."
	@powershell -NoProfile -Command "New-Item -ItemType Directory -Force '$(SHORTCUT_DIR_WIN)' | Out-Null"
	@go build -ldflags "$(LDFLAGS)" -o $(SHORTCUT_EXE_WIN) .
	@echo "Done: $(SHORTCUT_EXE_WIN)"
else
	@echo "shortcut target is supported only on Windows."
	@exit 1
endif

clean:
	@echo "Cleaning build and release directories..."
	@$(call RM_DIR,$(BUILD_DIR))
	@$(call RM_DIR,$(RELEASE_DIR))
