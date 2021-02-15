#!/bin/bash

while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -c=*|--custom-package=*)
    CUSTOM_PACKAGES+=("${1#*=}")
    shift # past argument=value
    ;;
    init|destroy|down|up|test)
    ACTION=$1
    shift # past value
    ;;
    -*)
    OTHER_FLAGS+=("${key} ${2}")
    shift # past argument value
    shift # past argument value
    ;;
    *)
    PACKAGES+=($1)
    shift # past value
    ;;
esac
done

printf "\nCommand Summary\n---------------\n"
echo "> Custom Package Locations: ${CUSTOM_PACKAGES[@]}"
echo "> Action: ${ACTION}"
echo "> Packages: ${PACKAGES[@]}"
echo "> Other Flags: ${OTHER_FLAGS[@]}"
echo

if [[ $ACTION = "init" ]]
then 
  echo "Delete a pre-existing instant volume..."

  docker volume rm instant
fi

echo "Creating fresh instant container with volumes..."

docker create -it --rm\
  --mount='type=volume,src=instant,dst=/instant' \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ~/.kube/config:/root/.kube/config:ro \
  -v ~/.minikube:/home/$USER/.minikube:ro \
  --network host \
  --name instant-openhie \
  openhie/instant:latest \
  $ACTION \
  ${OTHER_FLAGS[@]} \
  ${PACKAGES[@]}

echo "Adding 3rd party packages to instant volume:"

for customPackage in "${CUSTOM_PACKAGES[@]}"
do
  echo "- ${customPackage}"
  docker cp $customPackage instant-openhie:instant/
done

echo "Run Instant OpenHIE Installer Container"

docker start -a instant-openhie

if [[ $ACTION = "destroy" ]]
then 
  echo "Delete instant volume..."

  docker volume rm instant
fi
