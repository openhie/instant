#!/bin/bash

if [ "$1" == "up" ]; then
    kubectl apply -k .
    kubectl get jobs
elif [ "$1" == "clean" ]; then
    kubectl delete -k .
else
    echo "Valid options are: up, or clean"
fi
