# Instant OpenHIE Core Package - Kubernetes

The Instant OpenHIE Core Package is the base of the Instant OpenHIE architecture.

This package consists of two services:

- Interoperability Layer - [OpenHIM](http://openhim.org/)
- FHIR Server - [HAPI FHIR](https://hapifhir.io/)

## Getting Started

Useful Links:

- [Kubectl Cheat Sheet](https://www.digitalocean.com/community/cheatsheets/getting-started-with-kubernetes-a-kubectl-cheat-sheet)

Before we proceed with creating our `Core Package` services, we need to ensure we are on the correct directory containing our bash setup scripts.

Once you are in the correct working directory (`core/kubernetes`) we can proceed to create our core Instant OpenHIE deployment with the following command:

```bash
./main/k8s.sh init
```

The OpenHIM console url will be displayed in the terminal output when the script completes. The Url may take a few minutes to become active as the pod may not be fully initialized yet.

This bash script will apply the kubernetes `kustomization.yaml` file which controls the `Core Package` components (ie: OpenHIM and HAPI-FHIR).

> On first run the setup may take up to 10 minutes as the Docker images for each component will need to be pulled. This won't happen on future runs.

### View running Kubernetes resources

Execute the below commands to see the running Kubernetes resources and the state that they are in.

To display all resource: (Some new resources are not listed here)

```sh
kubectl get all
```

To tear down this deployment use the opposing command:

```bash
./main/k8s.sh down
```

To start up the services after a tear down, use the following command:

```bash
./main/k8s.sh down
```

To completely remove all project components use the following option:

```bash
./main/k8s.sh destroy
```

---

## Initial OpenHIM Config

We have included a useful set of scripts to initialize the OpenHIM and set it up to communicate with the HAPI-FHIR server. This will change the default user's password of the OpenHIM to `instant101`, and create a channel configured to route traffic to the HAPI-FHIR instance. From the `kubernetes` directory, use the following command to implement:

```bash
./importer/k8s.sh up
```

> These scripts can be duplicated and modified to implement custom imports

To clean up the remaining job and pods from a successful setup, run the following:

```bash
./importer/k8s.sh clean
```

Once the config is done HAPI FHIR will be accessible on: <https://OPENHIM_CORE_HOSTNAME/fhir/>

You may test that the OpenHIM is routing requests to HAPI FHIR by running:

```bash
./test.sh <OPENHIM_CORE_HOSTNAME>
```

---

## Development mode for exposed services

To run in development mode, where the OpenHIM mongo database, HAPI FHIR server, and the MySQL database can be accessed directly through their urls, run the following command

```bash
./dev/k8s.dev.sh
```

---
