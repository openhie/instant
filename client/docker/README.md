# OpenCR Client Registry

> Note: This is a docker-compose file to support the inclusion of a client registry interface to Instant OpenHIE. This initial version is for testing.

Clone the repo and start the docker apps.

> There is a long lag waiting for HAPI FHIR to start so the bash script uses a sleep to wait for it. If using the docker-compose script directly, OpenCR will fail before HAPI is done.

To run:
```
bash compose.sh up
# the bash support script also includes:
bash compose.sh stop
bash compose.sh down
bash compose.sh destroy
```

* Visit the UI at: [https://localhost:3000/crux](https://localhost:3000/crux)
    * **Default username**: root@intrahealth.org 
    * **Default password**: intrahealth

# Todo
- [x] Initial version with compose bash script.
- [ ] Migrate FHIR backend to HAPI when complete in home repo.
- [ ] Add OpenHIM user, roles, channel, route.