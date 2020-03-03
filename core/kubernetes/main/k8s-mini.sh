#!/bin/bash

kustomizationFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "up" ]; then
    minikube addons enable ingress
    kubectl apply -k $kustomizationFilePath

    kubectl get services
    kubectl get ingress

    # create HOST entry for ingress
    sudo sed -i "/HOST alias for kubernetes Minikube/d" /etc/hosts
    echo -e "\n$(minikube ip) $(kubectl get ingress -o jsonpath="{..host}") # HOST alias for kubernetes Minikube" | sudo tee -a /etc/hosts
elif [ "$1" == "down" ]; then
    kubectl delete deployment openhim-console-deployment
    kubectl delete deployment openhim-core-deployment
    kubectl delete deployment openhim-mongo-deployment
    kubectl delete deployment hapi-fhir-server-deployment
    kubectl delete deployment hapi-fhir-mysql-deployment
elif [ "$1" == "destroy" ]; then
    # delete host entry on destroy
    sudo sed -i "/HOST alias for kubernetes Minikube/d" /etc/hosts

    kubectl delete -k $kustomizationFilePath
else
    echo "Valid options are: up, down, or destroy"
fi
