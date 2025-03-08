#!/usr/bin/env bash

_about='
################################################################################
# This Script is used to install kubero-cli binaries.                          #
#                                                                              #
# Supported OS: Linux, macOS ---> Windows(not supported)                       #
# Supported Architecture: amd64, arm64                                         #
# Source: https://github.com/kubero-dev/kubero-cli                             #
# Binary Release: https://github.com/kubero-dev/kubero-cli/releases/latest     #
# License: Apache License 2.0                                                  #
# Usage:                                                                       #
#   curl -fsSL get.kubero.dev | bash                                           #
#   curl -fsSL get.kubero.dev | bash -s -- [version]                           #
#   bash <(curl -fsSL get.kubero.dev) [version]                                #
#                                                                              #
# Notes:                                                                       #
# - [version] is optional; if omitted, the latest version will be used.        #
# - If the script is run locally, it will try to resolve the version from the  #
#   repo tags if no version is provided.                                       #
# - The script will install the binary in the ~/.local/bin directory if the    #
#   user is not root. Otherwise, it will install in /usr/local/bin.            #
# - The script will add the installation directory to the PATH in the shell    #
#   configuration file.                                                        #
# - The script will also install UPX if it is not already installed.           #
# - The script will build the binary if the build option is provided.          #
# - The script will download the binary from the release URL if the install    #
#   option is provided.                                                        #
# - The script will clean up build artifacts if the clean option is provided.  #
# - The script will check if the required dependencies are installed.          #
# - The script will validate the Go version before building the binary.        #
# - The script will check if the installation directory is in the PATH.        #
# - The script will print a summary of the installation.                       #
################################################################################
'

set -e

# Define variables
APP_NAME="$(basename "$(dirname "$(realpath "$0")")")"
CMD_PATH="$(dirname "$(realpath "$(dirname "$0")")")/cmd"
BUILD_PATH="$(dirname "$CMD_PATH")"
BINARY="$BUILD_PATH/$APP_NAME"
LOCAL_BIN="$HOME/.local/bin"
GLOBAL_BIN="/usr/local/bin"
SUCCESS="\033[0;32m"
WARN="\033[0;33m"
ERROR="\033[0;31m"
INFO="\033[0;36m"
NC="\033[0m"

# Log messages with different levels
# Arguments:
#   $1 - log level (info, warn, error, success)
#   $2 - message to log
log() {
  local type=
  type=${1:-info}
  local message=
  message=${2:-}

  # With colors
  case $type in
    info|INFO|-i|-I)
      printf '%b[INFO]%b ℹ️  %s\n' "$INFO" "$NC" "$message"
      ;;
    warn|WARN|-w|-W)
      printf '%b[WARN]%b ⚠️  %s\n' "$WARN" "$NC" "$message"
      ;;
    error|ERROR|-e|-E)
      printf '%b[ERROR]%b ❌  %s\n' "$ERROR" "$NC" "$message"
      ;;
    success|SUCCESS|-s|-S)
      printf '%b[SUCCESS]%b ✅  %s\n' "$SUCCESS" "$NC" "$message"
      ;;
    *)
      log "info" "$message"
      ;;
  esac
}

# Detect the operating system
# Returns:
#   OS name (linux, darwin, unsupported)
get_os() {
    case "$(uname -s)" in
        Linux*) echo "linux" ;;
        Darwin*) echo "darwin" ;;
        *) echo "unsupported" ;;
    esac
}

# Detect the architecture
# Returns:
#   Architecture name (amd64, arm64, unsupported)
get_arch() {
    case "$(uname -m)" in
        x86_64) echo "amd64" ;;
        arm*|aarch64) echo "arm64" ;;
        *) echo "unsupported" ;;
    esac
}

# Detect the shell configuration file
# Returns:
#   Shell configuration file path
detect_shell_rc() {
    shell_rc_file=""
    user_shell=$(basename "$SHELL")
    case "$user_shell" in
        bash) shell_rc_file="$HOME/.bashrc" ;;
        zsh) shell_rc_file="$HOME/.zshrc" ;;
        sh) shell_rc_file="$HOME/.profile" ;;
        fish) shell_rc_file="$HOME/.config/fish/config.fish" ;;
        *)
            log "warn" "Unsupported shell" "$($user_shell), modify PATH manually."
            return 1
            ;;
    esac
    log "info" "$shell_rc_file"
}

# Add a directory to the PATH in the shell configuration file
# Arguments:
#   $1 - target path to add to PATH
add_to_path() {
    target_path="$1"
    shell_rc_file=$(detect_shell_rc)
    if [ -z "$shell_rc_file" ]; then
        log "error" "Could not determine shell configuration file."
        return 1
    fi

    if grep -q "export PATH=.*$target_path" "$shell_rc_file" 2>/dev/null; then
        log "success" "$target_path is already in $shell_rc_file."
        return 0
    fi

    log "success" "export PATH=$target_path:\$PATH" >> "$shell_rc_file"
    log "success" "Added $target_path to PATH in $shell_rc_file."
    log "success" "Run 'source $shell_rc_file' to apply changes."
}

