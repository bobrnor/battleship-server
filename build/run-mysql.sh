#!/usr/bin/env bash

docker stop battleship-db | true
docker rm battleship-db | true
docker run -d --name battleship-db \
    -v /private/var/battleship/mysql:/var/lib/mysql \
    -e MYSQL_DATABASE=battleship \
    -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
    -e MYSQL_ROOT_HOST=172.25.0.3 \
    --network=battleship-network \
    --ip=172.25.0.2 \
    mysql/mysql-server
docker logs -f battleship-db
