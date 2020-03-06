#!/bin/bash

kustomizationFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "up" ]; then
    # Create the namespace
    kubectl apply -f $kustomizationFilePath/healthworkforce-namespace.yaml
    kubectl apply -k $kustomizationFilePath
elif [ "$1" == "down" ]; then
    kubectl delete deployment mapper-deployment
elif [ "$1" == "destroy" ]; then
    kubectl delete namespaces healthworkforce-component
else
    echo "Valid options are: up, down, or destroy"
fi
