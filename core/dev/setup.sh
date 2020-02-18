#!/bin/bash
set -eu

echo "Setup Kubernetes..."

# Create config map for the console
kubectl create configmap console-config --from-file=./config/console-config.json

# Run the deployments and services
echo 'Deploying mongo'
kubectl apply -f ./mongo

if [ $? -eq 0 ]; then
  echo "Mongo deployment successful"
  echo "Deploying the openhim"
  kubectl apply -f ./openhim

  success=$?

  if [ $success -eq 0 ]; then
    echo 'Openhim deployment successful'
    consoleUrl=$(minikube service console --url)
    coreUrls=$(minikube service core --url)

    echo -e "The urls of the core are:\n$coreUrls"
    echo "Please store the host and port of the url at the top somewhere.\nThese are needed for the configuration of the console"

    read 'Enter yes/no to configure: ' proceed

    if [ "$proceed" == "yes"]; then
      chmod u+x ./configure-console.sh
      ./configure-console.sh
    fi

    echo -e "The console is accessible at url: $consoleUrl\n
      First time login credentials:\n\
      Username: root@openhim.org\n\
      Password: openhim-password"
  fi
fi
