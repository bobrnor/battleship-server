#!/usr/bin/env bash

docker stop battleship-server | true
docker rm battleship-server | true
docker run -d --name=battleship-server \
		--volume /private/var/battleship/server:/external \
		--network=bridge \
		--network=battleship-network \
        --ip=172.25.0.3 \
		battleship-server
docker logs -f battleship-server