# **kubero-cli**

A Command Line Interface (CLI) for Kubero, the simplest PaaS for Kubernetes.

The main repository is available [here](https://github.com/kubero-dev/kubero). Please use the main repository to report issues.

---

## **Installation**
### Supported Platforms: macOS, Linux

### **1. Shortcut (Script)**
Install using a single command:
```shell
curl -fsSL get.kubero.dev | bash
```

### **2. Homebrew (Package Manager)**
Install via Homebrew:
```shell
brew tap kubero-dev/kubero
brew install kubero-cli
```

### **3. Build from Source**
Manually build and package the binary.

#### **Requirements**
- [Git](https://git-scm.com/downloads)
- [GoLang](https://go.dev/doc/install)
- [UPX](https://github.com/upx/upx/releases/)

#### **Steps**
1. Clone the project repository (you may fork it first for proper tagging):
    ```shell
    git clone git@github.com:kubero-dev/kubero-cli.git
    ```
2. Navigate to the project directory:
    ```shell
    cd ./kubero-cli
    ```
3. Create a version tag:
    ```shell
    git tag -a v1.0 -m "Version 1.0"
    ```
4. Build and package the binary:
    ```shell
    cd ./cmd && \
    go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath && \
    upx ./cmd && \
    mv ./cmd ../kubero && \
    cd ../
    ```
5. Move the binary to a directory in your PATH (e.g., `/usr/bin` or `/usr/local/bin`):
    ```shell
    sudo mv kubero /usr/local/bin
    ```
6. Reload your shell profile:
    ```shell
    . "$HOME/.$(basename ${SHELL})rc"
    ```
7. Verify the installation by running:
    ```shell
    kubero --version
    ```

---

## **Supported Installers**
- **Scaleway**
- **Linode**
- **Digital Ocean**
- **Google GKE**
- **Kind (localhost)**

Coming soon:
- **Vultr**
- **Oracle Cloud OCI/OKE**
- **Exoscale**
- **Swissflow**

---

## **Usage**
### Command Map
```jsonc
kubero
├── install                // Create a Kubernetes cluster and install Kubero with all required components
├── list                   // List all running clusters
├── login                  // Log in to Kubero and save the instance to a credentials file
├── logout                 // Log out from Kubero and remove the instance from the credentials file
├── create                 // Create a new pipeline and app configuration
│   ├── app
│   └── pipeline
├── up                     // Deploy an app and pipeline
│   ├── app
│   └── pipeline
├── down                   // Delete an app and pipeline
│   ├── app
│   └── pipeline
├── fetch                  // Sync app and pipeline configurations to local files
│   ├── app
│   └── pipeline
├── dashboard              // Open the Kubero dashboard
├── tunnel                 // Open a tunnel to a NAT-ed cluster
├── instance               // Manage Kubero instances
│   ├── create             // Create a configuration for a Kubero instance
│   ├── delete             // Delete a configuration for a Kubero instance
│   └── select             // Select a Kubero instance
├── config                 // Print configurations
│   ├── addons             // List available addons
│   ├── buildpacks         // List available buildpacks
│   └── podsizes           // List pod size configurations
└── help                   // Display help for any command
```

---

## **Environment Variables for Credentials**

### **Scaleway**
Set the following environment variables:
```shell
export SCALEWAY_ACCESS_TOKEN=your_access_token
export SCALEWAY_PROJECTID=your_project_id
export SCALEWAY_ORGANIZATIONID=your_organization_id
```

### **Linode**
Set the following environment variable:
```shell
export LINODE_ACCESS_TOKEN=your_access_token
```

### **Digital Ocean**
Set the following environment variable:
```shell
export DIGITALOCEAN_ACCESS_TOKEN=your_access_token
```

### **Google GKE**
Set the following environment variable:
```shell
export GOOGLE_API_KEY=your_api_key
```

---

## **Development**
### Create a Development Version File
To enable development mode, create a `VERSION` file:
```shell
echo "dev" > cmd/kuberoCli/VERSION
```

This file ensures that the application operates in development mode during testing.

