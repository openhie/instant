#!/bin/bash

kustomizationFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "up" ]; then
    kubectl apply -k $kustomizationFilePath

    kubectl get services
    kubectl get ingress

elif [ "$1" == "down" ]; then
    kubectl delete deployment openhim-console-deployment
    kubectl delete deployment openhim-core-deployment
    kubectl delete deployment openhim-mongo-deployment
    kubectl delete deployment hapi-fhir-server-deployment
    kubectl delete deployment hapi-fhir-mysql-deployment
elif [ "$1" == "destroy" ]; then
    kubectl delete -k $kustomizationFilePath
else
    echo "Valid options are: up, down, or destroy"
fi
