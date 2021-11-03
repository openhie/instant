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
    chmod -R a+rwx "$composeFilePath"/importer/volume/tomcat/openelis/conf/
    docker cp "$composeFilePath"/importer/volume/tomcat/openelis/conf/. openelisglobal-webapp-helper:/tomcat/
    docker rm openelisglobal-webapp-helper

    docker create --name openelisglobal-properties-helper -v openelis-properties-data:/properties busybox
    docker cp "$composeFilePath"/importer/volume/properties/common.properties openelisglobal-properties-helper:/properties/common.properties
    docker rm openelisglobal-properties-helper

    docker create --name external-fhir-api-helper -v openelis-fhir-data:/tomcat busybox
    chmod -R a+rwx "$composeFilePath"/importer/volume/tomcat/fhir_server/conf/
    docker cp "$composeFilePath"/importer/volume/tomcat/fhir_server/conf/. external-fhir-api-helper:/tomcat/ 
    docker rm external-fhir-api-helper

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

