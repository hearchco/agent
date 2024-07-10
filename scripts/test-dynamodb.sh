#!/usr/bin/env bash
go test $(go list ./... | grep /dynamodb) -count=1
