#!/bin/bash

composeFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "init" ]; then

    docker create --name opencr-helper -v opencr-data:/config busybox
    docker cp "$composeFilePath"/importer/volume/. opencr-helper:/config/
    docker rm opencr-helper

    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml -f "$composeFilePath"/importer/docker-compose.config.yml up -d

elif [ "$1" == "up" ]; then

    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml up -d

elif [ "$1" == "down" ]; then

    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml -f "$composeFilePath"/importer/docker-compose.config.yml stop

elif [ "$1" == "destroy" ]; then

    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml -f "$composeFilePath"/importer/docker-compose.config.yml down -v

    docker volume rm opencr-data instant_elasticsearch-data

else
    echo "Valid options are: init, up, down, or destroy"
fi