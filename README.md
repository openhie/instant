# Instant OpenHIE

> Note: This repo is not for production. Rather, it contains strawpersons to facilitate discussion and demonstrations of progress on the Instant OpenHIE project. As such, it implies no endorsement or support from any institution especially and including the OpenHIE Community of Practice. This is for open discussion in the community. Please join the OpenHIE Dev-Ops Sub-community and give us your thoughts!

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