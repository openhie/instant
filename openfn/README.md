# OpenFn Microservice

A demo package using OpenFn Microservice, OpenHIM and HAPI FHIR using Instant OpenHIE.

The `microservice` project.yaml file is located here: [`openfn/docker/config/project.yaml`](./docker/config/project.yaml).

## Getting Started

### Docker

To start up the service

```sh
yarn
yarn docker:build
yarn docker:instant init -t docker openfn
```

To tear down the deployment use the opposing command:

```bash
yarn docker:instant down -t docker openfn
```

To start up the service after a tear down, use the following command:

```bash
yarn docker:instant up -t docker openfn
```

To completely remove all package components use the following option:

```bash
yarn docker:instant destroy -t docker openfn
```

### Kubernetes

To start up the service

```sh
yarn
yarn docker:build
yarn docker:instant init -t k8s openfn
```

To tear down the deployment use the opposing command:

```bash
yarn docker:instant down -t k8s openfn
```

To start up the service after a tear down, use the following command:

```bash
yarn docker:instant up -t k8s openfn
```

To completely remove all package components use the following option:

```bash
yarn docker:instant destroy -t k8s openfn
```

## How data gets to HAPI FHIR

Using the example payload [commcare_sample.json](./fixtures/commcare_sample.json)
we send that to the configured OpenFn Microservice.

Microservice is configurated to run a job based on the shape of the incoming 
payload _see [project.yaml](./docker/config/project.yaml)_.

The job `commcare-to-him` will match against this message and will be invoked
performing the following actions:

- creates a payload in the FHIR standard containing
  - a Encounter resource that contains (`contained` resource field) a Patient resource
- sends the payload to OpenHIM
- which in turn sends the payload to HAPI FHIR


### In Docker

* URL - `http://localhost:4001`

### In Kubernetes

Use the host ip or DNS, and the port `4001`. If the deployment has been done to a [minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/) cluster, get the external ip and port for the `openfn-service` service by running the following command

```sh
kubectl get services
```

## Notes

- The HAPI FHIR service runs on port `3447`
- The OpenHIM channel that we go through is on port `5001`.  
  The API is identical, with the exception of a required `Authorization` header.

