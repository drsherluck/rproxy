#!/bin/bash

echo "Creating root certificate"

openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -sha256 -days 1024 -key ca.key -subj "/CN=localhost-ca" -out ca.crt

echo "Creating tls certificates"

# proxy
proxy_dir=k8s/proxy
openssl genrsa -out $proxy_dir/server.key 2048
openssl req -new -key $proxy_dir/server.key -out $proxy_dir/server.csr -config $proxy_dir/csr.conf
openssl x509 -req -in $proxy_dir/server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out $proxy_dir/server.crt -days 365 -extfile $proxy_dir/csr.conf -extensions usr_ext

# ln -sf ../ca.crt $proxy_dir/ca.crt

