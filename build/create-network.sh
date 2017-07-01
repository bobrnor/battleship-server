#!/usr/bin/env bash

#docker network rm battleship-network

docker network create -d bridge --subnet 172.25.0.0/16 battleship-network