#!/bin/bash

composeFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "init" ]; then

    docker create --name openelisglobal-database-helper -v openelis-db-data:/database/dbInit busybox
    docker cp "$composeFilePath"/importer/volume/database/dbInit/. openelisglobal-database-helper:/database/dbInit/
    docker rm openelisglobal-database-helper

    docker create --name openelis-plugins-helper -v openelis-plugins-data:/plugins/ busybox
    docker cp "$composeFilePath"/importer/volume/plugins/. openelis-plugins-helper:/plugins/
    docker rm openelis-plugins-helper

    docker create --name openelisglobal-webapp-helper -v openelis-server-data:/tomcat busybox
    docker cp "$composeFilePath"/importer/volume/tomcat/oe_server.xml openelisglobal-webapp-helper:/tomcat/server.xml
    docker rm openelisglobal-webapp-helper

    docker create --name external-fhir-api-helper -v openelis-fhir-data:/tomcat busybox
    docker cp "$composeFilePath"/importer/volume/tomcat/hapi_server.xml external-fhir-api-helper:/tomcat/server.xml
    docker rm external-fhir-api-helper

    # docker secret create  datasource.password "$composeFilePath"/importer/volume/properties/datasource.password
    # docker secret create  common.properties "$composeFilePath"/importer/volume/properties/common.properties

    docker-compose -f "$composeFilePath"/docker-compose.yml up 

elif [ "$1" == "up" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml up -d

elif [ "$1" == "down" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml stop

elif [ "$1" == "destroy" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml down -v

else
    echo "Valid options are: init, up, down, or destroy"
fi

