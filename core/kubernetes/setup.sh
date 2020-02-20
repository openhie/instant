#!/bin/bash

echo "Setup Kubernetes..."
kubectl create configmap console-config --from-file=./config/openhim-console/default.json
kubectl create configmap core-config --from-file=./config/openhim-core/development.json
./k8s.sh up
echo "ELBs might take longer to provision (run kubectl get services)..."
