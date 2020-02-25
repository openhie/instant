# InstantHIE Core Package - Kubernetes

The InstantHIE Core Package is the base of the InstantHIE architecture.

This package consists of two services:

* Interoperability Layer - [OpenHIM](http://openhim.org/)
* FHIR Server - [HAPI FHIR](https://hapifhir.io/)

## Minikube (local)

For the Kubernetes deployment to work as expected, we need to ensure we have `minikube` installed on our local machine running the deployment. Follow these steps to [install minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

### Getting Started

Before we proceed with creating our `Core Package` services, we need to ensure we are on the correct directory containing our bash setup scripts.

Once you are in the correct working directory (`core/kubernetes/main`) we can proceed to create our core instant ohie deployment with the following command:

```bash
./k8s.sh up
```

This bash script will enable `ingress` on minikube then proceed to apply the kubernetes `kustomization.yaml` file which controls the `Core Package` components (ie: OpenHIM and HAPI-FHIR). This script will also implement the HOST mapping which is needed to access the OpenHIM Core and Console locally (on linux).

> On first run the setup may take up to 10 minutes as the Docker images for each component will need to be pulled. This won't happen on future runs.

To tear down this deployment use the opposing command:

```bash
./k8s.sh down
```

To completely remove all project components use the following option:

```bash
./k8s.sh destroy
```

### Initial OpenHIM Config

We have included a useful script to initialise the OpenHIM and set it up to communicate with the HAPI-FHIR server. This will change the default user's password of the OpenHIM to `instant101`, and create a channel configured to route traffic to the HAPI-FHIR instance. Use the following command to implement:

```bash
kubectl apply -k ./importer/
```

> This script can be duplicated and modified to implement custom imports
