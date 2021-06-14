#!/bin/bash

sleep 10

composeFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "init" ]; then
    docker create --name hmis-helper -v hmis-app-data:/DHIS2_home busybox
    docker cp "$composeFilePath"/importer/volume/dhis.conf hmis-helper:/DHIS2_home/dhis.conf
    docker rm hmis-helper

    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml up -d
elif [ "$1" == "up" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml up -d
elif [ "$1" == "down" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml stop
elif [ "$1" == "destroy" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml down

    docker volume rm hmis-app-data instant_hmis-db-data
else
    echo "Valid options are: init, up, down, or destroy"
fi
