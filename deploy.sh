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

getRepoName () {
  readarray -d . -t urlArray <<< "$1"
  len=${#urlArray[*]}
  readarray -d / -t urlPathArray <<< "${urlArray[len-2]}"
  len1=${#urlPathArray[*]}
  echo "${urlPathArray[len1-1]}"
}

isUrl () {
  if [[ $1 =~ "git" || $1 =~ "http" ]]; then
    echo "true"
  else
    echo ""
  fi
}

downloadPackage () {
  if [[ $1 =~ ".git" ]]; then
    git clone $1

    # Get the folder name
    repoName=$(getRepoName $1)

    docker cp "./$repoName" instant-openhie:instant/

    # Remove the downloaded folder
    rm -rf "./$repoName"
  elif [[ $1 =~ ".zip" ]]; then
    curl -L "$1" --output "./temp.zip"
    unzip "temp.zip" -d "./temp"
    docker cp "./temp/." instant-openhie:instant/

    # Remove the downloaded folders
    rm -rf "./temp.zip" "./temp"
  elif [[ $1 =~ ".tar" ]]; then
    curl -L "$1" --output "./temp.tar.gz"
    mkdir "temp"
    tar -xf "temp.tar.gz" -C "./temp"
    docker cp "./temp/." instant-openhie:instant/

    # Remove the downloaded folders
    rm -rf "./temp.tar.gz" "./temp"
  else
    echo "Download url not supported. Only github repos and zip or tar files are supported"
  fi
}

for customPackage in "${CUSTOM_PACKAGES[@]}"
do
  echo "${customPackage}"
  if [ $(isUrl $customPackage) ]; then
    downloadPackage $customPackage
  else
    docker cp $customPackage instant-openhie:instant/
  fi
done

echo "Run Instant OpenHIE Installer Container"

docker start -a instant-openhie

if [[ $ACTION = "destroy" ]]
then 
  echo "Delete instant volume..."

  docker volume rm instant
fi
