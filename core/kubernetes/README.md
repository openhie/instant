# InstantHIE Core Package - Kubernetes

The InstantHIE Core Package is the base of the InstantHIE architecture.

This package consists of two services:

* Interoperability Layer - [OpenHIM](http://openhim.org/)
* FHIR Server - [HAPI FHIR](https://hapifhir.io/)

## Minikube (local)

For the Kubernetes deployment to work as expected, we need to ensure we have `minikube` installed on our local machine running the deployment. Follow these steps to [install minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

Once installed, we can start the minikube service be executing the below command:

```sh
minikube start --cpus 4 --memory 8192
```

This also updates the VM settings to make use of 4 CPU's and 8GB of RAM, instead of the default 2 CPU's and 4GB of RAM

### Getting Started

Before we proceed with creating our `Core Package` services, we need to ensure we are on the correct directory containing our bash setup scripts.

Once you are in the correct working directory (`core/kubernetes`) we can proceed to create our core instant ohie deployment with the following command:

```bash
./main/k8s.sh up
```

This bash script will enable `ingress` on minikube then proceed to apply the kubernetes `kustomization.yaml` file which controls the `Core Package` components (ie: OpenHIM and HAPI-FHIR). This script will also implement the HOST mapping which is needed to access the OpenHIM Core and Console locally (on linux).

> On first run the setup may take up to 10 minutes as the Docker images for each component will need to be pulled. This won't happen on future runs.

#### View running Kubernetes resources

Execute the below commands to see the running Kubernetes resources and the state that they are in.

To display all resource: (Some new resources are not listed here)

```sh
kubectl get all
```

To view all the ingress resources: (Services that are exposed publicly)

```sh
kubectl get ingress
```

To tear down this deployment use the opposing command:

```bash
./main/k8s.sh down
```

To completely remove all project components use the following option:

```bash
./main/k8s.sh destroy
```

The OpenHIM console will be accessible on http://openhim-console.instant/ and core will be accessible on:
* API: http://openhim-core.api.instant/
* HTTPS routing: https://openhim-core.ssl.instant/
* HTTP routing: http://openhim-core.non-ssl.instant/

### Initial OpenHIM Config

We have included a useful set of scripts to initialise the OpenHIM and set it up to communicate with the HAPI-FHIR server. This will change the default user's password of the OpenHIM to `instant101`, and create a channel configured to route traffic to the HAPI-FHIR instance. From the `kubernetes` directory, use the following command to implement:

```bash
./importer/k8s.sh up
```

> These scripts can be duplicated and modified to implement custom imports

To clean up the remaining job and pods from a successful setup run the following:

```bash
./importer/k8s.sh clean
```

Once the config is done HAPI FHIR will be accessible on: https://openhim-core.ssl.instant/hapi-fhir-jpaserver/fhir/

You may test that the OpenHIM is routing requests to HAPI FHIR by running:

```bash
./test.sh openhim-core.ssl.instant
```

### Development mode for exposed services

To run in development mode, where the OpenHIM mongo database, HAPI fhir server and the MySQL database can be accessed directly through their urls, run the following command

```bash
./dev/k8s.dev.sh
```
