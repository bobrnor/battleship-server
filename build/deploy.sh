#!/usr/bin/env bash

set -e

./build-app.sh
./build-image.sh
./run-mysql.sh
./run-container.sh