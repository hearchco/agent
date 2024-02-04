#!/usr/bin/env bash
go test $(go list ./... | grep /engines/)
