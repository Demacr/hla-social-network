#!/bin/bash

docker-compose down
docker volume rm deploy_social-network-db
docker volume rm deploy_social-network-db-slave
