# URL Shortner


## Getting Started

To set up docker containers running and ready for you to run, run below command

```shell
make infra
```

To create new table
```shell
   aws dynamodb create-table \
    --table-name Url \
    --attribute-definitions \
        AttributeName=ShortUrl,AttributeType=S \
    --global-secondary-indexes
    --key-schema AttributeName=ShortUrl,KeyType=HASH \
    --table-class STANDARD \
    --provisioned-throughput \
            ReadCapacityUnits=10,WriteCapacityUnits=10 \
    --endpoint-url=http://localhost:8000
```

Run below command to check if table is created
```shell
aws dynamodb describe-table --table-name Url --endpoint-url=http://localhost:8000 | grep TableStatus
```

To scan all the rows in table
```shell
aws dynamodb scan --endpoint-url=http://localhost:8000 --table-name Url
```