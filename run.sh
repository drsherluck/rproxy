#!/bin/bash

# start and build
minikube start
eval $(minikube discord-env)
docker build -t server server/.
docker build -t proxy proxy/.

# create pods and serivces
kubectl apply -k server
kubectl apply -k proxy

# start service
minikube service proxy-service
