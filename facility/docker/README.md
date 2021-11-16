# GOFR (Facility Reconciliation)

> Note: This is a docker-compose file to support the inclusion of a facility registry interface to Instant OpenHIE. This initial version is for testing.

The docker-compose runs the facility-recon app, Redis, MongoDB, and the Hearth FHIR server.

Clone the repo and start the docker apps.

To run:
```
docker-compose up
```
Or:
```
bash compose.sh init
bash compose.sh up
bash compose.sh down
bash compose.sh destroy
```

Visit: http://localhost:4000

The default admin user is `root@gofr.org` and pass is `gofr`. For production, a different admin user should immediately be created, the default deleted, and ordinary users added as well.

## Load all the Required Fhir Data and Meta data
Down load the [Hapi_FHIR CLI](https://github.com/hapifhir/hapi-fhir/releases/download/v5.5.2/hapi-fhir-5.5.2-cli.tar.bz2) and Extract it into a given directory .

Then from that Diretory run 

    ./hapi-fhir-cli upload-definitions -v r4 -t http://localhost:8092/fhir/DEFAULT

see more https://hapifhir.io/hapi-fhir/docs/tools/hapi_fhir_cli.html
# Todo

- [x] Initial version with compose bash script.
- [ ] Migrate FHIR backend to HAPI when complete in home repo.
- [ ] Add OpenHIM user, roles, channel, route.