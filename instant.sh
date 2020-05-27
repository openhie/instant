#!/bin/bash

COMMAND=$1
TARGET=$2

if ! [[ "$TARGET" =~ ^(docker|kubernetes|k8s)$ ]]; then
    TARGET=docker
    echo "Defaulting to docker as a target, target either not specified or invalid"
fi

if [ "$TARGET" == "docker" ]; then
    if [ "$COMMAND" == "init" ]; then
        ./core/docker/compose.sh init
        ./healthworkforce/docker/compose.sh init
    elif [ "$COMMAND" == "up" ]; then
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
        echo "Valid options are: init, up, down, test or destroy"
    fi
fi


if [ "$TARGET" == "kubernetes" ] || [ "$TARGET" == "k8s" ]; then
    envContextName=$(kubectl config get-contexts | grep '*' | awk '{print $2}')
    printf "\n\n>>> Applying to the '${envContextName}' context <<<\n\n\n"

    if [ "$COMMAND" == "init" ]; then
        ./core/kubernetes/main/k8s.sh init
        ./core/kubernetes/importer/k8s.sh up
        ./healthworkforce/kubernetes/main/k8s.sh up
        ./healthworkforce/kubernetes/importer/k8s.sh up
    elif [ "$COMMAND" == "up" ]; then
        ./core/kubernetes/main/k8s.sh up
        ./healthworkforce/kubernetes/main/k8s.sh up
    elif [ "$COMMAND" == "down" ]; then
        ./core/kubernetes/main/k8s.sh down
        ./healthworkforce/kubernetes/main/k8s.sh down
    elif [ "$COMMAND" == "destroy" ]; then
        ./core/kubernetes/main/k8s.sh destroy
        ./healthworkforce/kubernetes/main/k8s.sh destroy
    elif [ "$COMMAND" == "test" ]; then
        openhimCoreHostname=$(kubectl get service openhim-core-service -o=jsonpath="{.status.loadBalancer.ingress[*]['hostname', 'ip']}")
        openhimCoreTransactionSSLPort=$(kubectl get service openhim-core-service -o=jsonpath={.spec.ports[1].port})
        hostnameLength=$(expr length "$openhimCoreHostname")

        if [ "$hostnameLength" -le 0 ]; then
            openhimCoreHostname=$(minikube ip)
            openhimCoreTransactionSSLPort=$(kubectl get service openhim-core-service -o=jsonpath={.spec.ports[1].nodePort})
        fi

        ./core/test.sh $openhimCoreHostname:$openhimCoreTransactionSSLPort
        ./healthworkforce/test.sh $openhimCoreHostname:$openhimCoreTransactionSSLPort
    else
        echo "Valid options are: init, up, down, test or destroy"
    fi
fi
