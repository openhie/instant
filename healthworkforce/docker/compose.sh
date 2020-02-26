if [ "$1" == "up" ]; then
    docker-compose up -d
elif [ "$1" == "down" ]; then
    docker-compose stop
elif [ "$1" == "destroy" ]; then
    docker-compose down
else
    echo "Valid options are: up, down, or destroy"
fi