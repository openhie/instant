#!/bin/bash

if [ "$1" == "up" ]; then
    ./core/kubernetes/main/k8s.sh up
    ./core/kubernetes/importer/k8s.sh up
    ./healthworkforce/kubernetes/main/k8s.sh up
elif [ "$1" == "down" ]; then
    ./core/kubernetes/main/k8s.sh down
    ./healthworkforce/kubernetes/main/k8s.sh down
elif [ "$1" == "destroy" ]; then
    ./core/kubernetes/main/k8s.sh destroy
    ./core/kubernetes/importer/k8s.sh clean
    ./healthworkforce/kubernetes/main/k8s.sh destroy
elif [ "$1" == "test" ]; then
    ./core/test.sh
else
    echo "Valid options are: up, down, test or destroy"
fi
