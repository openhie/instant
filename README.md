# Instant OpenHIE

> Note: This repo is not for production. Rather, it contains strawpersons to facilitate discussion and demonstrations of progress on the Instant OpenHIE project. As such, it implies no endorsement or support from any institution especially and including the OpenHIE Community of Practice. This is for open discussion in the community. Please join the OpenHIE Dev-Ops Sub-community and give us your thoughts!

## Getting Started

The services can be deployed using docker-compose or kubernetes.

### Docker-compose

Navigate to the main folder to execute the commands.

To set the Instant OpenHIE services run the following command:

```sh
./instant.sh init docker
```

To tear down the deployments use the opposing command:

```bash
./instant.sh down docker
```

To start up the services after a tear down, use the following command:

```bash
./instant.sh up docker
```

To completely remove all project components use the following option:

```bash
./instant.sh destroy docker
```

## Kubernetes

A kubernetes deployment can either be to AWS using [eksctl](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html) and [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) or locally using [minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/) and `kubectl`.

For a quick setup of these environments navigate to `/core/kubernetes/README.md`

Navigate to the main folder to execute the commands.

To set the Instant OpenHIE services run the following command:

```sh
./instant.sh init k8s
```

To tear down the deployments use the following command:

```bash
./instant.sh down k8s
```

To start up the services after a tear down, use the following command:

```bash
./instant.sh up k8s
```

To completely remove all project components use the following option:

```bash
./instant.sh destroy k8s
```
