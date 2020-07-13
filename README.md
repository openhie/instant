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

## Docker outside of Docker (DooD) - proof of concept

First build the image for the Orchestration container. This container will interact with the hosts docker socket and run the scripts to start and stop the packages:

```sh
docker build -t instant .
# eventually this won't be necessary as we can build and push an
# image to docker hub, this can be automatically downloaded when
# running the docker commands in the go app.
```

Next, run the go app. This app will run the orchestrator container with a particular command (up, down, destroy) which will run each package's script with that command. It will effectively start or stop each of the configured packages.

```sh
cd goinstant
go run goinstant.go
```

Select start to start up the infrastructure and stop to shut it all down. Note the containers all run on the host system even though they are kicked off in a sibling container. This is because we pass in the docker socket from the host via a volume. A similar strategy will work for k8s, we must just point kubectl to a new API server by passing in via volume `.kube/config`.

This method will still work (perhaps with slight adjustments) on Windows and OSx as both of those use a tiny linux vm to run docker containers on. This means that we can still pass the docker socket file via `-v /var/run/docker.sock:/var/run/docker.sock`. I have managed to verify that this works on Windows.

If we want to add more packages that aren't included in the base docker orchestration container then we can make the orchestration container able to detect extra (3rd party) packages in a particular folder and execute their bash scripts to start and stop them. The go app can manage downloading these extra packages and incudes them in the orchestration container via a volume `-v ~/.instant/extra-packages:/extra-packages`.

## Scratchpad

The commands that the go app runs in the background are similar to these:

```sh
docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock instant up
```

```sh
docker run -it -v /var/run/docker.sock:/var/run/docker.sock -v ~/.kube/config:/root/.kube/config -v ~/.minikube:/home/ryan/.minikube -v ~/.minikube:/root/.minikube -v /usr/local/bin/minikube:/usr/local/bin/minikube --network host instant up k8s
```
