#!/bin/sh

set -eu

echo 'Initiating the mongo replica set'

config='{"_id":"mongo-set","members":[{"_id":0,"priority":1,"host":"mongo-0.mongo-service:27017"},
{"_id":1,"priority":0.5,"host":"mongo-1.mongo-service:27017"},{"_id":2,"priority":0.5,"host":"mongo-2.mongo-service:27017"}]}'

# Sleep to ensure all the mongo instances for the replica set are up and running
sleep 90

kubectl exec -i mongo-0 -- mongo --eval "rs.initiate($config)"

## Allow the replica members to connect to each other and be in sync
sleep 60

echo 'Mongo Replica setup finished'
