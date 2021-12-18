#!/bin/bash
export MINIKUBE_IN_STYLE=false

# create certificates
./certificates.sh

# cp ca minikube 
mkdir -p $HOME/.minikube/certs
cp ca.crt $HOME/.minikube/certs/testca.crt

# start and build
minikube start --embed-certs
eval $(minikube docker-env)
docker build -t server server/.
docker build -t proxy proxy/.

# add tls secrets
kubectl delete secret tls-secret-server tls-secret-proxy
kubectl create secret tls tls-secret-server --cert=server/server.crt --key=server/server.key
kubectl create secret tls tls-secret-proxy  --cert=proxy/server.crt --key=proxy/server.key

# create pods and serivces
kubectl apply -k server
kubectl apply -k proxy

# start ingress
kubectl apply -f ingress.yaml

# start service
#minikube service proxy-service 
