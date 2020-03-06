#!/bin/bash

COMMAND=$1
TARGET=$2

if ! [[ "$TARGET" =~ ^(docker|kubernetes|k8s)$ ]]; then
    TARGET=docker
    echo "Defaulting to docker as a target, target either not specified or invalid"
fi

if [ "$TARGET" == "docker" ]; then
    if [ "$COMMAND" == "up" ]; then
        ./core/docker/compose.sh up
        ./healthworkforce/docker/compose.sh up
    elif [ "$COMMAND" == "down" ]; then
        ./core/docker/compose.sh down
        ./healthworkforce/docker/compose.sh down
    elif [ "$COMMAND" == "destroy" ]; then
        ./core/docker/compose.sh destroy
        ./healthworkforce/docker/compose.sh destroy
    elif [ "$COMMAND" == "test" ]; then
        ./core/test.sh localhost:5000
        ./healthworkforce/test.sh localhost:5000
    else
        echo "Valid options are: up, down, test or destroy"
    fi
fi


if [ "$TARGET" == "kubernetes" ] || [ "$TARGET" == "k8s" ]; then
    minikubeIP=$(minikube ip)
    openhimCoreTransactionSSLPort=$(kubectl get service openhim-core-service --namespace=core-component -o=jsonpath={.spec.ports[1].nodePort})

    if [ "$COMMAND" == "up" ]; then
        ./core/kubernetes/main/k8s.sh up
        ./core/kubernetes/importer/k8s.sh up
        ./healthworkforce/kubernetes/main/k8s.sh up
        ./healthworkforce/kubernetes/importer/k8s.sh up
    elif [ "$COMMAND" == "down" ]; then
        ./core/kubernetes/main/k8s.sh down
        ./healthworkforce/kubernetes/main/k8s.sh down
    elif [ "$COMMAND" == "destroy" ]; then
        ./core/kubernetes/main/k8s.sh destroy
        ./healthworkforce/kubernetes/main/k8s.sh destroy
    elif [ "$COMMAND" == "test" ]; then
        ./core/test.sh $minikubeIP:$openhimCoreTransactionSSLPort
        ./healthworkforce/test.sh $minikubeIP:$openhimCoreTransactionSSLPort
    else
        echo "Valid options are: up, down, test or destroy"
    fi
fi
