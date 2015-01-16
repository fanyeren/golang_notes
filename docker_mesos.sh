#!/bin/bash

HOST_IP=$(host `hostname` | awk '{print $NF}')

docker run -d --name zookeeper -p 2181:2181 -p 2888:2888 -p 3888:3888 garland/zookeeper:latest

docker run -d --name mesos_master --net="host" -p 5050:5050 -e "MESOS_HOSTNAME=${HOST_IP}" -e "MESOS_IP=${HOST_IP}" -e "MESOS_ZK=zk://${HOST_IP}:2181/mesos" -e "MESOS_PORT=5050" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_QUORUM=1" -e "MESOS_REGISTRY=in_memory" -e "MESOS_WORK_DIR=/var/lib/mesos" garland/mesosphere-docker-mesos-master:latest

docker run -d --name marathon -p 8080:8080 garland/mesosphere-docker-marathon:latest --master zk://${HOST_IP}:2181/mesos --zk zk://${HOST_IP}:2181/marathon

docker run -d --name mesos_slave_1 --entrypoint="mesos-slave" -e "MESOS_MASTER=zk://${HOST_IP}:2181/mesos" -e "MESOS_LOG_DIR=/var/log/mesos" -e "MESOS_LOGGING_LEVEL=INFO" garland/mesosphere-docker-mesos-master:latest
