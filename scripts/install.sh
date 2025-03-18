#!/usr/bin/env bash

set -e

_OWNER="kubero-dev"
_APP_NAME="kubero"
_PROJECT_NAME="kubero-cli"
_LICENSE="Apache License 2.0"

# Function to customize the release URL logic
get_release_url() {
    local version=$1
    local os=$2
    local arch=$3
    # Default logic for constructing the release URL
    echo "https://github.com/${_OWNER}/${_PROJECT_NAME}/releases/download/${version}/${_PROJECT_NAME}_${os}_${arch}.tar.gz"
}

_ABOUT="'
################################################################################
  This Script is used to install logz project.

  Supported OS: Linux, macOS ---> Windows(not supported)
  Supported Architecture: amd64, arm64
  Source: https://github.com/${_OWNER}/${_PROJECT_NAME}
  Binary Release: https://github.com/${_OWNER}/${_PROJECT_NAME}/releases/latest
  License: ${_LICENSE}
  Notes:
    - [version] is optional; if omitted, the latest version will be used.
    - If the script is run locally, it will try to resolve the version from the
      repo tags if no version is provided.
    - The script will install the binary in the ~/.local/bin directory if the
      user is not root. Otherwise, it will install in /usr/local/bin.
    - The script will add the installation directory to the PATH in the shell
      configuration file.
    - The script will also install UPX if it is not already installed.
    - The script will build the binary if the build option is provided.
    - The script will download the binary from the release URL if the install
      option is provided.
    - The script will clean up build artifacts if the clean option is provided.
    - The script will check if the required dependencies are installed.
    - The script will validate the Go version before building the binary.
    - The script will check if the installation directory is in the PATH.
    - The script will print a summary of the installation.
################################################################################
'"

_CMD_PATH="$(dirname "$(realpath "$(dirname "$0")")")/cmd"
_BUILD_PATH="$(dirname "${_CMD_PATH}")"
_BINARY="${_BUILD_PATH}/${_APP_NAME}"
_LOCAL_BIN="${HOME:-"~"}/.local/bin"
_GLOBAL_BIN="/usr/local/bin"
_SUCCESS="\033[0;32m"
_WARN="\033[0;33m"
_ERROR="\033[0;31m"
_INFO="\033[0;36m"
_NC="\033[0m"

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
    info|_INFO|-i|-I)
      printf '%b[_INFO]%b ℹ️  %s\n' "$_INFO" "$_NC" "$message"
      ;;
    warn|_WARN|-w|-W)
      printf '%b[_WARN]%b ⚠️  %s\n' "$_WARN" "$_NC" "$message"
      ;;
    error|_ERROR|-e|-E)
      printf '%b[_ERROR]%b ❌  %s\n' "$_ERROR" "$_NC" "$message"
      ;;
    success|_SUCCESS|-s|-S)
      printf '%b[_SUCCESS]%b ✅  %s\n' "$_SUCCESS" "$_NC" "$message"
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
            log "warn" "Unsupported shell, modify PATH manually."
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

    echo "export PATH=$target_path:\$PATH" >> "$shell_rc_file"
    log "success" "Added $target_path to PATH in $shell_rc_file."
    log "success" "Run 'source $shell_rc_file' to apply changes."
}

# Clean up build artifacts
clean() {
    log "info" "Cleaning up build artifacts..."
    rm -f "$_BINARY" || true
    log "success" "Cleaned up build artifacts."
}

