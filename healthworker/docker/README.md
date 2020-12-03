# Docker

## Using `docker-compose` directly

For more information, see the iHRIS [technical documentation](https://ihris.github.io/iHRIS/admin/docker/).

Once you are in the correct working directory (`healthworker/docker/`) run the following commands to create all the services and print logs in the terminal.

```bash
docker-compose -f docker-compose.hapi.yml up -d
# the next two containers are config only
docker-compose -f docker-compose.hapi.config.yml up
docker-compose -f docker-compose.ihris.config.yml up
# demo data is optional
docker-compose -f docker-compose.ihris.data.yml up
# must launch ES and Kibana before ihris
docker-compose -f docker-compose.elastic.yml up -d
# includes redis which must run before ihris
docker-compose -f docker-compose.ihris.yml up
```

## Using `./compose.sh`

```bash
bash compose.sh init
bash compose.sh up
bash compose.sh down
bash compose.sh destroy
```