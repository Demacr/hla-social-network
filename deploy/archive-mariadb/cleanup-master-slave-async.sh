#!/bin/bash

docker-compose -f docker-compose-async.yaml down
docker volume rm deploy_social-network-db-async-master
docker volume rm deploy_social-network-db-slave-1
docker volume rm deploy_social-network-db-slave-2
