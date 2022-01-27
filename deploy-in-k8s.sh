#!/bin/bash

echo "Building Redis DB and Redis Insight"
docker build -t redis:v1.0 -f ./Dockerfile-Redis ./
docker tag redis:v1.0 10.101.10.244:5000/redis:v1.0
docker push 10.101.10.244:5000/redis:v1.0
echo "Running Redis DB"
kubectl create -f redis-db.yaml
kubectl create -f redis-insight.yaml

echo "Building IP Anonymization Service Backend"
docker build -t ip-anonymization:v1.0 -f ./Dockerfile ./
docker tag ip-anonymization:v1.0 10.101.10.244:5000/ip-anonymization:v1.0
docker push 10.101.10.244:5000/ip-anonymization:v1.0
echo "Running IP Anonymization Service Backend"
kubectl create -f deployment.yaml