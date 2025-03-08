APP_NAME := kubero
ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BINARY_NAME := $(ROOT_DIR)$(APP_NAME)
CMD_DIR := $(ROOT_DIR)cmd
INSTALL_SCRIPT=$(ROOT_DIR)scripts/install.sh
ARGS :=

# Colors
COLOR_RESET := \033[0m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_RED := \033[31m
COLOR_BLUE := \033[34m

# Logging Functions
log = @printf "$(COLOR_BLUE)[LOG]$(COLOR_RESET) %s\n" "$(1)"
success = @printf "$(COLOR_GREEN)[SUCCESS]$(COLOR_RESET) %s\n" "$(1)"
warning = @printf "$(COLOR_YELLOW)[WARNING]$(COLOR_RESET) %s\n" "$(1)"
break = @printf "$(COLOR_BLUE)[LOG]$(COLOR_RESET)\n"
error = @printf "$(COLOR_RED)[ERROR]$(COLOR_RESET) %s\n" "$(1)" && exit 1

# Build the binary using the install script
build:
	@if [ -f $(BINARY_NAME) ]; then rm $(BINARY_NAME); fi
	$(call log, Building $(BINARY_NAME) )
	$(call break, b )
	@go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o $(BINARY_NAME) ${CMD_DIR} || exit 1
	$(call break, b )
	$(call success, Build process completed successfully)
	$(call break, b )
	$(call log, Compressing $(BINARY_NAME)... This may take a while)
	$(call break, b )
	@upx $(BINARY_NAME) --force-overwrite --lzma --no-progress --no-color -qqq || exit 1
	$(call break, b )
	$(call success, Compressed $(BINARY_NAME) )


# Build the binary using the install script
build-dev:
	@if [ -f $(BINARY_NAME) ]; then rm $(BINARY_NAME); fi
	$(call log, Building $(BINARY_NAME) )
	$(call break, b )
	@go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o $(BINARY_NAME) ${CMD_DIR} || exit 1
	$(call break, b )
	$(call success, Built $(BINARY_NAME) )

# Install the binary and configure environment
install:
	$(call log, Installing $(BINARY_NAME) )
	$(call break, b )
	@sh $(INSTALL_SCRIPT) install $(ARGS) || exit 1
	$(call break, b )
	$(call success, Installation completed )

# Clean up build artifacts
clean:
	$(call log, Cleaning up build artifacts)
	$(call break, b )
	@if [ -f $(BINARY_NAME) ]; then rm $(BINARY_NAME); fi
	$(call break, b )
	$(call success, Cleaned up build artifacts)

# Display this help message
help:
	$(call log, $(APP_NAME) Makefile )
	$(call break, b )
	$(call log, Usage: )
	$(call log,   make [target] [ARGS='--custom-arg value'] )
	$(call break, b )
	$(call log, Available targets: )
	$(call log,   make build      - Build the binary using install script)
	$(call log,   make build-dev  - Build the binary without compressing it)
	$(call log,   make install    - Install the binary and configure environment)
	$(call log,   make clean      - Clean up build artifacts)
	$(call log,   make help       - Display this help message)
	$(call break, b )
	$(call log, Usage with arguments: )
	$(call log,   make install ARGS='--custom-arg value' - Pass custom arguments to the install script)
	$(call break, b )
	$(call log, Example: )
	$(call log,   make install ARGS='--prefix /usr/local')
	$(call break, b )
	$(call log, $(APP_NAME) is a tool for managing Kubernetes resources)
	$(call break, b )
	$(call log, For more information, visit: )
	$(call log, 'https://github.com/kubero-dev/kubero-cli' )
	$(call break, b )
	$(call success, End of help message)
