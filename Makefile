# =========================
# Makefile — blocd (daemon)
# =========================

APP_NAME := blocd
CMD_PATH := ./cmd/blocd
BIN_DIR  := bin

GOOS ?= $(shell go env GOOS)

# -------------------------
# OS-specific settings
# -------------------------
ifeq ($(GOOS),windows)
	EXT := .exe
	INSTALL_SUPPORTED := false
else ifeq ($(GOOS),darwin)
	EXT :=
	INSTALL_SUPPORTED := false
else
	# linux
	EXT :=
	INSTALL_SUPPORTED := true
	INSTALL_DIR := /usr/local/bin
endif

BIN := $(BIN_DIR)/$(APP_NAME)$(EXT)

.PHONY: build run install clean

# -------------------------
# Build
# -------------------------
build:
	@echo "Building for $(GOOS)..."
	go build -o "$(BIN)" "$(CMD_PATH)"

# -------------------------
# Run (foreground)
# -------------------------
run: build
	$(BIN)

# -------------------------
# Install (Linux only)
# -------------------------
install:
ifeq ($(INSTALL_SUPPORTED),true)
	@echo "Installing $(APP_NAME) to $(INSTALL_DIR)"
	sudo cp "$(BIN)" "$(INSTALL_DIR)/$(APP_NAME)"
else
	@echo "Install is not supported on $(GOOS). Build output is: $(BIN)"
endif

# -------------------------
# Clean
# -------------------------
clean:
	rm -rf "$(BIN_DIR)"
