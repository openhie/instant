#!/bin/bash

k8sMainRootFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "init" ]; then

    kubectl apply -k $k8sMainRootFilePath

elif [ "$1" == "up" ]; then

    kubectl apply -k $k8sMainRootFilePath

elif [ "$1" == "down" ]; then

    kubectl delete -k $k8sMainRootFilePath

elif [ "$1" == "destroy" ]; then

    kubectl delete -k $k8sMainRootFilePath

else
    echo "Valid options are: init, up, down, or destroy"
fi
