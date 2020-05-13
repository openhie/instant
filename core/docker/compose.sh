composeFilePath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

if [ "$1" == "up" ]; then
    docker-compose -f "$composeFilePath"/docker-compose-mongo.yml up -d

    # Set up the replica set
    "$composeFilePath"/initiateReplicaSet.sh

    docker-compose -f "$composeFilePath"/docker-compose.yml -f "$composeFilePath"/docker-compose.dev.yml -f "$composeFilePath"/importer/docker-compose.config.yml up -d
elif [ "$1" == "down" ]; then
    docker-compose -f "$composeFilePath"/docker-compose.yml -f "$composeFilePath"/docker-compose.dev.yml -f "$composeFilePath"/importer/docker-compose.config.yml stop
elif [ "$1" == "destroy" ]; then
    docker-compose -f "$composeFilePath"/docker-compose.yml -f "$composeFilePath"/docker-compose.dev.yml -f "$composeFilePath"/importer/docker-compose.config.yml down -v
else
    echo "Valid options are: up, down, or destroy"
fi
