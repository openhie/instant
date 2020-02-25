# InstantHIE Core Package - Kubernetes

The InstantHIE Core Package is the base of the InstantHIE architecture.

This package consists of two services:

* Interoperability Layer - [OpenHIM](http://openhim.org/)
* FHIR Server - [HAPI FHIR](https://hapifhir.io/)

## Minikube (local)

For the Kubernetes deployment to work as expected, we need to ensure we have `minikube` installed on our local machine running the deployment. Follow these steps to [install minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

Once `minikube` has been installed we need to enable the `ingress` module for handling some of our networking configuration.

Execute the below command:

```sh
minikube addons enable ingress
```

### Getting Started

Before we proceed with creating our Core Package services, we need to ensure we are on the correct directory containing our bash setup scripts.

Once you are in the correct working directory (`core/kubernetes/`) we can proceed to create our core instant ohie deployment with the following command:

```bash
kubectl apply -k ./main/
```

This will deploy the OpenHIM and HAPI-FHIR along with all their services and dependencies.

To tear down this deployment use the opposing command:

```bash
kubectl delete -k ./main/
```

### Initial OpenHIM Config

We have included a useful script to initialise the OpenHIM and set it up to communicate with the HAPI-FHIR server. This will change the default user's password of the OpenHIM to `instant101`, and create a channel configured to route traffic to the HAPI-FHIR instance. Use the following command to implement:

```bash
kubectl apply -k ./importer/
```

> This script can be duplicated and modified to implement custom imports
