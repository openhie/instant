#!/usr/bin/env bash

currentDir=$(dirname $(readlink -f $0))

if [ "$1" == "init" ]; then

    kubectl apply -f $currentDir/manifest.yaml --wait

    # Copy the example project.yaml to the created PVC
    IMAGE=openfn/microservice:v0.3.3
    VOL_MOUNTS='[{"mountPath": "/opt/app/project", "name": "openfn-microservice-data"}]'
    VOLS='[{"name": "openfn-microservice-data", "persistentVolumeClaim": {"claimName": "openfn-microservice-data"}}]'
    SUFFIX=$(date +%s | shasum | base64 | fold -w 10 | head -1 | tr '[:upper:]' '[:lower:]')

    cd $currentDir/../../docker/config

    tar cf - project.yaml | kubectl run -n openfn \
        -i --rm --restart=Never \
        --image=${IMAGE} \
        openfn-microservice-${SUFFIX} \
        --overrides "{
        \"spec\": {
          \"hostNetwork\": true,
          \"containers\":[
            {
              \"args\": [\"sh\", \"-c\", \"tar xvf - -C /opt/app/project\"],
              \"stdin\": true,
              \"name\": \"openfn-microservice\",
              \"image\": \"${IMAGE}\",
              \"volumeMounts\": ${VOL_MOUNTS}
            }
          ],
          \"volumes\": ${VOLS}
        } }"



elif [ "$1" == "up" ]; then

    kubectl apply -f $currentDir/manifest.yaml

elif [ "$1" == "down" ]; then

    kubectl delete -f $currentDir/manifest.yaml

elif [ "$1" == "destroy" ]; then

    kubectl delete -f $currentDir/manifest.yaml
    
else
    echo "Valid options are: init, up, down, or destroy"
fi
