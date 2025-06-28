package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (db *DB) UpsertUrl(ctx context.Context, url Url) error {
	item, err := attributevalue.MarshalMap(url)
	if err != nil {
		return fmt.Errorf("marshal url into attibute value: %w", err)
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String(db.TableName),
		Item:      item,
	}

	_, err = db.Client.PutItem(ctx, params)
	if err != nil {
		return fmt.Errorf("upserting item into db: %w", err)
	}
	return nil
}

func (db *DB) GetUrl(ctx context.Context, url Url) (string, error) {
	response, err := db.Client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: url.GetKey(), TableName: aws.String(db.TableName),
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
