#!/bin/bash

echo "Setup Kubernetes..."
kubectl create configmap console-config --from-file=../config/openhim-console/default.json
kubectl create configmap core-config --from-file=../config/openhim-core/development.json
kubectl create configmap hapi-fhir-config --from-file=../config/hapi-fhir/hapi.properties
./k8s.sh up "$1"
echo "ELBs might take longer to provision (run kubectl get services)..."
