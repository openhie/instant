#!/bin/bash

k8sMainRootFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

openCrApiUrl=''

cloud_setup () {
    openCrPort=$(kubectl get service opencr -o=jsonpath={.spec.ports[0].port})

    while
        openCrHostname=$(kubectl get service opencr -o=jsonpath="{.status.loadBalancer.ingress[*]['hostname', 'ip']}")
        openCrUrlLength=$(expr length "$openCrHostname")
        (( openCrUrlLength <= 0 ))
    do
        echo "OpenCR not ready. Sleep 5"
        sleep 5
    done

    openCrApiUrl="http://$openCrHostname:$openCrPort"
}

local_setup () {
    minikubeIP=$(kubectl config view -o=jsonpath='{.clusters[?(@.name=="minikube")].cluster.server}' | awk '{ split($0,A,/:\/*/) ; print A[2] }')
    openCrPort=$(kubectl get service opencr -o=jsonpath={.spec.ports[0].nodePort})

    openCrApiUrl="http://$minikubeIP:$openCrPort"
}

print_services_url () {
    printf "\nOpenCr Server Url\n--------------------\n"$openCrApiUrl"\n\n"
    printf "\nOpenCr UI\n--------------------\n"$openCrApiUrl\/crux"\n\n"
}

if [ "$1" == "init" ]; then

    kubectl apply -k $k8sMainRootFilePath
    bash "$k8sMainRootFilePath"/../importer/k8s.sh up

    envContextName=$(kubectl config get-contexts | grep '*' | awk '{print $2}')
    envContextMinikube=$(echo $envContextName | grep 'minikube')

    if [ $(expr length "$envContextMinikube") -le 0 ]; then
        cloud_setup
    else
        local_setup
    fi

    print_services_url
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

    kubectl delete -k $k8sMainRootFilePath

elif [ "$1" == "destroy" ]; then

    kubectl delete -k $k8sMainRootFilePath
    bash "$k8sMainRootFilePath"/../importer/k8s.sh clean

else
    echo "Valid options are: init, up, down, or destroy"
fi
