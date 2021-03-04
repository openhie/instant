# OpenCR Client Registry

> Note: This is a docker-compose file to support the inclusion of a client registry interface to Instant OpenHIE. This initial version is for testing.

Clone the repo and start the docker apps. The `Client` Package relies on the `Core` package.
Therefore this compose script can't be run alone unless the Core package is already running.

To run:

```sh
bash compose.sh init
bash compose.sh up
bash compose.sh down
bash compose.sh destroy
```

* Visit the UI at: <https://localhost:3004/crux>
  * **Default username**: root@intrahealth.org
  * **Default password**: intrahealth
