#!/bin/bash

#
# This script is based on `build.sh` located in
# https://github.com/vbabak/docker-mysql-master-slave repo
#

docker-compose -f docker-compose-async.yaml up mariadb mariadb-slave-1 mariadb-slave-2 -d

until docker-compose -f docker-compose-async.yaml exec mariadb sh -c 'mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e ";"'
do
    echo "Waiting for mariadb database connection..."
    sleep 4
done

echo "--- Master is up"

MS_STATUS=`docker exec mariadb sh -c 'mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e "SHOW MASTER STATUS"'`
CURRENT_LOG=`echo $MS_STATUS | awk '{print $5}'`
CURRENT_POS=`echo $MS_STATUS | awk '{print $6}'`

echo "--- Get log file and position: $CURRENT_LOG $CURRENT_POS"

until docker-compose -f docker-compose-async.yaml exec mariadb-slave-1 sh -c 'mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e ";"'
do
    echo "Waiting for mariadb-slave-1 database connection..."
    sleep 4
done

until docker-compose -f docker-compose-async.yaml exec mariadb-slave-2 sh -c 'mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e ";"'
do
    echo "Waiting for mariadb-slave-2 database connection..."
    sleep 4
done

echo "--- Slave is up"

start_slave_stmt="CHANGE MASTER TO MASTER_HOST='mariadb',MASTER_USER='replication',MASTER_PASSWORD='replication-password',master_use_gtid=current_pos; START SLAVE;"
start_slave_cmd='mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e "'
start_slave_cmd+="$start_slave_stmt"
start_slave_cmd+='"'
docker exec mariadb-slave-1 sh -c "$start_slave_cmd"
docker exec mariadb-slave-2 sh -c "$start_slave_cmd"

echo "--- slave status:"
docker exec mariadb-slave-1 sh -c 'mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e "SHOW SLAVE STATUS \G"'
docker exec mariadb-slave-2 sh -c 'mysql -h localhost -u root -p$MYSQL_ROOT_PASSWORD -e "SHOW SLAVE STATUS \G"'

echo "--- starting social-network service:"
docker-compose -f docker-compose-async.yaml up -d
