#!/usr/bin/env bash

export AWS_REGION=hearchco-test-1
export AWS_ACCESS_KEY_ID=hearchco
export AWS_SECRET_ACCESS_KEY=hearchco
export DYNAMODB_TABLE=hearchco_test
export DYNAMODB_ENDPOINT=http://localhost:8000

docker run --rm --name hearchco-dynamodb -d -p 8000:8000 docker.io/amazon/dynamodb-local && \
sleep 5 && \
aws dynamodb create-table \
    --table-name $DYNAMODB_TABLE \
    --attribute-definitions AttributeName=Key,AttributeType=S \
    --key-schema AttributeName=Key,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url $DYNAMODB_ENDPOINT && \
aws dynamodb update-time-to-live \
    --table-name $DYNAMODB_TABLE \
    --time-to-live-specification "Enabled=true, AttributeName=TTL" \
    --endpoint-url $DYNAMODB_ENDPOINT && \
go test $(go list ./... | grep /dynamodb) -count=1

docker stop hearchco-dynamodb