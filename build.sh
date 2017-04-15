#!/usr/bin/env bash

GOOS=linux go build -o ./bin/linux_service_now_connector *.go

docker build -t service_now_api .

docker run --rm -it -p 8080:8080 --env-file .env service_now_api
