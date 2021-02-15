# Instant OpenHIE

> Note: This repo is not for production. Rather, it contains strawpersons to facilitate discussion and demonstrations of progress on the Instant OpenHIE project. As such, it implies no endorsement or support from any institution especially and including the OpenHIE Community of Practice. This is for open discussion in the community. Please join the OpenHIE Dev-Ops Sub-community and give us your thoughts!

View the [user documentation](https://openhie.github.io/instant/) for more information.

## Getting Started

The services can be deployed using docker-compose or kubernetes.

### Docker-compose

Navigate to the main folder to execute the commands.

To set the Instant OpenHIE services run the following command:

```sh
yarn
yarn docker:build
yarn docker:instant init -t docker
```

To tear down the deployments use the opposing command:

```bash
yarn docker:instant down -t docker
```

To start up the services after a tear down, use the following command:

```bash
yarn docker:instant up -t docker
```

To completely remove all project components use the following option:

```bash
yarn docker:instant destroy -t docker
```

Each command also takes a list of package IDs to operate on. If this is left out then all packages are run by default.

E.g only run `core` package: `yarn docker:instant init -t docker core`

### Kubernetes

A kubernetes deployment can either be to AWS using [eksctl](https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html) and [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) or locally using [minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/) and `kubectl`.

For a quick setup of these environments navigate to [the Kubernetes development readme](kubernetes.md)

Navigate to the main folder to execute the commands.

To set the Instant OpenHIE services run the following command:

```sh
yarn
yarn docker:build
yarn docker:instant init -t k8s
```

To tear down the deployments, use the following command:

```bash
yarn docker:instant down -t k8s
```

To start up the services after a tear down, use the following command:

```bash
yarn docker:instant up -t k8s
```

To completely remove all project components, use the following option:

```bash
yarn docker:instant destroy -t k8s
```

Each command also takes a list of package IDs to operate on. If this is left out then all packages are run by default.

E.g only run `core` package: `yarn docker:instant init -t k8s core`

## Custom packages

To add a custom package to your instant instance use the following flag

`-c="/path/to/package"`

For example, if you had downloaded the [who-covid19-surveillance-package](https://github.com/jembi/who-covid19-surveillance-package) repository onto your machine you could reference it as follows:

```sh
yarn docker:instant init core covid19surveillance -c="../who-covid19-surveillance-package"
```

> We hope to support package urls soon
