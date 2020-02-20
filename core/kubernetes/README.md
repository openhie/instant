# InstantHIE Core Package - Kubernetes

The InstantHIE Core Package is the base of the InstantHIE architecture.

This package consists of two services:

* Interoperability Layer - [OpenHIM](http://openhim.org/)
* FHIR Server - [HAPI FHIR](https://hapifhir.io/)

## Minikube (local)

### Getting Started

Before we proceed with creating our Core Package services, we need to ensure we are on the correct directory containing our bash setup scripts.

Once you are in the correct working directory (`core/kubernetes/`) we can proceed to execute our bash script by running the below command which will create all the pods, volumes and services.

```bash
./setup.sh
```

#### Configure OpenHIM Console

A manual configuration is required to enter the Public URL of OpenHIM Core into the configuration of OpenHIM Console.

Whilst on your local setup the services won't be assigned external IPs by default. To do this on **minikube** we have to request URLs by which we can access the desired services.

The first service we will expose is the OpenHIM Core service:

```sh
minikube service -n default openhim-core-service
```

Copy the `core-8080` url.

Edit the console-config configMap:

```sh
kubectl edit configmap console-config
```

Go to the line "host" and replace the value "{enter-openhim-core-service-host-here}" with the **host** and the `port` value of the openhim-core you copied. Example:

```json
"host": "192.168.99.101",
"port": 32289
```

> Remember to remove the protocol and path from the url

Save the file

```sh
press escape
:wq
```

Restart the deployments

```sh
./k8s.sh down
./k8s.sh up
```

In a new terminal, expose the OpenHIM Console service with the following:

```sh
minikube service -n default openhim-core-service
```

Navigate to that url and login to the OpenHIM
