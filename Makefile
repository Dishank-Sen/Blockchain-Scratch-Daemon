# =========================
# Makefile — blocd (daemon)
# =========================

APP_NAME=blocd
CMD_PATH=./cmd/blocd
BIN_DIR=bin

GOOS ?= $(shell go env GOOS)

ifeq ($(GOOS),windows)
	EXT=.exe
	INSTALL_DIR=C:/Program Files/blocd
else
	EXT=
	INSTALL_DIR=/usr/local/bin
endif

BIN=$(BIN_DIR)/$(APP_NAME)$(EXT)

.PHONY: build run install clean

# -------------------------
# Build
# -------------------------
build:
	go build -o "$(BIN)" "$(CMD_PATH)"

# -------------------------
# Run (foreground)
# -------------------------
run: build
	$(BIN)

# -------------------------
# Install (global)
# -------------------------
install: build
ifeq ($(GOOS),windows)
	@echo "Installing $(APP_NAME) (Administrator required)"
	powershell -Command " \
	  if (-not ([Security.Principal.WindowsPrincipal] \
	    [Security.Principal.WindowsIdentity]::GetCurrent() \
	    ).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) { \
	      Write-Error 'Run make install as Administrator'; exit 1 \
	  }"
	powershell -Command "New-Item -ItemType Directory -Force '$(INSTALL_DIR)'"
	powershell -Command "Copy-Item '$(BIN)' '$(INSTALL_DIR)/$(APP_NAME)$(EXT)' -Force"
	powershell -Command " \
	  $$path = [Environment]::GetEnvironmentVariable('Path','Machine'); \
	  if ($$path -notlike '*$(INSTALL_DIR)*') { \
	    [Environment]::SetEnvironmentVariable('Path', $$path + ';$(INSTALL_DIR)', 'Machine') \
	  }"
	@echo ""
	@echo "SUCCESS:"
	@echo "  Installed $(APP_NAME) to $(INSTALL_DIR)"
	@echo "  Restart all terminals or reboot for PATH changes to apply"
else
	sudo cp "$(BIN)" "$(INSTALL_DIR)/$(APP_NAME)"
endif

# -------------------------
# Clean
# -------------------------
clean:
	rm -rf "$(BIN_DIR)"
