#!/bin/bash
set -u

echo "Setup Kubernetes..."

# Check if console configMap exists
kubectl get configmap console-config

if [ $? -eq 0 ]; then 
  # Delete configMap for the console if it exists
  kubectl delete configmap console-config
fi

kubectl create configmap console-config --from-file=./config/console-config.json

# Check if hapi-fhir configMap exists
kubectl get configmap console-config

if [ $? -eq 0 ]; then
  # Delete configMap for the hapi fhir if it exists
  kubectl delete configmap hapi-fhir-config
fi

kubectl create configmap hapi-fhir-config --from-file=./config/hapi.properties

error='none'

# Run the deployments and services
echo 'Deploying mongo'
kubectl apply -f ./mongo

if [ $? -eq 0 ]; then
  echo "Mongo deployment successful"
  echo "Deploying the OpenHIM"

  kubectl apply -f ./openhim
  if [ $? -eq 0 ]; then
    echo 'Openhim deployment successful'
  else
    error='Openhim deployment failure'
  fi
else
  error='Mongo deployment failure'
fi

echo 'Deploying MySQL database'

kubectl apply -f ./mysql
if [ $? -eq 0 ]; then
  echo 'MySQL database deployed'
else
  error='MySQL deployment failure'
fi

echo 'Deploying HAPI-fhir'

kubectl apply -f ./fhir
if [ $? -eq 0 ]; then
  echo 'HAPI-fhir deployment successful'
else
  error='HAPI-fhir deployment failure'
fi

if [ $error == 'none' ]; then
  echo 'Setup completed successfully'
else
  echo "Setup failure: $error"
fi
