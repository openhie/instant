# Mounts pipeline directory into the docker container so we can dev locally and automatically sync with the running Logstash container
docker-compose -p instant -f docker-compose.yml -f docker-compose.dev.yml --project-directory . up -d logstash
