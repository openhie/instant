# mCSD endpoint flow

This guide will assist in getting the mCSD endpoint setup and tested with the mock service to return the relevant lookup results

## Create Endpoint record

```sh
curl -X POST -d "@endpoint.json" -H "Content-Type: application/json" http://localhost:3003/endpoints
```

## Start the mock service

Navigate to the `mcsdMockServices` directory and run the below commands to start the services

Install all the required dependencies:

```sh
yarn
```

Start the service:

```sh
yarn start
```

### Endpoints in mock service

* `GET /gofr-location-mock/_history` - Returns a GOFR FHIR Location bundle
* `GET /gofr-organization-mock/_history` - Returns a GOFR FHIR Organization bundle
* `GET /ihris-practitioner-mock/_history` - Returns a iHRIS FHIR Practitioner bundle
* `GET /ihris-practitionerRole-mock/_history` - Returns a iHRIS FHIR PractitionerRole bundle
* `POST /fhir-mock` - Returns a FHIR transaction bundle

## Testing the endpoint

To test the endpoint, execute the below curl command to trigger the endpoint for processing

```sh
curl -X POST -d "{}" -H "Content-Type: application/json" http://localhost:3003/mcsd
```

> NB! We are sending an empty JSON POST request to this endpoint. This is because the mapping mediator currently only accepts incoming POST requests, and its needs to be a valid JSON object
