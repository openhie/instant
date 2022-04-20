#!/bin/bash

composeFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$PROD" == "true" ]; then
    printf "\nRunning core package in PROD mode\n"
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

    PASSWORD_SALT=$(openssl rand -hex 16)
    PASSWORD_HASH=$(echo -n $PASSWORD_SALT${OPENHIM_ROOT_PASSWORD-instant101} | openssl sha512 | awk '{print $2}')

    cp "$composeFilePath"/importer/volume/openhim-import-temp.json "$composeFilePath"/importer/volume/openhim-import.json
    sed -i "s/{{PASSWORD_SALT}}/$PASSWORD_SALT/g" "$composeFilePath"/importer/volume/openhim-import.json
    sed -i "s/{{PASSWORD_HASH}}/$PASSWORD_HASH/g" "$composeFilePath"/importer/volume/openhim-import.json
    sed -i "s/{{OPENHIM_CLIENT_TOKEN}}/${OPENHIM_CLIENT_TOKEN-test}/g" "$composeFilePath"/importer/volume/openhim-import.json

    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml -f "$composeFilePath"/importer/docker-compose.config.yml $devComposeParam up -d

    echo "Sleep 30s to allow core config importer to complete"
    sleep 30

    docker rm core-config-importer
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
