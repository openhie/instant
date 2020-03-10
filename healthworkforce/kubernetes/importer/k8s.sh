#!/bin/bash

k8sImporterRootFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "up" ]; then
    kubectl apply -k $k8sImporterRootFilePath
    kubectl get jobs --namespace=hwf-package
elif [ "$1" == "clean" ]; then
    kubectl delete -k $k8sImporterRootFilePath
else
    echo "Valid options are: up, or clean"
fi
