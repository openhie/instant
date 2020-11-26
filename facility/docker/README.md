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

Visit: http://localhost:3000

The default admin user is `root@gofr.org` and pass is `gofr`. For production, a different admin user should immediately be created, the default deleted, and ordinary users added as well.

# Todo

- [x] Initial version with compose bash script.
- [ ] Migrate FHIR backend to HAPI when complete in home repo.
- [ ] Add OpenHIM user, roles, channel, route.