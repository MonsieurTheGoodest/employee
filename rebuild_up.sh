#!/bin/bash
echo docker compose down
echo docker rmi employee:v01
echo docker build -t employee:v01 ./
echo docker compose up

docker compose down
docker rmi employee:v01
docker build -t  employee:v01 ./
docker compose up