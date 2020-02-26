#!/bin/bash

getExposedServices () {
	mongoUrl="$(minikube service openhim-mongo-dev-service --url)"
	trimmedMongoUrl=$(sed 's/http:\/\///g' <<< "$mongoUrl")
	echo -e "The OpenHIM MongoDB url:\n $trimmedMongoUrl\n"

	mysqlUrl="$(minikube service hapi-fhir-mysql-dev-service --url)"
	trimmedMysqlUrl=$(sed 's/http:\/\///g' <<< "$mysqlUrl")
	echo -e "The HAPI FHIR MySQL url:\n $trimmedMysqlUrl\n"

	echo -e "HAPI FHIR server is accessible at:\n https://hapi-fhir-server-dev.instant/hapi-fhir-jpaserver/fhir/"
}

applyDevScripts () {
	if [ "$1" == "up" ]; then
		kubectl apply -k .
		echo -e "\nCurrently in development mode\n"

		getExposedServices
	elif [ "$1" == "destroy" ]; then
		kubectl delete -k .
	else
		echo "Valid options are: up or destroy"
	fi
}

while true; do
	read -p "Are the Core Package services running? Dev mode depends on these services (y/n)" yn
	case $yn in
		[Yy]* ) applyDevScripts $1; break;;
		[Nn]* ) exit;;
		* ) echo "Please answer yes or no.";;
	esac
done