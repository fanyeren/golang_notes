#!/bin/bash

HOST_IP=$(host `hostname` | awk '{print $NF}')
# HOST_IP_1 = $(host mesos_master_1 | awk '{print $NF}')
# HOST_IP_2 = $(host mesos_master_2 | awk '{print $NF}')

docker run -d --name zookeeper -p 2181:2181 -p 2888:2888 -p 3888:3888 garland/zookeeper:latest
#docker run -d --name zookeeper_1 --net="host" -e SERVER_ID=1 -e ADDITIONAL_ZOOKEEPER_1=server.1=${HOST_IP_1}:2888:3888 -e ADDITIONAL_ZOOKEEPER_2=server.2=${HOST_IP_2}:2888:3888 -p 2181:2181 -p 2888:2888 -p 3888:3888 garland/zookeeper:latest
#docker run -d --name zookeeper_2 --net="host" -e SERVER_ID=2 -e ADDITIONAL_ZOOKEEPER_1=server.1=${HOST_IP_1}:2888:3888 -e ADDITIONAL_ZOOKEEPER_2=server.2=${HOST_IP_2}:2888:3888 -p 2181:2181 -p 2888:2888 -p 3888:3888 garland/zookeeper:latest

# mesos: http://${HOST_IP}:5050
docker run -d --name mesos_master --net="host" -p 5050:5050 -e "MESOS_HOSTNAME=${HOST_IP}" -e "MESOS_IP=${HOST_IP}" -e "MESOS_ZK=zk://${HOST_IP}:2181/mesos" -e "MESOS_PORT=5050" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_QUORUM=1" -e "MESOS_REGISTRY=in_memory" -e "MESOS_WORK_DIR=/var/lib/mesos" garland/mesosphere-docker-mesos-master:latest
#docker run -d --name mesos_master --net="host" -p 5050:5050 -e "MESOS_HOSTNAME=${HOST_IP_1}" -e "MESOS_IP=${HOST_IP_1}" -e "MESOS_ZK=zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/mesos" -e "MESOS_PORT=5050" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_QUORUM=1" -e "MESOS_REGISTRY=in_memory" -e "MESOS_WORK_DIR=/var/lib/mesos" garland/mesosphere-docker-mesos-master:latest
#docker run -d --name mesos_master --net="host" -p 5050:5050 -e "MESOS_HOSTNAME=${HOST_IP_2}" -e "MESOS_IP=${HOST_IP_2}" -e "MESOS_ZK=zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/mesos" -e "MESOS_PORT=5050" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_QUORUM=1" -e "MESOS_REGISTRY=in_memory" -e "MESOS_WORK_DIR=/var/lib/mesos" garland/mesosphere-docker-mesos-master:latest

# marathon: http://${HOST_IP}:8080
docker run -d --name marathon -p 8080:8080 garland/mesosphere-docker-marathon:latest --master zk://${HOST_IP}:2181/mesos --zk zk://${HOST_IP}:2181/marathon
#docker run -d --name marathon -p 8080:8080 garland/mesosphere-docker-marathon:latest --master zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/mesos --zk zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/marathon

docker run -d --name mesos_slave_1 --entrypoint="mesos-slave" -e "MESOS_MASTER=zk://${HOST_IP}:2181/mesos" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_LOGGING_LEVEL=INFO" garland/mesosphere-docker-mesos-master:latest
#docker run -d --name mesos_slave_1 --entrypoint="mesos-slave" -e "MESOS_MASTER=zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/mesos" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_LOGGING_LEVEL=INFO" garland/mesosphere-docker-mesos-master:latest
#docker run -d --name mesos_slave_2 --entrypoint="mesos-slave" -e "MESOS_MASTER=zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/mesos" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_LOGGING_LEVEL=INFO" garland/mesosphere-docker-mesos-master:latest
# HOST_IP_2 = 

docker run -d --name zookeeper -p 2181:2181 -p 2888:2888 -p 3888:3888 garland/zookeeper:latest
#docker run -d --name zookeeper_1 --net="host" -e SERVER_ID=1 -e ADDITIONAL_ZOOKEEPER_1=server.1=${HOST_IP_1}:2888:3888 -e ADDITIONAL_ZOOKEEPER_2=server.2=${HOST_IP_2}:2888:3888 -p 2181:2181 -p 2888:2888 -p 3888:3888 garland/zookeeper:latest
#docker run -d --name zookeeper_2 --net="host" -e SERVER_ID=2 -e ADDITIONAL_ZOOKEEPER_1=server.1=${HOST_IP_1}:2888:3888 -e ADDITIONAL_ZOOKEEPER_2=server.2=${HOST_IP_2}:2888:3888 -p 2181:2181 -p 2888:2888 -p 3888:3888 garland/zookeeper:latest

# mesos: http://${HOST_IP}:5050
docker run -d --name mesos_master --net="host" -p 5050:5050 -e "MESOS_HOSTNAME=${HOST_IP}" -e "MESOS_IP=${HOST_IP}" -e "MESOS_ZK=zk://${HOST_IP}:2181/mesos" -e "MESOS_PORT=5050" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_QUORUM=1" -e "MESOS_REGISTRY=in_memory" -e "MESOS_WORK_DIR=/var/lib/mesos" garland/mesosphere-docker-mesos-master:latest
#docker run -d --name mesos_master --net="host" -p 5050:5050 -e "MESOS_HOSTNAME=${HOST_IP_1}" -e "MESOS_IP=${HOST_IP_1}" -e "MESOS_ZK=zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/mesos" -e "MESOS_PORT=5050" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_QUORUM=1" -e "MESOS_REGISTRY=in_memory" -e "MESOS_WORK_DIR=/var/lib/mesos" garland/mesosphere-docker-mesos-master:latest
#docker run -d --name mesos_master --net="host" -p 5050:5050 -e "MESOS_HOSTNAME=${HOST_IP_2}" -e "MESOS_IP=${HOST_IP_2}" -e "MESOS_ZK=zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/mesos" -e "MESOS_PORT=5050" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_QUORUM=1" -e "MESOS_REGISTRY=in_memory" -e "MESOS_WORK_DIR=/var/lib/mesos" garland/mesosphere-docker-mesos-master:latest

# marathon: http://${HOST_IP}:8080
docker run -d --name marathon -p 8080:8080 garland/mesosphere-docker-marathon:latest --master zk://${HOST_IP}:2181/mesos --zk zk://${HOST_IP}:2181/marathon
#docker run -d --name marathon -p 8080:8080 garland/mesosphere-docker-marathon:latest --master zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/mesos --zk zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/marathon

docker run -d --name mesos_slave_1 --entrypoint="mesos-slave" -e "MESOS_MASTER=zk://${HOST_IP}:2181/mesos" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_LOGGING_LEVEL=INFO" garland/mesosphere-docker-mesos-master:latest
#docker run -d --name mesos_slave_1 --entrypoint="mesos-slave" -e "MESOS_MASTER=zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/mesos" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_LOGGING_LEVEL=INFO" garland/mesosphere-docker-mesos-master:latest
#docker run -d --name mesos_slave_2 --entrypoint="mesos-slave" -e "MESOS_MASTER=zk://${HOST_IP_1}:2181,${HOST_IP_2}:2181/mesos" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_LOGGING_LEVEL=INFO" garland/mesosphere-docker-mesos-master:latest
