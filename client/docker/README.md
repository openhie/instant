# OpenCR Client Registry

> Note: This is a docker-compose file to support the inclusion of a client registry interface to Instant OpenHIE. This initial version is for testing.

Clone the repo and start the docker apps.

> There is a long lag waiting for HAPI FHIR to start so the bash script uses a sleep to wait for it. If using the docker-compose script directly, OpenCR will fail before HAPI is done.

To run:

```sh
bash compose.sh init
bash compose.sh up
bash compose.sh down
bash compose.sh destroy
```

* Visit the UI at: [https://localhost:3003/crux](https://localhost:3003/crux)
  * **Default username**: root@intrahealth.org
  * **Default password**: intrahealth
