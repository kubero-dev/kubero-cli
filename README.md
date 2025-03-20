# kubero-cli

![Version](https://img.shields.io/github/v/release/kubero-dev/kubero-cli)
![Build Status](https://img.shields.io/github/actions/workflow/status/kubero-dev/kubero-cli/build.yml?branch=main)
![License](https://img.shields.io/github/license/kubero-dev/kubero-cli)

A powerful and user-friendly Command Line Interface (CLI) for [Kubero](https://github.com/kubero-dev/kubero), the simplest Platform as a Service (PaaS) for Kubernetes.

> **Note:** Please report any issues in the [main repository](https://github.com/kubero-dev/kubero).

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
   - [Supported Platforms](#supported-platforms)
   - [1. Shortcut Installation](#1-shortcut-installation)
   - [2. Homebrew Installation](#2-homebrew-installation)
   - [3. Build from Source](#3-build-from-source)
- [Supported Providers](#supported-providers)
- [Usage](#usage)
   - [Command Overview](#command-overview)
- [Provider Credentials](#provider-credentials)
   - [Scaleway](#scaleway)
   - [Linode](#linode)
   - [DigitalOcean](#digitalocean)
   - [Google GKE](#google-gke)
- [Development Guide](#development-guide)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)

---

## Features

- **Easy Cluster Deployment:** Quickly create Kubernetes clusters on supported providers.
- **App Management:** Simplify application deployment and management.
- **Pipeline Integration:** Seamlessly integrate CI/CD pipelines.
- **User-Friendly Commands:** Intuitive CLI commands for efficient workflows.
- **Dashboard Access:** Easy access to the Kubero dashboard for monitoring.

---

## Installation

### Supported Platforms

- **macOS**
- **Linux**

### 1. Shortcut Installation

Install Kubero CLI with a single command:

```shell
curl -fsSL get.kubero.dev | bash
```

### 2. Homebrew Installation

If Homebrew is not installed, install it first:

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Install Kubero CLI via Homebrew:

```shell
brew tap kubero-dev/kubero
brew install kubero-cli
```

### 3. Build from Source

For advanced use cases, build and package the binary manually.

#### Requirements

- [Git](https://git-scm.com/downloads)
- [Go](https://go.dev/doc/install)
- [UPX](https://github.com/upx/upx/releases/)

#### Steps

1. **Clone the Repository:**

   ```shell
   git clone https://github.com/kubero-dev/kubero-cli.git
   ```

2. **Navigate to the Project Directory:**

   ```shell
   cd kubero-cli
   ```

3. **Create a Version Tag (Optional):**

   ```shell
   git tag -a v1.0 -m "Version 1.0"
   ```

4. **Build and Package the Binary:**

   ```shell
   cd cmd
   go build -ldflags "-s -w -X main.version=$(git describe --tags --abbrev=0) -X main.commit=$(git rev-parse --short HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o kubero-cli
   upx kubero-cli
   mv kubero-cli ../kubero
   cd ..
   ```

5. **Move the Binary to Your PATH:**

   ```shell
   sudo mv kubero /usr/local/bin/
   ```

6. **Reload Shell Configuration:**

   ```shell
   source "$HOME/.$(basename ${SHELL})rc"
   ```

7. **Verify Installation:**

   ```shell
   kubero version
   ```

8. **Check for Updates:** (Optional)

   ```shell
   kubero version check
   ```

---

## Supported Providers

Kubero CLI currently supports the following cloud providers:

- **Scaleway**
- **Linode**
- **DigitalOcean**
- **Google GKE**
- **Kind** (local clusters)

### Coming Soon

- **Vultr**
- **Oracle Cloud OCI/OKE**
- **Exoscale**
- **Swissflow**

---

## Usage

### Command Overview

```plaintext
kubero
├── install                # Create a Kubernetes cluster and install Kubero with all required components
├── list                   # List all running clusters
├── login                  # Log in to Kubero and save credentials
├── logout                 # Log out from Kubero and remove saved credentials
├── create                 # Create new app and pipeline configurations
│   ├── app
│   └── pipeline
├── up                     # Deploy apps and pipelines
│   ├── app
│   └── pipeline
├── down                   # Remove apps and pipelines
│   ├── app
│   └── pipeline
├── fetch                  # Sync configurations to local files
│   ├── app
│   └── pipeline
├── dashboard              # Open the Kubero dashboard
├── tunnel                 # Open a tunnel to a NAT-ed cluster
├── instance               # Manage Kubero instances
│   ├── create             # Create an instance configuration
│   ├── delete             # Delete an instance configuration
│   └── select             # Select an active instance
├── config                 # View available configurations
│   ├── addons             # List addons
│   ├── buildpacks         # List buildpacks
│   └── podsizes           # List pod size configurations
└── help                   # Display help for commands
```

---

## Provider Credentials

Set the appropriate environment variables for your cloud provider before using Kubero CLI.

### Scaleway

```shell
export SCALEWAY_ACCESS_TOKEN=your_access_token
export SCALEWAY_PROJECT_ID=your_project_id
export SCALEWAY_ORGANIZATION_ID=your_organization_id
```

### Linode

```shell
export LINODE_ACCESS_TOKEN=your_access_token
```

### DigitalOcean

```shell
export DIGITALOCEAN_ACCESS_TOKEN=your_access_token
```

### Google GKE

```shell
export GOOGLE_API_KEY=your_api_key
```

---

## Development Guide

### Enable Development Mode

To enable development mode for testing and debugging, create a `VERSION` file:

```shell
echo "dev" > cmd/kuberoCli/VERSION
```

---

## Contributing

We welcome contributions from the community! Please check out our [Contributing Guidelines](https://github.com/kubero-dev/kubero/blob/main/CONTRIBUTING.md) for more information.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## Acknowledgments

- **[Kubero](https://github.com/kubero-dev/kubero):** The simplest PaaS for Kubernetes.
- **[Go](https://golang.org/):** The programming language used for development.
- **Community Contributors:** Thank you to all who have contributed to this project.

---

Thank you for using **kubero-cli**! If you have suggestions or encounter issues, please open an issue in the [main repository](https://github.com/kubero-dev/kubero).

---
