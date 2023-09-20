# kubero-cli
A CLI for kubero. The simplest PaaS for Kubernetes.
The main repository is [here](https://github.com/kubero-dev/kubero). Please use the main repository to open issues. 

## Installation

Download the latest release [here](https://github.com/kubero-dev/kubero-cli/releases/latest) and extract the binary.

## Supported installer
- Scaleway
- Linode
- Digital Ocean
- Google GKE
- Kind (localhost)
- Vultr (soon)
- Oracle Cloud OCI/OKE (soon)
- Exoscale (soon)
- Swissflow (soon)


## Usage
Command map
```
    kubero
    ├── install                 // create kubernetes cluster and install kubero with all required components
*    ├── login                  // login to kubero, safe instance to credentials file
*    ├── logout                 // logout from kubero, remove instance from credentials file
*    ├── instance               // print current kubero instance
*    │   ├── create             // create a configuration to a kubero instance
*    │   ├── delete             // delete a configuration to a kubero instance
*    │   ├── select             // select a kubero instance
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
    ├── config                  // print configurations
    │   ├── addons
    │   ├── buildpacks
    │   └── podsizes
    ├── dashboard               // Open the kubero dashboard
    ├── tunnel                  // Open a tunnel to a natted cluster
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
