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
docker build -t server src/server/.
docker build -t proxy src/proxy/.

# add tls secrets
sdir=k8s/server
pdir=k8s/proxy
kubectl delete secret tls-secret-server tls-secret-proxy
kubectl create secret tls tls-secret-server --cert=$sdir/server.crt --key=$sdir/server.key
kubectl create secret tls tls-secret-proxy  --cert=$pdir/server.crt --key=$pdir/server.key

# create pods and serivces
kubectl apply -k k8s/server
kubectl apply -k k8s/proxy

# start ingress
kubectl apply -f k8s/ingress.yaml

# start service
#minikube service proxy-service 
