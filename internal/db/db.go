package db

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

const (
	DynamoRegion = "eu-west-1"
)

type Url struct {
	ShortUrl string `dynamodbav:"ShortUrl"`
	LongUrl  string `dynamodbav:"LongUrl"`
}

func (url Url) GetKey() map[string]types.AttributeValue {
	shortUrl, err := attributevalue.Marshal(url.ShortUrl)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"ShortUrl": shortUrl}
}

type Database interface {
	UpsertUrl(ctx context.Context, url Url) error
	GetLongUrl(ctx context.Context, url Url) (string, error)
	DeletUrl(ctx context.Context, url Url) error
	GetShortUrl(ctx context.Context, longUrl string) (string, error)
}

type DB struct {
	client    *dynamodb.Client
	tableName string
}

var _ Database = (*DB)(nil)

func New() (*DB, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(DynamoRegion))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolverV2 = &LocalDynamoEndpointResolver{}
		o.Credentials = credentials.NewStaticCredentialsProvider(
			"testkey",
			"testsecret",
			"testsession",
		)
	})

	return &DB{
		client:    svc,
		tableName: "Url",
	}, nil
}

// &LocalDynamoEndpointResolver implements a custom endpoint resolver for local DynamoDB.
type LocalDynamoEndpointResolver struct{}

func (*LocalDynamoEndpointResolver) ResolveEndpoint(context.Context, dynamodb.EndpointParameters) (smithyendpoints.Endpoint, error) {
	u, err := url.Parse("http://localhost:8000")
	if err != nil {
		// return endpoints.Endpoint{}, fmt.Errorf("parsing DynamoDB local URL: %w", err)
		return smithyendpoints.Endpoint{}, fmt.Errorf("parsing DynamoDB local URL: %w", err)
	}

	return smithyendpoints.Endpoint{
		URI: *u,
	}, nil
}
