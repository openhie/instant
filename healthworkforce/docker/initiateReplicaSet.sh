#!/bin/sh

set -eu

echo 'Initiating the mongo replica set'

configPart1='{"_id":"mapper-mongo-set","members":[{"_id":0,"priority":1,"host":"mapper-mongo-1:27017"},'
configPart2='{"_id":1,"priority":0.5,"host":"mapper-mongo-2:27017"},{"_id":2,"priority":0.5,"host":"mapper-mongo-3:27017"}]}'

config="$configPart1$configPart2"

docker exec -i mapper-mongo-1 mongo --eval "rs.initiate($config)"

echo 'Replica set successfully set up'

docker restart mcsdMediator

echo 'Mediator successfully restarted'
