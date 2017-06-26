#!/usr/bin/env bash

GOPATH="$(pwd)/../../../../../"
GOOS=linux GOARCH=amd64 go build ../battleship.go
