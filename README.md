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
    ├── login
    ├── install
    ├── create
    │   ├── app
    │   └── pipeline
    ├── list
    │   ├── app
    │   └── pipeline
    ├── up
    ├── down
    ├── sync
    │   ├── app
    │   └── pipeline
    ├── config
    │   ├── addons
    │   ├── buildpacks
    │   └── podsizes
    └── help
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