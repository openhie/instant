#!/bin/bash

composeFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$PROD" == "true" ]; then
    devComposeParam=""
else
    printf "\nRunning core package in DEV mode\n"
    devComposeParam="-f ${composeFilePath}/docker-compose.dev.yml"
    devComposeMongoParam="-f ${composeFilePath}/docker-compose.dev-mongo.yml"
fi

if [ "$1" == "init" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose-mongo.yml $devComposeMongoParam up -d

    # Set up the replica set
    "$composeFilePath"/initiateReplicaSet.sh

    docker create --name hapi-mysql-helper -v hapi-mysql-config:/conf.d busybox
    docker cp "$composeFilePath"/importer/volume/mysql.cnf hapi-mysql-helper:/conf.d/mysql.cnf
    docker rm hapi-mysql-helper

    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml -f "$composeFilePath"/importer/docker-compose.config.yml $devComposeParam up -d
elif [ "$1" == "up" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose-mongo.yml $devComposeMongoParam up -d

    # Wait for mongo replica set to be set up
    sleep 20

    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml $devComposeParam up -d
elif [ "$1" == "down" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose-mongo.yml -f "$composeFilePath"/docker-compose.yml -f "$composeFilePath"/importer/docker-compose.config.yml stop
elif [ "$1" == "destroy" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose-mongo.yml -f "$composeFilePath"/docker-compose.yml -f "$composeFilePath"/importer/docker-compose.config.yml down -v
else
    echo "Valid options are: init, up, down, or destroy"
fi
