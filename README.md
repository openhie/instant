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

> The ports of the services will change when the services are brought back up.

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

### Docker or Kubernetes without the Instant OpenHIE repo

The Instant OpenHIE project is available as a Docker image therefore we do not need the whole GitHub repository to run the containers.

For a minimum Instant OpenHIE set up, download [this deploy script from GitHub](https://raw.githubusercontent.com/openhie/instant/master/deploy.sh).
Once downloaded make sure it's executable: `sudo chmod +x deploy.sh`

Then, run the following command to add your custom package and initialise the system in docker.

```sh
./deploy init -t docker core <your_package_ids> -c="../path/to/your/package"
```

To remove the instant project, run the following:

./deploy destroy -t docker core covid19surveillance

> The custom package location is not needed for `up`, `down`, or `destroy` commands on an existing system.

To initialise kubernetes, run the following:

```sh
./deploy init -t k8s core <your_package_ids> -c="../path/to/your/package"
```

Multiple custom packages can be chained together as follows:

```sh
./deploy init test1 test2 test3 -c="../test1" -c="../test2" -c="../test3"
```

### Running tests on running packages

The [Cucumber](https://cucumber.io/) framework is used for testing the instantiated packages.

#### Docker-Compose

Run the following for the default tests:

```sh
yarn test:local <PACKAGE_IDs>
```

If you want to make changes to the tests, you can run your changes without rebuilding anything by using the *dev* version of the command:

> Remember to update the volume file path in the `package.json`

```sh
yarn test:local:dev <PACKAGE_IDs>
```

#### Kubernetes

Update the `.env.remote` file with the instances' host urls and ports.
Then update the file path in the [`package.json`](./package.json) file on line 22 in the scripts section.
This file path needs to reference your `.env.remote` file to volume in your updates.
Finally, run the following command for the default tests:

```sh
yarn test:remote <PACKAGE_IDs>
```

If you want to make changes to the tests, you can run your changes without rebuilding anyting by using the *dev* version of the command:

> Remember to update the volume file path in the `package.json`

```sh
yarn test:remote:dev <PACKAGE_IDs>
```

> The `PACKAGE_IDs` is a string of the package ids separated by space.
