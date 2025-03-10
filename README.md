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
   - [1. Installation with Makefile](#1-installation-with-makefile)
   - [2. Homebrew Installation](#2-homebrew-installation)
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
- **Automated RSA Certificate Generation:** Securely generate encrypted RSA certificates with random passwords stored in a keyring.
- **Separate Test Flow:** Handle test logic separately from the core logic for better modularity.

<br/>

---

## Installation

### Supported Platforms

- **macOS**
- **Linux**

### 1. Installation with Makefile

With the updated workflow, installation is streamlined using the `Makefile`. Follow these steps:

1. **Clone the Repository:**

   ```shell
   git clone https://github.com/kubero-dev/kubero-cli.git
   ```

[//]: # (   ![Clone the Repository]&#40;#&#41;)

2. **Navigate to the Project Directory:**

   ```shell
   cd kubero-cli
   ```

3. **Run the Build Process:**

   ```shell
   make build
   ```

[//]: # (   ![Make Build Process]&#40;#&#41;)

4. **Install the Binary:**

   ```shell
   make install
   ```

[//]: # (   ![Make Install Process]&#40;#&#41;)

5. **Verify Installation:**

   ```shell
   kubero --version
   ```

[//]: # (   ![Version Check]&#40;#&#41;)

### 2. Homebrew Installation


&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;[Coming soon!!](https://github.com/kubero-dev/kubero-cli)
 

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
