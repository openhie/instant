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

Once you are in the correct working directory (`core/kubernetes/`) we can proceed to execute our bash script by running the below command which will create all the pods, volumes and services.

```bash
./setup.sh
```

