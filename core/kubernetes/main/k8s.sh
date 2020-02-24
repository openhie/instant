#!/bin/bash

if [ "$1" == "up" ]; then
    kubectl apply -f ./openhim
    kubectl apply -f ./hapi-fhir

    if [ "$2" == "dev" ]; then
        kubectl apply -f ../dev/openhim
        kubectl apply -f ../dev/hapi-fhir
        
        echo -e "\nCurrently in development mode\n"
    fi

    kubectl get services
    kubectl get ingress

    # create HOST entry for ingress
    sudo sed -i "/HOST alias for kubernetes Minikube/d" /etc/hosts
    echo -e "\n$(minikube ip) $(kubectl get ingress -o jsonpath="{..host}") # HOST alias for kubernetes Minikube" | sudo tee -a /etc/hosts
elif [ "$1" == "down" ]; then
    kubectl delete deployment openhim-console-deployment
    kubectl delete deployment openhim-core-deployment
    kubectl delete deployment hapi-fhir-server-deployment
    kubectl delete deployment hapi-fhir-mysql-deployment
elif [ "$1" == "destroy" ]; then
    kubectl delete deployment --all --grace-period 2
    kubectl delete service --all --grace-period 2
    kubectl delete configmap --all --grace-period 2
    kubectl delete persistentvolumeclaims --all --grace-period 2
    kubectl delete cronjobs --all --grace-period 2
    kubectl delete pods --all --grace-period 2
    kubectl delete ingress --all --grace-period 2
else
    echo "Valid options are: up, down, or destroy"
fi
