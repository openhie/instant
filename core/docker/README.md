
# InstantHIE Core Component - docker-compose

The InstantHIE Core Component is the base of the InstantHIE architecture.

This component consists of two services:

* Interoperability Layer - [OpenHIM](http://openhim.org/)
* FHIR Server - [HAPI FHIR](https://hapifhir.io/)

## Getting Started

Before we proceed with creating our Core Component services, we need to ensure we are on the correct directory containing our `docker-compose` script.

Once you are in the correct working directory (`core/docker/`) we can proceed to execute our `docker-compose` script by running the below command which will create all the services and print our their logs in the terminal.

```bash
docker-compose up
```

### Useful compose flags

Some additional flags can be passed to the `docker-compose` command making it a bit easier to work with.

* `-d`: Run the services in a detached mode. This means that when you close or exit your terminal, the services will still be running in the background.
* `--force-recreate`: This will force the container/image to be re-created if a newer version is found. This is useful when a new images has been released but not yet pulled onto the host machine.

```bash
docker-compose up -d --force-recreate
```

## Accessing the services

### OpenHIM

Console: http://localhost:9000<br />
Username: root@openhim.org<br />
Password: instant101

### HAPI FHIR

This service should not be publicly accessible and only accessed via the Interoperability Layer

## Testing the Core Component

As part of the Core Component setup we also do some initial importation of config for connecting the services together. 

* OpenHIM: Import a public channel configuration that routes requests to the HAPI FHIR services
* HAPI FHIR: *Not config import yet*

For testing this Core Component we will be making use of `curl` for sending our request, but any client could be used to achieve the same result.

Execute the below `curl` request to successfully route a request through the OpenHIM to query the HAPI FHIR server.

```bash
curl http://localhost:5001/hapi-fhir-jpaserver/fhir/Patient
```