# Clean up build artifacts
clean() {
    log "info" "Cleaning up build artifacts..."
    rm -f "$BINARY" || true
    log "success" "Cleaned up build artifacts."
}

# Install the binary to the appropriate directory
install_binary() {
    if [ "$(id -u)" -ne 0 ]; then
        log "info" "You are not root. Installing in $LOCAL_BIN..."
        mkdir -p "$LOCAL_BIN"
        cp "$BINARY" "$LOCAL_BIN/$APP_NAME" || exit 1
        add_to_path "$LOCAL_BIN"
    else
        log "info" "Root detected. Installing in $GLOBAL_BIN..."
        cp "$BINARY" "$GLOBAL_BIN/$APP_NAME" || exit 1
        add_to_path "$GLOBAL_BIN"
    fi
    clean
}

# Install UPX if it is not already installed
install_upx() {
    if ! command -v upx > /dev/null; then
        log "info" "Installing UPX..."
        if [ "$(uname)" = "Darwin" ]; then
            brew install upx
        elif command -v apt-get > /dev/null; then
            sudo apt-get install -y upx
        else
            log "error" 'Install UPX manually from https://upx.github.io/'
            exit 1
        fi
    else
        log "success" ' UPX is already installed.'
    fi
}

# Check if the required dependencies are installed
# Arguments:
#   $@ - list of dependencies to check
check_dependencies() {
    for dep in "$@"; do
        if ! command -v "$dep" > /dev/null; then
            log "error" "$dep is not installed."
            exit 1
        else
            log "success" "$dep is installed."
        fi
    done
}

# Build the binary
build_binary() {
    log "info" "Building the binary..."
    go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o "$BINARY" "$CMD_PATH"
    install_upx
    upx "$BINARY" --force-overwrite --lzma --no-progress --no-color -qqq
}

# Validate the Go version
validate_versions() {
    REQUIRED_GO_VERSION="1.18"
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_GO_VERSION" ]; then
        log "error" "Go version must be >= $REQUIRED_GO_VERSION. Detected: $GO_VERSION"
        exit 1
    fi
    log "success" "Go version is valid: $GO_VERSION"
}

# Print a summary of the installation
summary() {
    install_dir="$BINARY"
    log "success" "Build and installation complete!"
    log "success" "Binary: $BINARY"
    log "success" "Installed in: $install_dir"
    check_path "$install_dir"
}

# Build the binary and validate the Go version
build_and_validate() {
    validate_versions
    build_binary
}

# Check if the installation directory is in the PATH
# Arguments:
#   $1 - installation directory
check_path() {
    log "info" "Checking if the installation directory is in the PATH..."
    if ! echo "$PATH" | grep -q "$1"; then
        log "warn" "$1 is not in the PATH."
        log "warn" "Add the following to your ~/.bashrc, ~/.zshrc, or equivalent file:"
        log "warn" "export PATH=$1:\$PATH"
    else
        log "success" "$1 is already in the PATH."
    fi
}

# Download the binary from the release URL
download_binary() {
    os=$(get_os)
    arch=$(get_arch)

    if [ "$os" == "unsupported" ] || [ "$arch" == "unsupported" ]; then
        log "error" "Unsupported OS or architecture."
        exit 1
    fi

    version=$(curl -s https://github.com/kubero-dev/kubero-cli/releases/latest | grep "tag_name" | cut -d '"' -f 4) || return 1
    release_url="https://github.com/kubero-dev/kubero-cli/releases/${version}/download/kubero-cli_${os}_${arch}.tar.gz" || return 1

    log "info" "Downloading binary for $os/$arch (version $version)..."
    if curl -L -o "${APP_NAME}.tar.gz" "$release_url"; then
        log "error" "Failed to download binary."
        exit 1
    fi

    log "info" "Extracting binary..."
    tar -xzf "${APP_NAME}.tar.gz" -C "$(dirname "$BINARY")"
    rm -f "${APP_NAME}.tar.gz"

    log "success" "Download and extraction complete!"
}

# Install the binary from the release URL
install_from_release() {
    download_binary
    install_binary
}

# Clear the screen before beginning
printf '\33c\e[3J'

# Print the about message
printf "\n%s\n\n" "$_about"

# Check if the user has provided a command
case "$1" in
    build)
        build_and_validate
        ;;
    install)
        log info "Do you want to build locally or download the precompiled binary? (build/download)"
        read -r choice
        if [ "$choice" == "download" ]; then
            install_from_release
        else
            build_and_validate
            install_binary
        fi
        summary
        ;;
    clean)
        clean
        ;;
    *)
        echo "Usage: $0 {build|install|clean}"
        exit 1
        ;;
esac