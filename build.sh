#!/usr/bin/env bash -ex

GOOS=linux go build -o ./bin/linux_service_now_proxy *.go

docker build -t service_now_api .

 docker run --rm -it -p 8080:8080 --env-file .env service_now_api
