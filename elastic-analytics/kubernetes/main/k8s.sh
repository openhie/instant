#!/bin/bash

k8sMainRootFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

# Set default password which can be overwritten if there the env var is already set
export ES_LOGSTASH_SYSTEM=${ES_LOGSTASH_SYSTEM:-dev_password_only}
export ES_APM_SYSTEM=${ES_APM_SYSTEM:-dev_password_only}
export ES_REMOTE_MONITORING_USER=${ES_REMOTE_MONITORING_USER:-dev_password_only}
export ES_ELASTIC=${ES_ELASTIC:-dev_password_only}
export ES_KIBANA_SYSTEM=${ES_KIBANA_SYSTEM:-dev_password_only}
export ES_BEATS_SYSTEM=${ES_BEATS_SYSTEM:-dev_password_only}

if [ "$1" == "init" ]; then

    kubectl apply -k $k8sMainRootFilePath

    echo "Waiting for elasticsearch to start before automatically setting built-in passwords..."
    sleep 40
    apt-get install -y expect >/dev/null 2>&1
    echo "Setting passwords..."
    "$k8sMainRootFilePath"/set-pwds.exp
    echo "Done"

elif [ "$1" == "up" ]; then

    kubectl apply -k $k8sMainRootFilePath

elif [ "$1" == "down" ]; then

    kubectl delete -k $k8sMainRootFilePath

elif [ "$1" == "destroy" ]; then

    kubectl delete -k $k8sMainRootFilePath

else
    echo "Valid options are: init, up, down, or destroy"
fi
