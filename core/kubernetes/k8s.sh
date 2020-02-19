#!/bin/bash

if [ "$1" == "up" ]; then
    kubectl apply -f ./openhim
    kubectl get services
elif [ "$1" == "down" ]; then
    kubectl delete deployment openhim-console-deployment
    kubectl delete deployment openhim-core-deployment
elif [ "$1" == "destroy" ]; then
    kubectl delete deployment --all --grace-period 2
    kubectl delete service --all --grace-period 2
    kubectl delete configmap --all --grace-period 2
    kubectl delete persistentvolumeclaims --all --grace-period 2
    kubectl delete cronjobs --all --grace-period 2
    kubectl delete pods --all --grace-period 2
else
    echo "Valid options are: up, down, or destroy"
fi