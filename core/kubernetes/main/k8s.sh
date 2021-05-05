#!/bin/bash

k8sMainRootFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
openhimConsoleVolumePath="${k8sMainRootFilePath}/openhim/volume/openhim-console/default.json"

hapiFhirServerUrl=''
openhimConsoleUrl=''
openhimCoreMediatorApiUrl=''
openhimCoreTransactionApiUrl=''
openhimCoreTransactionSSLApiUrl=''

hapiFhirPort=''
openhimConsolePort=''
openhimCoreMediatorSSLPort=''
openhimCoreTransactionPort=''
openhimCoreTransactionSSLPort=''

cloud_setup () {
    openhimCoreMediatorSSLPort=$(kubectl get service openhim-core-service -o=jsonpath={.spec.ports[0].port})
    openhimCoreTransactionPort=$(kubectl get service openhim-core-service -o=jsonpath={.spec.ports[2].port})
    openhimCoreTransactionSSLPort=$(kubectl get service openhim-core-service -o=jsonpath={.spec.ports[1].port})

    while
        openhimCoreHostname=$(kubectl get service openhim-core-service -o=jsonpath="{.status.loadBalancer.ingress[*]['hostname', 'ip']}")
        coreUrlLength=$(expr length "$openhimCoreHostname")
        (( coreUrlLength <= 0 ))
    do
        echo "OpenHIM Core not ready. Sleep 5"
        sleep 5
    done

    openhimCoreMediatorApiUrl="https://$openhimCoreHostname:$openhimCoreMediatorSSLPort"
    openhimCoreTransactionApiUrl="http://$openhimCoreHostname:$openhimCoreTransactionPort"
    openhimCoreTransactionSSLApiUrl="https://$openhimCoreHostname:$openhimCoreTransactionSSLPort"

    # Injecting OpenHIM Core Api url into Console config file
    sed -i -E "s/(\"host\": \")\S*(\")/\1${openhimCoreHostname}\2/" $openhimConsoleVolumePath
    # Injecting OpenHIM Core port into Console config file
    sed -i -E "s/(\"port\": )\S*(,)/\1${openhimCoreMediatorSSLPort}\2/" $openhimConsoleVolumePath

    kubectl apply -k $k8sMainRootFilePath/openhim

    hapiFhirPort=$(kubectl get service hapi-fhir-server-service -o=jsonpath={.spec.ports[0].port})

    while
        hapiFhirServerHostname=$(kubectl get service hapi-fhir-server-service -o=jsonpath="{.status.loadBalancer.ingress[0]['hostname', 'ip']}")
        fhirUrlLength=$(expr length "$hapiFhirServerHostname")
        (( fhirUrlLength <= 0 ))
    do
        echo "HAPI-FHIR not ready. Sleep 5"
        sleep 5
    done

    hapiFhirServerUrl="http://$hapiFhirServerHostname:$hapiFhirPort"

    openhimConsolePort=$(kubectl get service openhim-console-service -o=jsonpath={.spec.ports[0].port})

    while
        openhimConsoleHostname=$(kubectl get service openhim-console-service -o=jsonpath="{.status.loadBalancer.ingress[0]['hostname', 'ip']}")
        consoleUrlLength=$(expr length "$openhimConsoleHostname")
        (( consoleUrlLength <= 0 ))
    do
        echo "OpenHIM Console not ready. Sleep 5"
        sleep 5
    done

    openhimConsoleUrl="http://$openhimConsoleHostname:$openhimConsolePort"
}

