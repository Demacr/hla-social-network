#!/bin/bash

#
# This script is based on `build.sh` located in
# https://github.com/vbabak/docker-mysql-master-slave repo
#

docker-compose up mariadb mariadb-slave -d

until docker-compose exec mariadb sh -c 'mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e ";"'
do
    echo "Waiting for mariadb database connection..."
    sleep 4
done

echo "--- Master is up"

MS_STATUS=`docker exec mariadb sh -c 'mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e "SHOW MASTER STATUS"'`
CURRENT_LOG=`echo $MS_STATUS | awk '{print $5}'`
CURRENT_POS=`echo $MS_STATUS | awk '{print $6}'`

echo "--- Get log file and position: $CURRENT_LOG $CURRENT_POS"

until docker-compose exec mariadb-slave sh -c 'mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e ";"'
do
    echo "Waiting for mariadb-slave database connection..."
    sleep 4
done

echo "--- Slave is up"

start_slave_stmt="CHANGE MASTER TO MASTER_HOST='mariadb',MASTER_USER='replication',MASTER_PASSWORD='replication-password',MASTER_LOG_FILE='$CURRENT_LOG',MASTER_LOG_POS=$CURRENT_POS; START SLAVE;"
start_slave_cmd='mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e "'
start_slave_cmd+="$start_slave_stmt"
start_slave_cmd+='"'
docker exec mariadb-slave sh -c "$start_slave_cmd"

echo "--- slave status:"
docker exec mariadb-slave sh -c 'mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e "SHOW SLAVE STATUS \G"'

echo "--- starting social-network service:"
docker-compose up -d
