#!/bin/bash

k8sMainRootFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "up" ]; then
    kubectl apply -k $k8sMainRootFilePath
elif [ "$1" == "down" ]; then
    kubectl delete deployment mapper-deployment
elif [ "$1" == "destroy" ]; then
    kubectl delete deployments,services -l package=hwf
else
    echo "Valid options are: up, down, or destroy"
fi
