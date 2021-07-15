#!/bin/bash

composeFilePath=$(dirname $(readlink -f $0))

# Create a volume and copy the project configuration in before starting the stack.
docker volume create microservice_him_hapi
docker container create --name dummy -v microservice_him_hapi:/mnt tianon/true 
docker cp $composeFilePath/config/project.yaml dummy:/mnt/
docker rm dummy

if [ "$1" == "init" ]; then
    docker-compose -p openfn \
    -f "$composeFilePath"/docker-compose.yml \
    -f "$composeFilePath"/docker-compose.dev.yml \
    -f "$composeFilePath"/docker-compose.config.yml \
    up -d
elif [ "$1" == "up" ]; then
    docker-compose -p openfn \
    -f "$composeFilePath"/docker-compose.yml \
    -f "$composeFilePath"/docker-compose.dev.yml \
    -f "$composeFilePath"/docker-compose.config.yml \
    up -d
elif [ "$1" == "down" ]; then
    docker-compose -p openfn \
    -f "$composeFilePath"/docker-compose.yml \
    -f "$composeFilePath"/docker-compose.dev.yml \
    -f "$composeFilePath"/docker-compose.config.yml \
    stop
elif [ "$1" == "destroy" ]; then
    docker-compose -p openfn \
    -f "$composeFilePath"/docker-compose.yml \
    -f "$composeFilePath"/docker-compose.dev.yml \
    -f "$composeFilePath"/docker-compose.config.yml \
    down
else
    echo "Valid options are: init, up, down, or destroy"
fi
