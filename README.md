# kubero-cli

![Version](https://img.shields.io/github/v/release/kubero-dev/kubero-cli)
![Build Status](https://img.shields.io/github/actions/workflow/status/kubero-dev/kubero-cli/build.yml?branch=main)
![License](https://img.shields.io/github/license/kubero-dev/kubero-cli)

A powerful and user-friendly Command Line Interface (CLI) for [Kubero](https://github.com/kubero-dev/kubero), the simplest Platform as a Service (PaaS) for Kubernetes.

> **Note:** Please report any issues in the [main repository](https://github.com/kubero-dev/kubero).

---

## Installation

### Using Makefile

With the updated workflow, installation is streamlined using the `Makefile`. Follow these steps:

1. **Clone the Repository:**

   ```shell
   git clone https://github.com/faelmori/kubero-cli.git
   ```

2. **Navigate to the Project Directory:**

   ```shell
   cd kubero-cli
   ```

3. **Run the Build Process:**

   ```shell
   make build
   ```

4. **Install the Binary:**

   ```shell
   make install
   ```

5. **Verify Installation:**

   ```shell
   kubero --version
   ```

### Using Install Script

You can also use the provided install script for a quick setup:

```shell
curl -fsSL get.kubero.dev | bash
```

---

## Usage

### As a CLI

The `kubero-cli` can be used to manage Kubernetes clusters and applications with various commands. Here are some examples:

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

### As a Go Module

To use `kubero-cli` as a module in your Go project, you can import the packages as needed. For example:

```go
import (
    "github.com/faelmori/kubero-cli/internal/api"
    "github.com/faelmori/kubero-cli/internal/config"
    "github.com/faelmori/kubero-cli/internal/db"
    "github.com/faelmori/kubero-cli/internal/utils"
)
```

You can then create instances of the clients and use them in your code:

```go
func main() {
    // Initialize API client
    client := api.NewClient("https://api.kubero.dev", "your_token")

    // Load configuration
    configLoader := &config.ViperConfig{}
    err := configLoader.LoadConfigs("/path/to/config", "config_name")
    if err != nil {
        log.Fatal(err)
    }

    // Initialize database
    dbClient := &db.GormDB{}
    db, err := dbClient.InitDB()
    if err != nil {
        log.Fatal(err)
    }
    dbClient.AutoMigrateDB(db, &db.Instance{})

    // Use utilities
    prompt := &utils.ConsolePrompt{}
    answer := prompt.PromptLine("Are you sure?", "[y/n]", "n")
    fmt.Println("Answer:", answer)
}
```

---

## Development

### Development Guide

To enable development mode for testing and debugging, create a `VERSION` file:

```shell
echo "dev" > cmd/kuberoCli/VERSION
```

### Contributing

We welcome contributions from the community! Please check out our [Contributing Guidelines](https://github.com/kubero-dev/kubero/blob/main/CONTRIBUTING.md) for more information.

### License

This project is licensed under the [MIT License](LICENSE).

### Acknowledgments

- **[Kubero](https://github.com/kubero-dev/kubero):** The simplest PaaS for Kubernetes.
- **[Go](https://golang.org/):** The programming language used for development.
- **Community Contributors:** Thank you to all who have contributed to this project.

---

Thank you for using **kubero-cli**! If you have suggestions or encounter issues, please open an issue in the [main repository](https://github.com/kubero-dev/kubero).
