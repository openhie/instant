#!/bin/bash

composeFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "init" ]; then
    docker stack deploy  -c "$composeFilePath"/docker-compose-mongo.yml -c "$composeFilePath"/docker-compose-mongo.stack-0.yml instant

    # Set up the replica set
    "$composeFilePath"/initiateReplicaSet.sh

    docker stack deploy -c "$composeFilePath"/docker-compose.yml -c "$composeFilePath"/docker-compose.dev.yml -c "$composeFilePath"/importer/docker-compose.config.yml -c "$composeFilePath"/docker-compose.stack-1.yml -c "$composeFilePath"/importer/docker-compose.config.stack-0.yml instant
elif [ "$1" == "up" ]; then
    docker stack deploy -c "$composeFilePath"/docker-compose.mongo.yml -c "$composeFilePath"/docker-compose-mongo.stack-0.yml instant

    # Wait for mongo replica set to be set up
    sleep 20
    docker stack deploy -c "$composeFilePath"/docker-compose.yml -c "$composeFilePath"/docker-compose.dev.yml -c "$composeFilePath"/docker-compose.stack-1.yml instant
elif [ "$1" == "down" ]; then
    docker service scale instant_openhim-core=0 instant_openhim-console=0 instant_hapi-fhir=0 instant_hapi-mysql=0
elif [ "$1" == "destroy" ]; then
    docker service rm instant_core instant_console instant_fhir instant_mysql

    docker volume rm hapi-mysql hapi-mysql-config instant_openhim-mongo1 instant_openhim-mongo2 instant_openhim-mongo3
else
    echo "Valid options are: init, up, down, or destroy"
fi
