#!/bin/bash

res="$(kubectl get service openhim-mongo-dev-service)"

if [ "$?" == "0" ]; then
  url="$(minikube service openhim-mongo-dev-service --url)"

  echo -e "The mongo url:\n $url\n"
else
  echo "Mongo dev not started (in dev mode)"
fi

res1="$(kubectl get service hapi-fhir-mysql-dev-service)"

if [ "$?" == "0" ]; then
  url="$(minikube service hapi-fhir-mysql-dev-service --url)"

  echo -e "The mysql url:\n $url\n"
else
  echo "Mysql fhir server not started (in dev mode)"
fi

echo -e "Hapi fhir server is accessible at:\n https://hapi-fhir-server-dev.instant/hapi-fhir-jpaserver/fhir/"
