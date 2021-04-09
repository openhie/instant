#!/bin/bash

k8sMainRootFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "init" ]; then
    kubectl apply -k $k8sMainRootFilePath
elif [ "$1" == "up" ]; then
    kubectl apply -k $k8sMainRootFilePath
elif [ "$1" == "down" ]; then
    kubectl delete deployment dhis-web dhis-postgres
elif [ "$1" == "destroy" ]; then
    kubectl delete deployments,services,jobs,persistentvolumeclaim,configmap -l package=dhis
else
    echo "Valid options are: up, down, or destroy"
fi