# Install the binary to the appropriate directory
install_binary() {
    if [ "$(id -u)" -ne 0 ]; then
        log "info" "You are not root. Installing in $_LOCAL_BIN..."
        mkdir -p "$_LOCAL_BIN"
        cp "$_BINARY" "$_LOCAL_BIN/$_APP_NAME" || exit 1
        add_to_path "$_LOCAL_BIN"
    else
        log "info" "Root detected. Installing in $_GLOBAL_BIN..."
        cp "$_BINARY" "$_GLOBAL_BIN/$_APP_NAME" || exit 1
        add_to_path "$_GLOBAL_BIN"
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
    go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o "$_BINARY" "$_CMD_PATH"
    install_upx
    upx "$_BINARY" --force-overwrite --lzma --no-progress --no-color -qqq
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
    install_dir="$_BINARY"
    log "success" "Build and installation complete!"
    log "success" "Binary: $_BINARY"
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
    # Obtem o sistema operacional e a arquitetura
    os=$(get_os)
    arch=$(get_arch)

    # Validação: Verificar se o sistema operacional ou a arquitetura são suportados
    if [ "$os" = "unsupported" ] || [ "$arch" = "unsupported" ]; then
        log "error" "Unsupported OS or architecture: OS=$os ARCH=$arch."
        exit 1
    fi

    # Obter a versão mais recente de forma robusta (fallback para "latest")
    version=$(curl -s "https://api.github.com/repos/${_OWNER}/${_PROJECT_NAME}/releases/latest" | \
        grep "tag_name" | cut -d '"' -f 4 || echo "latest")

    if [ -z "$version" ]; then
        log "error" "Failed to determine the latest version."
        exit 1
    fi

    # Construir a URL de download usando a função customizável
    release_url=$(get_release_url "$version" "$os" "$arch")
    log "info" "Downloading ${_APP_NAME} binary for OS=$os, ARCH=$arch, Version=$version..."
    log "info" "Release URL: $release_url"

    # Diretório temporário para baixar o arquivo
    temp_dir=$(mktemp -d || exit 1)
    archive_path="${temp_dir}/${_APP_NAME}.tar.gz"

    # Realizar o download e validar sucesso
    if ! curl -L -o "$archive_path" "$release_url"; then
        log "error" "Failed to download the binary from: $release_url"
        rm -rf "$temp_dir"
        exit 1
    fi
    log "success" "Binary downloaded successfully."

    # Extração do arquivo para o diretório binário
    log "info" "Extracting binary to: $(dirname "$_BINARY")"
    if ! tar -xzf "$archive_path" -C "$(dirname "$_BINARY")"; then
        log "error" "Failed to extract the binary from: $archive_path"
        rm -rf "$temp_dir"
        exit 1
    fi

    # Limpar artefatos temporários
    rm -rf "$temp_dir"
    log "success" "Binary extracted successfully."

    # Verificar se o binário foi extraído com sucesso
    if [ ! -f "$_BINARY" ]; then
        log "error" "Binary not found after extraction: $_BINARY"
        exit 1
    fi

    log "success" "Download and extraction of ${_APP_NAME} completed!"
}

# Install the binary from the release URL
install_from_release() {
    download_binary
    install_binary
}

# Clear the screen before beginning
clear

check_dependencies "curl" "tar" "upx" "git" || exit 1

# Clear the screen before beginning
clear

# Print the ABOUT message
printf '\n%s\n\n' "${_ABOUT}"

# Check if the user has provided a command
case "$1" in
    build|BUILD|"-b"|"-B")
        log "info" "Executing build command..."
        build_and_validate || exit 1
        ;;
    install|INSTALL|"-i"|"-I")
        log "info" "Executing install command..."
        read -r -p "Do you want to download the precompiled binary? [y/N] (No will build locally): " c </dev/tty
        log "info" "User choice: ${c}"

        if [ "$c" = "y" ] || [ "$c" = "Y" ]; then
            log "info" "Downloading precompiled binary..."
            install_from_release || exit 1
        else
            log "info" "Building locally..."
            build_and_validate || exit 1
            install_binary || exit 1
        fi
        summary
        ;;
    clean|CLEAN|"-c"|"-C")
        log "info" "Executing clean command..."
        clean || exit 1
        ;;
    *)
        log "error" "Invalid command: $1"
        echo "Usage: $0 {build|install|clean}"
        exit 1
        ;;
esac

exit $?