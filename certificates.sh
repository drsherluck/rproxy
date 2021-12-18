#!/bin/bash

echo "Creating root certificate"

openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -sha256 -days 1024 -key ca.key -subj "/CN=localhost-ca" -out ca.crt

echo "Creating tls certificates"

# server 
server_dir=server
openssl genrsa -out $server_dir/server.key 2048
openssl req -new -key $server_dir/server.key -out $server_dir/server.csr -config $server_dir/csr.conf
openssl x509 -req -in $server_dir/server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out $server_dir/server.crt -days 365 -extfile $server_dir/csr.conf -extensions usr_ext

# proxy
proxy_dir=proxy
openssl genrsa -out $proxy_dir/server.key 2048
openssl req -new -key $proxy_dir/server.key -out $proxy_dir/server.csr -config $proxy_dir/csr.conf
openssl x509 -req -in $proxy_dir/server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out $proxy_dir/server.crt -days 365 -extfile $proxy_dir/csr.conf -extensions usr_ext

ln -sf ../ca.crt $proxy_dir/ca.crt

