#!/bin/bash

kustomizationFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
openhimConsoleVolumePath="${kustomizationFilePath}/openhim/volume/openhim-console/default.json"

hapiFhirServerUrl=''
openhimConsoleUrl=''
openhimCoreMediatorApiUrl=''
openhimCoreTransactionApiUrl=''
openhimCoreTransactionSSLApiUrl=''

hapiFhirPort='8080'
openhimConsolePort='80'
openhimCoreMediatorSSLPort='8082'
openhimCoreTransactionPort='5001'
openhimCoreTransactionSSLPort='5000'

cloud_setup () {
    coreUrlLength=$(expr length "$openhimCoreHostname")

    while [ $coreUrlLength -le 0 ];
    do
        echo "OpenHIM Core not ready. Sleep 10"
        sleep 10
        openhimCoreHostname=$(kubectl get service openhim-core-service --namespace=core-component -o=jsonpath={.status.loadBalancer.ingress[0].hostname})
        coreUrlLength=$(expr length "$openhimCoreHostname")
    done

    openhimCoreMediatorApiUrl="https://$openhimCoreHostname:$openhimCoreMediatorSSLPort"
    openhimCoreTransactionApiUrl="http://$openhimCoreHostname:$openhimCoreTransactionPort"
    openhimCoreTransactionSSLApiUrl="https://$openhimCoreHostname:$openhimCoreTransactionSSLPort"

    # Injecting OpenHIM Core Api url into Console config file
    sed -i -E "s/(\"host\": \")\S*(\")/\1${openhimCoreHostname}\2/" $openhimConsoleVolumePath
    # Injecting OpenHIM Core port into Console config file
    sed -i -E "s/(\"port\": )\S*(,)/\1${openhimCoreMediatorSSLPort}\2/" $openhimConsoleVolumePath

    kubectl apply -k $kustomizationFilePath/openhim

    fhirUrlLength=$(expr length "$hapiFhirServerHostname")

    while [ $fhirUrlLength -le 0 ];
    do
        echo "HAPI-FHIR not ready. Sleep 5"
        sleep 5
        hapiFhirServerHostname=$(kubectl get service hapi-fhir-server-service --namespace=core-component -o=jsonpath={.status.loadBalancer.ingress[0].hostname})
        fhirUrlLength=$(expr length "$hapiFhirServerHostname")
    done

    hapiFhirServerUrl="http://$hapiFhirServerHostname:$hapiFhirPort"

    consoleUrlLength=$(expr length "$openhimConsoleHostname")

    while [ $consoleUrlLength -le 0 ];
    do
        echo "OpenHIM Console not ready. Sleep 5"
        sleep 5
        openhimConsoleHostname=$(kubectl get service openhim-console-service --namespace=core-component -o=jsonpath={.status.loadBalancer.ingress[0].hostname})
        consoleUrlLength=$(expr length "$openhimConsoleHostname")
    done

    openhimConsoleUrl="http://$openhimConsoleHostname:$openhimConsolePort"
}

local_setup () {
    minikubeIP=$(minikube ip)
    openhimCoreMediatorSSLPort=$(kubectl get service openhim-core-service --namespace=core-component -o=jsonpath={.spec.ports[0].nodePort})
    openhimCoreTransactionPort=$(kubectl get service openhim-core-service --namespace=core-component -o=jsonpath={.spec.ports[2].nodePort})
    openhimCoreTransactionSSLPort=$(kubectl get service openhim-core-service --namespace=core-component -o=jsonpath={.spec.ports[1].nodePort})
    hapiFhirPort=$(kubectl get service hapi-fhir-server-service --namespace=core-component -o=jsonpath={.spec.ports[0].nodePort})

    hapiFhirServerUrl="http://$minikubeIP:$hapiFhirPort"
    openhimCoreMediatorApiUrl="https://$minikubeIP:$openhimCoreMediatorSSLPort"
    openhimCoreTransactionApiUrl="http://$minikubeIP:$openhimCoreTransactionPort"
    openhimCoreTransactionSSLApiUrl="https://$minikubeIP:$openhimCoreTransactionSSLPort"

    # Injecting minikube ip as the hostname of the OpenHIM Core into Console config file
    sed -i -E "s/(\"host\": \")\S*(\")/\1${minikubeIP}\2/" $openhimConsoleVolumePath

    # Injecting OpenHIM Core port into Console config file
    sed -i -E "s/(\"port\": )\S*(,)/\1${openhimCoreMediatorSSLPort}\2/" $openhimConsoleVolumePath

    kubectl apply -k $kustomizationFilePath/openhim

    openhimConsolePort=$(kubectl get service openhim-console-service --namespace=core-component -o=jsonpath={.spec.ports[0].nodePort})

    openhimConsoleUrl="http://$minikubeIP:$openhimConsolePort"
}

if [ "$1" == "up" ]; then
    envContextName=$(kubectl config get-contexts | grep '*' | awk '{print $2}')

    printf "\n\n>>> Deploying to the '${envContextName}' context <<<\n\n\n"

    # Create the component's namespace
    kubectl apply -f $kustomizationFilePath/core-namespace.yaml
    kubectl apply -k $kustomizationFilePath

    envContextMinikube=$(echo $envContextName | grep 'minikube')

    if [ $(expr length "$envContextMinikube") -le 0 ]; then
        cloud_setup
    else
        local_setup
    fi

    printf "\n\nHAPI FHIR Server Url\n--------------------\n"$hapiFhirServerUrl"\n\n"
    printf "OpenHIM Mediator API Url\n------------------------\n"$openhimCoreMediatorApiUrl"\n\n"
    printf "OpenHIM Transaction API Url\n---------------------------\n"$openhimCoreTransactionApiUrl"\n\n"
    printf "OpenHIM Transaction SSL API Url\n-------------------------------\n"$openhimCoreTransactionSSLApiUrl"\n\n"
    printf "OpenHIM Console Url\n===================\n"$openhimConsoleUrl"\n\n"
    printf ">>> The OpenHIM Console Url will take a few minutes to become active <<<\n\n"
elif [ "$1" == "down" ]; then
    kubectl delete deployment openhim-console-deployment
    kubectl delete deployment openhim-core-deployment
    kubectl delete deployment openhim-mongo-deployment
    kubectl delete deployment hapi-fhir-server-deployment
    kubectl delete deployment hapi-fhir-mysql-deployment
elif [ "$1" == "destroy" ]; then
    kubectl delete namespaces core-component
else
    echo "Valid options are: up, down, or destroy"
fi
