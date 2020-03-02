#!/bin/bash

kustomizationFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
openhimConsoleVolumePath="${kustomizationFilePath}/openhim/volume/openhim-console/default.json"
openhimApiUrl=''
fhirServerUrl=''
openhimConsoleUrl=''

if [ "$1" == "up" ]; then
    kubectl apply -k $kustomizationFilePath

    openhimApiUrl=$(kubectl get service openhim-core-service -o=jsonpath={.status.loadBalancer.ingress[0].hostname})
    fhirServerUrl=$(kubectl get service hapi-fhir-server-service -o=jsonpath={.status.loadBalancer.ingress[0].hostname})

    printf "\n\nOpenHIM Api Url\n---------------\n"$openhimApiUrl"\n"
    printf "\nHAPI FHIR Url\n--------------\n"$fhirServerUrl"\n\n\n"

    sed -i -E "s/(\"host\": \")\S*(\")/\1${openhimApiUrl}\2/" $openhimConsoleVolumePath
    sed -i -E "s/(\"port\": )\S*(,)/\18082\2/" $openhimConsoleVolumePath

    kubectl apply -k $kustomizationFilePath/openhim

    openhimConsoleUrl=$(kubectl get service openhim-console-service -o=jsonpath={.status.loadBalancer.ingress[0].hostname})

    printf "\n\nOpenHIM Console Url\n-------------------\nhttp://"$openhimConsoleUrl"\n\n"

elif [ "$1" == "down" ]; then
    kubectl delete deployment openhim-console-deployment
    kubectl delete deployment openhim-core-deployment
    kubectl delete deployment openhim-mongo-deployment
    kubectl delete deployment hapi-fhir-server-deployment
    kubectl delete deployment hapi-fhir-mysql-deployment
elif [ "$1" == "destroy" ]; then
    kubectl delete -k $kustomizationFilePath
else
    echo "Valid options are: up, down, or destroy"
fi
