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


set -eo pipefail
[[ $TRACE ]] && set -x

# Function to detect the operating system
get_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        #CYGWIN*|MINGW*) echo "windows" ;;
        *)          echo "unsupported" ;;
    esac
}

# Function to detect the architecture
get_arch() {
    case "$(uname -m)" in
        x86_64)    echo "amd64" ;;
        arm*|aarch64) echo "arm64" ;;
        *)          echo "unsupported" ;;
    esac
}

# Detect OS and architecture
os=$(get_os)
arch=$(get_arch)

if [ -z "$1" ]; then
    version="latest"
else
    version=$1
fi

if [ "$os" == "unsupported" ] || [ "$arch" == "unsupported" ]; then
    echo "Unsupported OS or architecture."
    exit 1
fi

# Check if there is already a binary installed
if [ -f "/usr/local/bin/kubero" ]; then
    echo "There is already a binary installed in /usr/local/bin/kubero ."
    read -p "Do you want to replace it? [y/n] " replaceBinary
    echo

    if [ "$replaceBinary" != "y" ] && [ "$replaceBinary" != "" ]; then
        echo "Aborting installation."
        rm -rf "$temp_dir"
        exit 1
    fi
fi

# Define the release URL
if [ $version == "latest" ]; then
    release_url="https://github.com/kubero-dev/kubero-cli/releases/latest/download/kubero-cli_${os}_${arch}.tar.gz"
else 
    release_url="https://github.com/kubero-dev/kubero-cli/releases/download/${version}/kubero-cli_${os}_${arch}.tar.gz"
fi

# Create a temporary directory
temp_dir=$(mktemp -d)

# Download and unpack the binary
echo "Downloading ${release_url} ..."
curl -L -s -o "${temp_dir}/kubero-cli.tar.gz" "$release_url"
if [ $? -ne 0 ]; then
    echo "Failed to download the binary."
    rm -rf "$temp_dir"
    exit 1
fi

echo "Unpacking the binary..."
tar -xzvf "${temp_dir}/kubero-cli.tar.gz" -C "$temp_dir"
if [ $? -ne 0 ]; then
    echo "Failed to unpack the binary."
    rm -rf "$temp_dir"
    exit 1
fi

# Check if the binary exists
if [ ! -f "${temp_dir}/kubero" ]; then
    echo "Failed to unpack the binary."
    rm -rf "$temp_dir"
    exit 1
fi

# Install the binary in /usr/local/bin
echo "Installing kubero in /usr/local/bin ..."
sudo mv "${temp_dir}/kubero" "/usr/local/bin/kubero"
if [ $? -ne 0 ]; then
    echo "Failed to install kubero."
    rm -rf "$temp_dir"
    exit 1
fi

# Clean up the temporary directory
rm -rf "$temp_dir"

echo 
echo "Kubero has been successfully installed."
echo "Run 'kubero install' to create a kubernetes cluster."