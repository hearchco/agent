#!/usr/bin/env bash
go test $(go list ./... | grep /redis) -count=1
