> [!WARNING]
> Instant OpenHIE v1 has been discontinued in favour of Instant OpenHIE v2. Please see the [Instant OpenHIE v2 docs](https://jembi.gitbook.io/instant-v2/) for more details.

# Instant OpenHIE

Welcome to Instant OpenHIE. Open Health Information Exchange (OpenHIE) is a global, mission-driven community of practice dedicated to improving the health of the underserved through open and collaborative development and support of country driven, large scale health information sharing architectures. 

The central technical output of OpenHIE is the consensus-driven OpenHIE Architecture, an open framework to develop health information exchanges (HIEs) to improve patient care, public health, and the management of health resources. The architecture is based upon real-world implementations and needs.

Instant OpenHIE is meant to help implementers to test and try OpenHIE Architecture components and workflows. The goals of Instant OpenHIE are to support:
* **Demonstrations**: A simple way to launch containers with OpenHIE components and demonstrate workflows.
* **Simplify life for developers**: A way for implementers to rapidly develop their tooling around the OpenHIE Architecture and the versions of it without having to unnecessary time configuring many components.
* **QA**: Run workflows in a CI/CD system like GitHub Actions.

Instant OpenHIE is not meant for production and is not a substitute for enterprise-level IT configuration. Implementers should consult the product owners' documentation for production set-up.

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

### Supported options

Instant supports a few options after the main command in the form: `instant <main_command> [options] [package-ids]` e.g. `yarn docker:instant init -t k8s`

* --target, -t : specify a target to deploy to, defaults to docker. Allowed values: docker, kubernetes, k8s
* --only, -o : if supplied only the specified packages will be acted upon in the order they are provided, ignoring dependencies

## Custom packages

To add a custom package to your instant instance use the following flag

`-c="/path/to/package"`

For example, if you had downloaded the [who-covid19-surveillance-package](https://github.com/jembi/who-covid19-surveillance-package) repository onto your machine you could reference it as follows:

```sh
yarn docker:instant init core covid19surveillance -c="../who-covid19-surveillance-package"
```

Urls are supported. The custom package will be downloaded and then mounted to the instant instance.

`-c="https://github.com/jembi/who-covid19-surveillance-package"`

> Only github repos and urls pointing to files with extensions "zip" and "tar.gz"

> Packages that are nested in folders can also be mounted by specifying the path or url of the root folder. Packages can be nested 5 levels down. For example one can mount a package within the following folder `test/test1/test2/test3/test4/package` by using the path to the `test` folder. This also allows us to mount multiple packages contained in a folder by specifying the path to that folder.

### Go CLI Binary (Cross platform)

The Instant OpenHIE project is available as a Docker image therefore we do not need the whole GitHub repository to run the containers.

For a minimum Instant OpenHIE set up, download [the Go Binary for your distribution](https://github.com/openhie/instant/releases).
The binaries are in the assets section of each release and you have a choice of linux, mac, and windows to choose from.

#### Binary Interactive Mode

To run the `go cli` binary, launch the project as follows:

* **Linux**. From terminal run: `./instant-linux`
* Mac. From terminal run: `./instant-macos`
  > Warning: Mac has an issue with the binary as it views the file as a security risk. See [this article](https://www.lifewire.com/fix-developer-cannot-be-verified-error-5183898) to bypass warning
* Windows. Double click: `instant.exe`

Then choose your options and deploy!

#### Binary Non-interactive Mode (Linux/Mac)

From a terminal, run the following command to add your custom package and initialise the system in docker.

```sh
./instant-{os} init -t docker core <your_package_ids> -c="../path/to/your/package"
```

To remove the instant project, run the following:

./instant-{os} destroy -t docker core covid19surveillance

> The custom package location is not needed for `up`, `down`, or `destroy` commands on an existing system.

To initialise kubernetes, run the following:

```sh
./instant-{os} init -t k8s core <your_package_ids> -c="../path/to/your/package"
```

Multiple custom packages can be chained together as follows:

```sh
./instant-{os} init test1 test2 test3 -c="../test1" -c="../test2" -c="../test3"
```

### Running tests on running packages

The [Cucumber](https://cucumber.io/) framework is used for testing the instantiated packages.

#### Docker-Compose

Run the following for the default tests:

```sh
yarn test:local <PACKAGE_IDs>
```

If you want to make changes to the tests, you can run your changes without rebuilding anything by using the *dev* version of the command:

> Remember to update the **volume file path** in the `package.json`

```sh
yarn test:local:dev <PACKAGE_IDs>
```

To run custom package tests in your local environment, no changes need to be made to your setup as the existing `instant` volume contains the custom package in the appropriate directory.
However, if you want to make changes to the custom package tests you will need to add a new volume reference to the npm test script.
For example, to experiment with the [WHO Covid19 Surveillance Package](https://github.com/jembi/who-covid19-surveillance-package) tests we will need to add the local file path of the package to the test command.
Line 21 of our package.json should look something like this (substituting in your specific file path):

```json
    "test:local:dev": "docker run --rm --name test-helper -v </absolute/path/to/instant>:/instant -v </absolute/path/to/who-covid19-surveillance-package>:/instant/who-covid19-surveillance-package --network instant_default openhie/package-test local",
```

This will allow us to make changes to the Covid19 Surveillance tests without having to rebuild the containers between runs.

#### Kubernetes

Update the `.env.remote` file with the instances' host urls and ports.
Then update the file path in the [`package.json`](./package.json) file on line 22 in the scripts section.
This file path needs to reference your `.env.remote` file to volume in your updates.
Finally, run the following command for the default tests:

```sh
yarn test:remote <PACKAGE_IDs>
```

If you want to make changes to the tests, you can run your changes without rebuilding anything by using the *dev* version of the command:

> Remember to update the **volume file path** in the `package.json`

```sh
yarn test:remote:dev <PACKAGE_IDs>
```

> The `PACKAGE_IDs` is a string of the package ids separated by space.

To run custom package tests in the remote environment, no changes need to be made to your setup as the existing `instant` volume contains the custom package in the appropriate directory.
However, if you want to make changes to the custom package tests you will need to add a new volume reference to the npm test script.
For example, to experiment with the [WHO Covid19 Surveillance Package](https://github.com/jembi/who-covid19-surveillance-package) tests we will need to add the local file path of the package to the test command.
Line 23 of our package.json should look something like this (substituting in your specific file path):

```json
    "test:remote:dev": "docker run --rm --name test-helper -v </absolute/path/to/instant>:/instant -v </absolute/path/to/who-covid19-surveillance-package>:/instant/who-covid19-surveillance-package --network host openhie/package-test remote",
```

This will allow us to make changes to the Covid19 Surveillance tests without having to rebuild the pods between runs.
