#!/bin/bash

getExposedServices () {
	mongoUrl="$(minikube service openhim-mongo-dev-service --url)"
	trimmedMongoUrl=$(sed 's/http:\/\///g' <<< "$mongoUrl")
	echo -e "The OpenHIM MongoDB url:\n $trimmedMongoUrl\n"

	mysqlUrl="$(minikube service hapi-fhir-mysql-dev-service --url)"
	trimmedMysqlUrl=$(sed 's/http:\/\///g' <<< "$mysqlUrl")
	echo -e "The HAPI FHIR MySQL url:\n $trimmedMysqlUrl\n"

	echo -e "HAPI FHIR server is accessible at:\n $(minikube service hapi-fhir-server-service --url)/fhir/"
}

applyDevScripts () {
	k8sDevRootFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

	if [ "$1" == "up" ]; then
		kubectl apply -k $k8sDevRootFilePath
		echo -e "\nCurrently in development mode\n"

		getExposedServices
	elif [ "$1" == "destroy" ]; then
		# delete host entry on destroy
		sudo sed -i "/HOST alias for kubernetes Minikube development/d" /etc/hosts

		kubectl delete -k $k8sDevRootFilePath
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