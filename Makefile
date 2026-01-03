# =========================
# Makefile — blocd (daemon)
# =========================

APP_NAME=blocd
CMD_PATH=./cmd/blocd
BIN_DIR=bin

GOOS ?= $(shell go env GOOS)

ifeq ($(GOOS),windows)
	EXT=.exe
	INSTALL_DIR=$(USERPROFILE)\bin
else
	EXT=
	INSTALL_DIR=/usr/local/bin
endif

BIN=$(BIN_DIR)/$(APP_NAME)$(EXT)

.PHONY: build run install clean

build:
	go build -o $(BIN) $(CMD_PATH)

run: build
	$(BIN)

install: build
	@echo "Installing $(APP_NAME) to $(INSTALL_DIR)"
ifeq ($(GOOS),windows)
	@if not exist "$(INSTALL_DIR)" mkdir "$(INSTALL_DIR)"
	copy $(BIN) "$(INSTALL_DIR)"
else
	sudo cp $(BIN) $(INSTALL_DIR)/$(APP_NAME)
endif

clean:
	rm -rf $(BIN_DIR)
