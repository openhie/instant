# InstantHIE Core Component - Kubernetes

The InstantHIE Core Component is the base of the InstantHIE architecture.

This component consists of two services:

* Interoperability Layer - [OpenHIM](http://openhim.org/)
* FHIR Server - [HAPI FHIR](https://hapifhir.io/)

## Getting Started

Before we proceed with creating our Core Component services, we need to ensure we are on the correct directory containing our **main.yml** script.

Once you are in the correct working directory (`core/kubernetes/`) we can proceed to execute our `kubectl` script by running the below command which will create all the services and print our their logs in the terminal.

```bash
kubectl apply -f main.yml
```
