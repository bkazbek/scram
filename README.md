# SCRAM TCP Server / Client

Written on Golang with help of https://github.com/xdg-go/scram

### How to run

### Using Makefile

``make server``

``make client``

### Using Docker

#### Build docker image

``docker build -t scram_tcp .``

#### Before running server and client, we should configure network

``docker network create scram_tcp_network``

#### Run server

``docker run -dit -p 9001:9001 --network=scram_tcp_network --name scram_server scram_tcp serve``

#### Run client

``docker run -dit --network=scram_tcp_network --name scram_client scram_tcp client --address scram_server:9001 --password tcp_password --username tcp``

## Notes

Currently, server credentials stored in ./internal/server.go as default creds. To customize use:

``docker run -dit -p 9001:9001 scram_tcp serve --username tcp_custom --password tcp_password_custom --port 9092``

## Clean up

``docker stop scram_server``

``docker rm scram_client``

``docker rm scram_server``

``docker rmi $(docker images | grep 'scram_tcp' -a -q)``

``docker network rm scram_tcp_network``

