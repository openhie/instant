#!/bin/bash

kustomizationFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
openhimConsoleVolumePath="${kustomizationFilePath}/openhim/volume/openhim-console/default.json"
openhimApiUrl=''
fhirServerUrl=''
openhimConsoleUrl=''

if [ "$1" == "up" ]; then
    kubectl apply -k $kustomizationFilePath

    coreUrlLength=0

    while [ $coreUrlLength -le 0 ];
    do
        echo "OpenHIM Core not ready. Sleep 10"
        sleep 10
        openhimApiUrl=$(kubectl get service openhim-core-service -o=jsonpath={.status.loadBalancer.ingress[0].hostname})
        coreUrlLength=$(expr length "$openhimApiUrl")
    done

    printf "\n\nOpenHIM Api Url\n---------------\n"$openhimApiUrl"\n\n\n"

    fhirUrlLength=$(expr length "$fhirServerUrl")

    while [ $fhirUrlLength -le 0 ];
    do
        echo "HAPI-FHIR not ready. Sleep 5"
        sleep 5
        fhirServerUrl=$(kubectl get service hapi-fhir-server-service -o=jsonpath={.status.loadBalancer.ingress[0].hostname})
        fhirUrlLength=$(expr length "$fhirServerUrl")
    done

    printf "\nHAPI FHIR Url\n--------------\n"$fhirServerUrl"\n\n\n"

    # Injecting OpenHIM Core Api url into Console config file
    sed -i -E "s/(\"host\": \")\S*(\")/\1${openhimApiUrl}\2/" $openhimConsoleVolumePath
    # Injecting OpenHIM Core port into Console config file
    sed -i -E "s/(\"port\": )\S*(,)/\18082\2/" $openhimConsoleVolumePath

    kubectl apply -k $kustomizationFilePath/openhim

    consoleUrlLength=0

    while [ $consoleUrlLength -le 0 ];
    do
        echo "OpenHIM Console not ready. Sleep 5"
        sleep 5
        openhimConsoleUrl=$(kubectl get service openhim-console-service -o=jsonpath={.status.loadBalancer.ingress[0].hostname})
        consoleUrlLength=$(expr length "$openhimConsoleUrl")
    done

    printf "\n\nOpenHIM Console Url\n-------------------\nhttp://"$openhimConsoleUrl"\n\n"
    printf "AWS is creating "

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
