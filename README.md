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
    AttributeName=LongUrl,AttributeType=S \
  --key-schema AttributeName=ShortUrl,KeyType=HASH \
  --global-secondary-indexes '[
    {
      "IndexName": "LongUrlIndex",
      "KeySchema":[
        {"AttributeName":"LongUrl","KeyType":"HASH"}
      ],
      "Projection":{
        "ProjectionType":"ALL"
      },
      "ProvisionedThroughput":{
        "ReadCapacityUnits": 5,
        "WriteCapacityUnits": 5
      }
    }
  ]' \
  --provisioned-throughput ReadCapacityUnits=10,WriteCapacityUnits=10 \
  --table-class STANDARD \
  --endpoint-url http://localhost:8000 \
  --region eu-west-1
```

Run below command to check if table is created

```shell
aws dynamodb describe-table --table-name Url --endpoint-url=http://localhost:8000 | grep TableStatus
```

To scan all the rows in table

```shell
aws dynamodb scan --endpoint-url=http://localhost:8000 --table-name Url
```
