#!/bin/bash
set -e

echo "Waiting for dynamodb local to start..."
sleep 5

echo "Setting up environments"
export AWS_ACCESS_KEY_ID=dummy
export AWS_SECRET_ACCESS_KEY=dummy
export AWS_DEFAULT_REGION=eu-west-1

echo "Creating DynamoDB table..."
aws dynamodb create-table \
    --table-name Url \
    --attribute-definitions \
        AttributeName=ShortUrl,AttributeType=S \
        AttributeName=LongUrl,AttributeType=S \
    --key-schema AttributeName=ShortUrl,KeyType=HASH \
    --global-secondary-indexes '[
        {
            "IndexName": "LongUrlIndex",
            "KeySchema": [
                {
                    "AttributeName": "LongUrl",
                    "KeyType": "HASH"
                }
            ],
            "Projection": {
                "ProjectionType": "ALL"
            },
            "ProvisionedThroughput": {
                "ReadCapacityUnits": 5,
                "WriteCapacityUnits": 5
            }
        }
    ]' \
    --provisioned-throughput ReadCapacityUnits=10,WriteCapacityUnits=10 \
    --endpoint-url http://dynamodb-local:8000 \
    --region eu-west-1

echo "Done creating table"
