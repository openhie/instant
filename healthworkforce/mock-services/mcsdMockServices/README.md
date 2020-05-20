# mCSD Mock Service

This service is not meant to be implemented into production!

This service is just a mock of the actual implementation that still needs to be incorporated.

This one service creates a few endpoints which simulate the external lookup requests

### Endpoints in mock service

* `GET /gofr-location-mock/_history` - Returns a GOFR FHIR Location bundle
* `GET /gofr-organization-mock/_history` - Returns a GOFR FHIR Organization bundle
* `GET /ihris-practitioner-mock/_history` - Returns a iHRIS FHIR Practitioner bundle
* `GET /ihris-practitionerRole-mock/_history` - Returns a iHRIS FHIR PractitionerRole bundle
* `POST /fhir-mock` - Returns a FHIR transaction bundle

## Setup

Install all the service dependencies

```sh
yarn install
```

Start the service

```sh
yarn start
```

### Docker

Build the docker image

```sh
docker build -t jembi/mcsd-mock-service:latest .
```

Run the docker container

```sh
docker run -it -p 4000:4000 jembi/mcsd-mock-service:latest
```