local_setup () {
    minikubeIP=$(kubectl config view -o=jsonpath='{.clusters[?(@.name=="minikube")].cluster.server}' | awk '{ split($0,A,/:\/*/) ; print A[2] }')
    openhimCoreMediatorSSLPort=$(kubectl get service openhim-core-service -o=jsonpath={.spec.ports[0].nodePort})
    openhimCoreTransactionPort=$(kubectl get service openhim-core-service -o=jsonpath={.spec.ports[2].nodePort})
    openhimCoreTransactionSSLPort=$(kubectl get service openhim-core-service -o=jsonpath={.spec.ports[1].nodePort})
    hapiFhirPort=$(kubectl get service hapi-fhir-server-service -o=jsonpath={.spec.ports[0].nodePort})

    hapiFhirServerUrl="http://$minikubeIP:$hapiFhirPort"
    openhimCoreMediatorApiUrl="https://$minikubeIP:$openhimCoreMediatorSSLPort"
    openhimCoreTransactionApiUrl="http://$minikubeIP:$openhimCoreTransactionPort"
    openhimCoreTransactionSSLApiUrl="https://$minikubeIP:$openhimCoreTransactionSSLPort"

    # Injecting minikube ip as the hostname of the OpenHIM Core into Console config file
    sed -i -E "s/(\"host\": \")\S*(\")/\1${minikubeIP}\2/" $openhimConsoleVolumePath

    # Injecting OpenHIM Core port into Console config file
    sed -i -E "s/(\"port\": )\S*(,)/\1${openhimCoreMediatorSSLPort}\2/" $openhimConsoleVolumePath

    kubectl apply -k $k8sMainRootFilePath/openhim

    openhimConsolePort=$(kubectl get service openhim-console-service -o=jsonpath={.spec.ports[0].nodePort})

    openhimConsoleUrl="http://$minikubeIP:$openhimConsolePort"
}

print_services_url () {
    printf "\n\nHAPI FHIR Server Url\n--------------------\n"$hapiFhirServerUrl"\n\n"
    printf "OpenHIM Mediator API Url\n------------------------\n"$openhimCoreMediatorApiUrl"\n\n"
    printf "OpenHIM Transaction API Url\n---------------------------\n"$openhimCoreTransactionApiUrl"\n\n"
    printf "OpenHIM Transaction SSL API Url\n-------------------------------\n"$openhimCoreTransactionSSLApiUrl"\n\n"
    printf "OpenHIM Console Url\n===================\n"$openhimConsoleUrl"\n\n"
}

if [ "$1" == "init" ]; then
    # Create persistence volume for the mongo replica set members
    kubectl apply -f $k8sMainRootFilePath/mongo/mongo-volume.yaml

    # Create the replica set
    kubectl apply -f $k8sMainRootFilePath/mongo/mongo-service.yaml -f $k8sMainRootFilePath/mongo/mongo-replica.yaml

    # Set up the replica set
    "$k8sMainRootFilePath"/mongo/initiateReplicaSet.sh

    kubectl apply -k $k8sMainRootFilePath

    envContextName=$(kubectl config get-contexts | grep '*' | awk '{print $2}')
    envContextMinikube=$(echo $envContextName | grep 'minikube')

    if [ $(expr length "$envContextMinikube") -le 0 ]; then
        cloud_setup
    else
        local_setup
    fi

    print_services_url
    printf ">>> The OpenHIM Console Url will take a few minutes to become active <<<\n\n"

    bash "$k8sMainRootFilePath"/../importer/k8s.sh up
elif [ "$1" == "up" ]; then
    kubectl apply -k $k8sMainRootFilePath

    envContextName=$(kubectl config get-contexts | grep '*' | awk '{print $2}')
    envContextMinikube=$(echo $envContextName | grep 'minikube')

    if [ $(expr length "$envContextMinikube") -le 0 ]; then
        cloud_setup
    else
        local_setup
    fi

    print_services_url
elif [ "$1" == "down" ]; then
    kubectl delete deployment openhim-console-deployment
    kubectl delete deployment openhim-core-deployment
    kubectl delete deployment hapi-fhir-server-deployment
    kubectl delete deployment hapi-fhir-mysql-deployment
elif [ "$1" == "destroy" ]; then
    kubectl delete -f $k8sMainRootFilePath/mongo/mongo-volume.yaml
    kubectl delete -f $k8sMainRootFilePath/mongo/mongo-service.yaml -f $k8sMainRootFilePath/mongo/mongo-replica.yaml
    kubectl delete -k $k8sMainRootFilePath
    kubectl delete -k $k8sMainRootFilePath/openhim
    bash "$k8sMainRootFilePath"/../importer/k8s.sh clean
    kubectl delete pvc -l package=core
else
    echo "Valid options are: init, up, down, or destroy"
fi
