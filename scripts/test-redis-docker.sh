#!/usr/bin/env bash
docker run --rm --name hearchco-redis -d -p 6379:6379 docker.io/library/redis && \
go test $(go list ./... | grep /redis) -count=1
docker stop hearchco-redis