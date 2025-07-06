package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (db *DB) UpsertUrl(ctx context.Context, url Url) error {
	item, err := attributevalue.MarshalMap(url)
	if err != nil {
		return fmt.Errorf("marshal url into attibute value: %w", err)
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String(db.tableName),
		Item:      item,
	}

	_, err = db.client.PutItem(ctx, params)
	if err != nil {
		return fmt.Errorf("upserting item into db: %w", err)
	}
	return nil
}

func (db *DB) GetLongUrl(ctx context.Context, url Url) (string, error) {
	response, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: url.GetKey(), TableName: aws.String(db.tableName),
	})
	if err != nil {
		return "", fmt.Errorf("getting item from db: %w", err)
	}

	var longUrl string

	err = attributevalue.Unmarshal(response.Item["LongUrl"], &longUrl)
	if err != nil {
		return "", fmt.Errorf("unmarshalling response from db: %w", err)
	}

	return longUrl, nil
}

func (db *DB) DeletUrl(ctx context.Context, url Url) error {
	_, err := db.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: url.GetKey(), TableName: aws.String(db.tableName),
	})
	if err != nil {
		return fmt.Errorf("deleting item from db: %w", err)
	}

	return nil
}

func (db *DB) GetShortUrl(ctx context.Context, longUrl string) (string, error) {
	output, err := db.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(db.tableName),
		IndexName:              aws.String("LongUrlIndex"),
		KeyConditionExpression: aws.String("LongUrl = :longUrl"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":longUrl": &types.AttributeValueMemberS{Value: longUrl},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return "", fmt.Errorf("query DynamoDB: %w", err)
	}

	if output.Count < 1 {
		return "", nil
	}

	if output.Count > 1 {
		return "", fmt.Errorf("more than 1 url found in db for this longUrl")
	}

	var shortUrl string
	err = attributevalue.Unmarshal(output.Items[0]["ShortUrl"], &shortUrl)
	if err != nil {
		return "", fmt.Errorf("unmarshalling response from db: %w", err)
	}

	return shortUrl, nil
}
