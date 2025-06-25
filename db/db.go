package db

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

const (
	DynamoRegion = "eu-west-1"
)

type DB struct {
	Client    *dynamodb.Client
	TableName string
}

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
		Client:    svc,
		TableName: "Url",
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
