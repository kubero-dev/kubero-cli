# kubero-cli
A CLI for kubero. The simplest PaaS for Kubernetes.
The main repository is [here](https://github.com/kubero-dev/kubero).

## Installation

Download the latest release [here](https://github.com/kubero-dev/kubero-cli/releases/latest) and extract the binary.

## Supported installer
- Scaleway
- Linode
- Digital Ocean
- Google GKE
- Kind (localhost)
- Oracle Cloud OCI/OKE
- Exoscale (soon)

## Usage
Command map
```
    kubero
    ├── apps
    │   ├── create
    │   ├── fetch
    │   ├── list
    │   └── delete
    ├── config
    │   ├── addons
    │   ├── buildpacks
    │   └── podsizes
    ├── help
    ├── init
    ├── install
    └── pipelines
        ├── create
        ├── fetch
        ├── list
        └── delete
```


## Usage
Command map
```
    kubero
    ├── install                 // create kubernetes cluster and install kubero with all required components
    ├── login                  // login to kubero, write credentials file
*    ├── logout                 // logout from kubero
    ├── create                 // create a new pipeline and app config
    │   ├── app
    │   └── pipeline
    ├── list                   // list all running pipelines and apps
    ├── up                     // deploy app and pipeline
    │   ├── app
    │   └── pipeline
    ├── down                   // delete app and pipeline
    │   ├── app
    │   └── pipeline
    ├── fetch                  // sync app and pipeline to local config
    │   ├── app
    │   └── pipeline
*    ├── status                 // show current Auth status
*    ├── dashboard              // Open the kubero dashboard
    ├── config                  // print configurations
    │   ├── addons
    │   ├── buildpacks
    │   └── podsizes
    └── help                    // Help about any command   
```


## Environment variables for credentials
### Scaleway
```
export SCALEWAY_ACCESS_TOKEN=xxx
export SCALEWAY_PROJECTID=xxx
export SCALEWAY_ORGANIZATIONID=xxx
```

### Linode
```
export LINODE_ACCESS_TOKEN=xxx
```

### Digital Ocean
```
export DIGITALOCEAN_ACCESS_TOKEN=xxx
```

### Google GKE
```
export GOOGLE_API_KEY=xxx
```

### Development
Create a dev VERSION File
```
echo "dev" > cmd/VERSION
```