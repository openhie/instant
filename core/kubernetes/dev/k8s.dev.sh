#!/bin/bash

getExposedServices () {
	mongoUrl="$(minikube service openhim-mongo-dev-service --url)"
	trimmedMongoUrl=$(sed 's/http:\/\///g' <<< "$mongoUrl")
	echo -e "The OpenHIM MongoDB url:\n $trimmedMongoUrl\n"

	mysqlUrl="$(minikube service hapi-fhir-mysql-dev-service --url)"
	trimmedMysqlUrl=$(sed 's/http:\/\///g' <<< "$mysqlUrl")
	echo -e "The HAPI FHIR MySQL url:\n $trimmedMysqlUrl\n"

	echo -e "HAPI FHIR server is accessible at:\n https://$(kubectl get ingress hapi-fhir-server-ingress -o jsonpath="{..host}")/hapi-fhir-jpaserver/fhir/"
}

applyDevScripts () {
	if [ "$1" == "up" ]; then
		kubectl apply -k .
		echo -e "\nCurrently in development mode\n"

		# create HOST entry for ingress
    sudo sed -i "/HOST alias for kubernetes Minikube development/d" /etc/hosts
    echo -e "\n$(minikube ip) $(kubectl get ingress hapi-fhir-server-ingress -o jsonpath="{..host}") # HOST alias for kubernetes Minikube development" | sudo tee -a /etc/hosts

		getExposedServices
	elif [ "$1" == "destroy" ]; then
		# delete host entry on destroy
		sudo sed -i "/HOST alias for kubernetes Minikube development/d" /etc/hosts

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