#!/bin/bash

kustomizationFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "up" ]; then
    kubectl apply -k $kustomizationFilePath
    kubectl get jobs --namespace=core-package
elif [ "$1" == "clean" ]; then
    kubectl delete -k $kustomizationFilePath
else
    echo "Valid options are: up, or clean"
fi
