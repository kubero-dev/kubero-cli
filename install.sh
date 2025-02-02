#!/usr/bin/env bash

set -eo pipefail
[[ $TRACE ]] && set -x

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

os=$(get_os)
arch=$(get_arch)
version=${1:-latest}

if [[ "$os" == "unsupported" || "$arch" == "unsupported" ]]; then
    echo "Unsupported OS or architecture."
    exit 1
fi

if [[ -f "/usr/local/bin/kubero" ]]; then
    read -r -p "Do you want to replace it? [y/n] " replaceBinary
    [[ "$replaceBinary" != "y" && "$replaceBinary" != "" ]] && echo "Aborting installation." && exit 1
fi

release_url="https://github.com/kubero-dev/kubero-cli/releases/${version}/download/kubero-cli_${os}_${arch}.tar.gz"
temp_dir=$(mktemp -d)

echo "Downloading ${release_url} ..."
curl -L -s -o "${temp_dir}/kubero-cli.tar.gz" "$release_url" || { echo "Failed to download the binary."; rm -rf "$temp_dir"; exit 1; }

echo "Unpacking the binary..."
tar -xzvf "${temp_dir}/kubero-cli.tar.gz" -C "$temp_dir" || { echo "Failed to unpack the binary."; rm -rf "$temp_dir"; exit 1; }

[[ ! -f "${temp_dir}/kubero" ]] && { echo "Failed to unpack the binary."; rm -rf "$temp_dir"; exit 1; }

echo "Installing kubero in /usr/local/bin ..."
sudo mv "${temp_dir}/kubero" "/usr/local/bin/kubero" || { echo "Failed to install kubero."; rm -rf "$temp_dir"; exit 1; }

rm -rf "$temp_dir"
echo "Kubero has been successfully installed."
echo "Run 'kubero install' to create a kubernetes cluster."