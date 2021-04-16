# Health Management Information System

A package that installs and starts up an instance of DHIS2

## Getting Started

### Docker

To start up the service

```sh
yarn
yarn docker:build
yarn docker:instant init -t docker hmis
```

To tear down the deployment use the opposing command:

```bash
yarn docker:instant down -t docker hmis
```

To start up the service after a tear down, use the following command:

```bash
yarn docker:instant up -t docker hmis
```

To completely remove all package components use the following option:

```bash
yarn docker:instant destroy -t docker hmis
```

### Kubernetes

To start up the service

```sh
yarn
yarn docker:build
yarn docker:instant init -t k8s hmis
```

To tear down the deployment use the opposing command:

```bash
yarn docker:instant down -t k8s hmis
```

To start up the service after a tear down, use the following command:

```bash
yarn docker:instant up -t k8s hmis
```

To completely remove all package components use the following option:

```bash
yarn docker:instant destroy -t k8s hmis
```

**NB** The commands above should be run in the instant folder

## Accessing the service

To access DHIS

### In Docker

* URL - `http://localhost:8081`

### In Kubernetes

Use the host ip or DNS, and the port `8081`. If the deployment has been done to a [minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/) cluster, get the external ip and port for the `dhis-web` service by running the following command

```sh
kubectl get services
```

The credentials for logging in are

* USERNAME - admin
* PASSWORD - district
