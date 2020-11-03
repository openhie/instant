#!/bin/bash

composeFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "init" ]; then
    docker-compose -p facility -f "$composeFilePath"/docker-compose.yml up -d

    # Set up the openhim
    # "$composeFilePath"/initiateReplicaSet.sh

elif [ "$1" == "up" ]; then
    docker-compose -p facility -f "$composeFilePath"/docker-compose.yml up -d

    # Wait
    # sleep 20

    # do something else needed

elif [ "$1" == "down" ]; then
    docker-compose -p facility -f "$composeFilePath"/docker-compose.yml stop

elif [ "$1" == "destroy" ]; then
    docker-compose -p facility -f "$composeFilePath"/docker-compose.yml down -v

else
    echo "Valid options are: init, up, down, or destroy"
fi

