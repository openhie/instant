#!/bin/bash

composeFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "init" ]; then

    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.hapi.yml up -d
    echo "Waiting for HAPI to start"
    sleep 60
    echo "Loading HL7 structure definitions..."
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.hapi.config.yml up -d
    sleep 60
    echo "Loading custom structure definitions..."
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.ihris.config.yml up -d
    sleep 60
    echo "Loading demo data..."
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.ihris.data.yml up -d
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.elastic.yml up -d
    echo "Launching iHRIS"
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.ihris.yml up -d


elif [ "$1" == "up" ]; then

    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.hapi.yml up -d
    echo "Waiting for HAPI to start"
    sleep 60
    echo "Loading HL7 structure definitions..."
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.hapi.config.yml up -d
    sleep 60
    echo "Loading custom structure definitions..."
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.ihris.config.yml up -d
    sleep 60
    echo "Loading demo data..."
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.ihris.data.yml up -d
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.elastic.yml up -d
    echo "Launching iHRIS"
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.ihris.yml up -d

elif [ "$1" == "down" ]; then

    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.hapi.yml stop
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.elastic.yml stop
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.ihris.yml stop


elif [ "$1" == "destroy" ]; then

    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.hapi.yml down -v
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.elastic.yml down -v
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.ihris.yml down -v

    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.hapi.config.yml down -v
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.ihris.config.yml down -v
    docker-compose -p healthworker -f "$composeFilePath"/docker-compose.ihris.data.yml down -v

else
    echo "Valid options are: init, up, down, or destroy"
fi
