.PHONY: all build install setup clean uninstall help dev setup-env

BINARY_NAME := mcc
INSTALL_PATH := $(HOME)/bin
SHELL_CONFIG := $(HOME)/.zshrc
MCC_CONFIG_DIR := $(HOME)/.mcc
CLAUDE_CONFIG_LINE := export CLAUDE_CONFIG_DIR="$$HOME/.mcc/current"

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .
	@echo "✓ Built $(BINARY_NAME)"

# Install binary to ~/bin with aliases
install: build
	@mkdir -p $(INSTALL_PATH)
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@ln -sf $(BINARY_NAME) $(INSTALL_PATH)/multicc
	@ln -sf $(BINARY_NAME) $(INSTALL_PATH)/multi-claude-code
	@echo "✓ Installed to $(INSTALL_PATH)/$(BINARY_NAME)"
	@echo "✓ Created aliases: multicc, multi-claude-code"

# Setup shell environment variable
setup-env:
	@echo "Setting up environment variables..."
	@if grep -q "CLAUDE_CONFIG_DIR" $(SHELL_CONFIG) 2>/dev/null; then \
		echo "✓ CLAUDE_CONFIG_DIR already configured in $(SHELL_CONFIG)"; \
	else \
		echo '' >> $(SHELL_CONFIG); \
		echo '# Claude Code multi-account support (mcc)' >> $(SHELL_CONFIG); \
		echo '$(CLAUDE_CONFIG_LINE)' >> $(SHELL_CONFIG); \
		echo "✓ Added CLAUDE_CONFIG_DIR to $(SHELL_CONFIG)"; \
	fi
	@if echo "$$PATH" | grep -q "$(INSTALL_PATH)"; then \
		echo "✓ $(INSTALL_PATH) already in PATH"; \
	elif grep -q 'HOME/bin' $(SHELL_CONFIG) 2>/dev/null; then \
		echo "✓ ~/bin already configured in $(SHELL_CONFIG)"; \
	else \
		echo 'export PATH="$$HOME/bin:$$PATH"' >> $(SHELL_CONFIG); \
		echo "✓ Added ~/bin to PATH in $(SHELL_CONFIG)"; \
	fi

# Full setup: build + install + setup environment
setup: install setup-env
	@echo ""
	@echo "=========================================="
	@echo "✓ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Run: source $(SHELL_CONFIG)"
	@echo "  2. Run: mcc"
	@echo "=========================================="

# Quick dev iteration: just build and copy
dev: build
	@cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@ln -sf $(BINARY_NAME) $(INSTALL_PATH)/multicc
	@ln -sf $(BINARY_NAME) $(INSTALL_PATH)/multi-claude-code
	@echo "✓ Updated $(INSTALL_PATH)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@echo "✓ Cleaned"

# Uninstall everything
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@rm -f $(INSTALL_PATH)/multicc
	@rm -f $(INSTALL_PATH)/multi-claude-code
	@echo "✓ Removed $(INSTALL_PATH)/$(BINARY_NAME) and aliases"
	@echo ""
	@echo "Note: $(MCC_CONFIG_DIR) and shell config not removed."
	@echo "To fully clean up:"
	@echo "  1. Remove 'export CLAUDE_CONFIG_DIR=...' from $(SHELL_CONFIG)"
	@echo "  2. rm -rf $(MCC_CONFIG_DIR)"

# Show help
help:
	@echo "Claude Code Multi-Account Manager (mcc) - Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build      - Build the binary"
	@echo "  make install    - Build and install to $(INSTALL_PATH)"
	@echo "  make setup-env  - Add CLAUDE_CONFIG_DIR to shell config"
	@echo "  make setup      - Full setup (build + install + env)"
	@echo "  make dev        - Quick rebuild and install (for development)"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make uninstall  - Remove installed binary"
	@echo "  make help       - Show this help"
