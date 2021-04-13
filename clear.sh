#!/bin/bash

docker kill $(docker ps -q)
docker rm $(docker ps -aq)
docker volume rm $(docker volume ls -q)
docker image rm cdc-client
docker image rm cdc-service