#!/usr/bin/env bash

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
#   curl -fsSL get.kubero.dev | bash -s -- v1.10.0                             #
#   bash <(curl -fsSL get.kubero.dev) v1.9.2                                   #
################################################################################

APP_NAME="kubero"
CMD_PATH="$(dirname "$(realpath "$(dirname "$0")")")/cmd"
BUILD_PATH="$(dirname "$CMD_PATH")"
BINARY="$BUILD_PATH/$APP_NAME"
LOCAL_BIN="$HOME/.local/bin"
GLOBAL_BIN="/usr/local/bin"

GREEN="\033[0;32m"
YELLOW="\033[0;33m"
RED="\033[0;31m"
NC="\033[0m"

set -e
# [[ $TRACE ]] && set -x
# [[ $TRACE ]] && set -x

get_os() {
    case "$(uname -s)" in
        Linux*) echo "linux" ;;
        Darwin*) echo "darwin" ;;
        *) echo "unsupported" ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64) echo "amd64" ;;
        arm*|aarch64) echo "arm64" ;;
        *) echo "unsupported" ;;
    esac
}

detect_shell_rc() {
    shell_rc_file=""
    user_shell=$(basename "$SHELL")
    case "$user_shell" in
        bash) shell_rc_file="$HOME/.bashrc" ;;
        zsh) shell_rc_file="$HOME/.zshrc" ;;
        sh) shell_rc_file="$HOME/.profile" ;;
        fish) shell_rc_file="$HOME/.config/fish/config.fish" ;;
        *)
            echo "${YELLOW}Warning: Unsupported shell ($user_shell). Modify PATH manually.${NC}"
            return 1
            ;;
    esac
    echo "$shell_rc_file"
}

add_to_path() {
    target_path="$1"
    shell_rc_file=$(detect_shell_rc)
    if [ -z "$shell_rc_file" ]; then
        echo "${RED}Error: Could not determine shell configuration file.${NC}"
        return 1
    fi

    if grep -q "export PATH=.*$target_path" "$shell_rc_file" 2>/dev/null; then
        echo "${GREEN}✅ $target_path is already in $shell_rc_file.${NC}"
        return 0
    fi

    echo "export PATH=$target_path:\$PATH" >> "$shell_rc_file"
    echo "${GREEN}Added $target_path to PATH in $shell_rc_file.${NC}"
    echo "Run 'source $shell_rc_file' to apply changes."
}

install_binary() {
    if [ "$(id -u)" -ne 0 ]; then
        echo "You are not root. Installing in $LOCAL_BIN..."
        mkdir -p "$LOCAL_BIN"
        cp "$BINARY" "$LOCAL_BIN/$APP_NAME" || exit 1
        add_to_path "$LOCAL_BIN"
    else
        echo "Root detected. Installing in $GLOBAL_BIN..."
        cp "$BINARY" "$GLOBAL_BIN/$APP_NAME" || exit 1
        add_to_path "$GLOBAL_BIN"
    fi
    clean
}

install_upx() {
    if ! command -v upx > /dev/null; then
        echo "${YELLOW}Installing UPX...${NC}"
        if [ "$(uname)" = "Darwin" ]; then
            brew install upx
        elif command -v apt-get > /dev/null; then
            sudo apt-get install -y upx
        else
            echo "${RED}Install UPX manually from https://upx.github.io/${NC}"
            exit 1
        fi
    else
        echo "${GREEN}✅ UPX is already installed.${NC}"
    fi
}

check_dependencies() {
    for dep in "$@"; do
        if ! command -v "$dep" > /dev/null; then
            echo "${RED}Error: $dep is not installed.${NC}"
            exit 1
        else
            echo "${GREEN}✅ $dep is installed.${NC}"
        fi
    done
}

build_binary() {
    echo "Building the binary..."
    go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o "$BINARY" "$CMD_PATH"
    install_upx
    upx "$BINARY" --force-overwrite --lzma --no-progress --no-color
}

clean() {
    echo "Cleaning up build artifacts..."
    rm -f "$BINARY"
}

validate_versions() {
    REQUIRED_GO_VERSION="1.18"
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_GO_VERSION" ]; then
        echo "${RED}Error: Go version must be >= $REQUIRED_GO_VERSION. Detected: $GO_VERSION${NC}"
        exit 1
    fi
    echo "${GREEN}✅ Go version is valid: $GO_VERSION${NC}"
}

summary() {
    install_dir="$BINARY"
    echo "${GREEN}Build and installation complete!${NC}"
    echo "Binary: $BINARY"
    echo "Installed in: $install_dir"
    check_path "$install_dir"
}

build_and_validate() {
    validate_versions
    build_binary
}

check_path() {
    echo "Checking if the installation directory is in the PATH..."
    if ! echo "$PATH" | grep -q "$1"; then
        echo "⚠️  Warning: $1 is not in the PATH."
        echo "Add the following to your ~/.bashrc, ~/.zshrc, or equivalent file:"
        echo "export PATH=$1:\$PATH"
    else
        echo "✅ $1 is already in the PATH."
    fi
}

download_binary() {
    os=$(get_os)
    arch=$(get_arch)

    if [ "$os" == "unsupported" ] || [ "$arch" == "unsupported" ]; then
        echo "${RED}Error: Unsupported OS or architecture.${NC}"
        exit 1
    fi

    version=$(curl -s https://github.com/kubero-dev/kubero-cli/releases/latest | grep "tag_name" | cut -d '"' -f 4)
    release_url="https://github.com/kubero-dev/kubero-cli/releases/${version}/download/kubero-cli_${os}_${arch}.tar.gz"

    echo "Downloading binary for $os/$arch (version $version)..."
    if curl -L -o "${APP_NAME}.tar.gz" "$release_url"; then
        echo "${RED}Error: Failed to download binary.${NC}"
        exit 1
    fi

    echo "Extracting binary..."
    tar -xzf "${APP_NAME}.tar.gz" -C "$(dirname "$BINARY")"
    rm -f "${APP_NAME}.tar.gz"

    echo "${GREEN}Download and extraction complete!${NC}"
}

install_from_release() {
    download_binary
    install_binary
}

case "$1" in
    build)
        build_and_validate
        ;;
    install)
        echo "Do you want to build locally or download the precompiled binary? (build/download)"
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