#!/bin/bash

k8sMainRootFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

DHIS2ServerUrl=''
DHIS2Port=''

cloud_setup () {
    DHIS2Port=$(kubectl get service dhis-web -o=jsonpath={.spec.ports[0].port})

    while
        DHIS2ServerHostname=$(kubectl get service dhis-web -o=jsonpath="{.status.loadBalancer.ingress[0]['hostname', 'ip']}")
        DHISUrlLength=$(expr length "$DHIS2ServerHostname")
        (( DHISUrlLength <= 0 ))
    do
        echo "DHIS2 not ready. Sleep 5"
        sleep 5
    done

    DHIS2ServerUrl="http://$DHIS2ServerHostname:$DHIS2Port"
}

local_setup () {
    minikubeIP=$(kubectl config view -o=jsonpath='{.clusters[?(@.name=="minikube")].cluster.server}' | awk '{ split($0,A,/:\/*/) ; print A[2] }')
    DHIS2Port=$(kubectl get service dhis-web -o=jsonpath={.spec.ports[0].nodePort})

    DHIS2ServerUrl="http://$minikubeIP:$DHIS2Port"
}

print_service_url () {
    printf "\n\nDHIS2 Server Url\n--------------------\n"$DHIS2ServerUrl"\n\n"
}

if [ "$1" == "init" ]; then
    kubectl apply -k $k8sMainRootFilePath

    envContextName=$(kubectl config get-contexts | grep '*' | awk '{print $2}')
    envContextMinikube=$(echo $envContextName | grep 'minikube')

    if [ $(expr length "$envContextMinikube") -le 0 ]; then
        cloud_setup
    else
        local_setup
    fi

    print_service_url
    printf ">>> DHIS2 will take a few minutes to become active <<<\n\n"

elif [ "$1" == "up" ]; then
    kubectl apply -k $k8sMainRootFilePath

    envContextName=$(kubectl config get-contexts | grep '*' | awk '{print $2}')
    envContextMinikube=$(echo $envContextName | grep 'minikube')

    if [ $(expr length "$envContextMinikube") -le 0 ]; then
        cloud_setup
    else
        local_setup
    fi

    print_service_url
elif [ "$1" == "down" ]; then
    kubectl delete deployment dhis-web dhis-postgres
elif [ "$1" == "destroy" ]; then
    kubectl delete deployments,services,jobs,persistentvolumeclaim,configmap -l package=dhis
else
    echo "Valid options are: up, down, or destroy"
fi
