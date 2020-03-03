#!/bin/bash

kustomizationFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
openhimConsoleVolumePath="${kustomizationFilePath}/openhim/volume/openhim-console/default.json"
minikubeIP=''

if [ "$1" == "up" ]; then
    kubectl apply -k $kustomizationFilePath

    minikubeIP=$(minikube ip)
    corePort=$(kubectl get service openhim-core-service -o=jsonpath={.spec.ports[0].nodePort})

    # Injecting minikube ip as the hostname of the OpenHIM Core into Console config file
    sed -i -E "s/(\"host\": \")\S*(\")/\1${minikubeIP}\2/" $openhimConsoleVolumePath

    # Injecting OpenHIM Core port into Console config file
    sed -i -E "s/(\"port\": )\S*(,)/\1${corePort}\2/" $openhimConsoleVolumePath

    kubectl apply -k $kustomizationFilePath/openhim

    consolePort=$(kubectl get service openhim-console-service -o=jsonpath={.spec.ports[0].nodePort})

    printf "\n\nOpenHIM Console Url\n-------------------\nhttp://"$minikubeIP":"$consolePort"\n\n"
elif [ "$1" == "down" ]; then
    kubectl delete deployment openhim-console-deployment
    kubectl delete deployment openhim-core-deployment
    kubectl delete deployment openhim-mongo-deployment
    kubectl delete deployment hapi-fhir-server-deployment
    kubectl delete deployment hapi-fhir-mysql-deployment
elif [ "$1" == "destroy" ]; then
    kubectl delete -k $kustomizationFilePath
    kubectl delete -k $kustomizationFilePath/openhim
else
    echo "Valid options are: up, down, or destroy"
fi
