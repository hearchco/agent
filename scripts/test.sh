#!/usr/bin/env bash
go test $(go list ./... | grep -v /engines/ | grep -v /redis | grep -v /dynamodb)
