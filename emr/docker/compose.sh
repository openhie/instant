#!/bin/bash

composeFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "init" ]; then
    docker create --name openmrs-db-helper -v openmrs-dbdump:/dbdump busybox
    docker cp "$composeFilePath"/importer/volume/dbdump/. openmrs-db-helper:/dbdump/
    docker rm openmrs-db-helper

    docker create --name openmrs-server-helper -v openmrs-modules:/modules busybox
    docker cp "$composeFilePath"/importer/volume/modules/. openmrs-server-helper:/modules/
    docker rm openmrs-server-helper

    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml up -d

elif [ "$1" == "up" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml up -d

elif [ "$1" == "down" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml stop

elif [ "$1" == "destroy" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml down -v

else
    echo "Valid options are: init, up, down, or destroy"
fi

