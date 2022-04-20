#!/bin/bash

composeFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

# Set default password which can be overwritten if there the env var is already set
export ES_LOGSTASH_SYSTEM=${ES_LOGSTASH_SYSTEM:-dev_password_only}
export ES_APM_SYSTEM=${ES_APM_SYSTEM:-dev_password_only}
export ES_REMOTE_MONITORING_USER=${ES_REMOTE_MONITORING_USER:-dev_password_only}
export ES_ELASTIC=${ES_ELASTIC:-dev_password_only}
export ES_KIBANA_SYSTEM=${ES_KIBANA_SYSTEM:-dev_password_only}
export ES_BEATS_SYSTEM=${ES_BEATS_SYSTEM:-dev_password_only}

if [ "$1" == "init" ]; then
    docker-compose -p instant -f "$composeFilePath"/docker-compose.yml up -d es-analytics

    # echo "Waiting for elasticsearch to start before automatically setting built-in passwords..."
    sleep 40
    apt-get install -y expect >/dev/null 2>&1
    echo "Setting passwords..."
    "$composeFilePath"/set-pwds.exp
    echo "Done"

    echo "Setting passwords in config files..."
    sed -i "s/dev_password_only/$ES_ELASTIC/g" "$composeFilePath"/../../elastic-pipeline/docker/pipeline/*.conf
    sed -i "s/dev_password_only/$ES_ELASTIC/g" "$composeFilePath"/../../elastic-pipeline/docker/logstash.yml
    echo "Done"

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
