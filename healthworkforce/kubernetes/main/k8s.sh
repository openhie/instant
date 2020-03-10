#!/bin/bash

k8sMainRootFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "up" ]; then
    # Create the namespace
    kubectl apply -f $k8sMainRootFilePath/healthworkforce-namespace.yaml
    kubectl apply -k $k8sMainRootFilePath
elif [ "$1" == "down" ]; then
    kubectl delete deployment mapper-deployment
elif [ "$1" == "destroy" ]; then
    kubectl delete namespaces hfw-package
else
    echo "Valid options are: up, down, or destroy"
fi
