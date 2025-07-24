package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"

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

type Config struct {
	Host string `envconfig:"DB_HOST" required:"true"`
	Port int    `envconfig:"DB_PORT" required:"true"`
}

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

func New(dbConfig Config) (*DB, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(DynamoRegion))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolverV2 = &CustomDynamoEndpointResolver{
			Host: dbConfig.Host,
			Port: dbConfig.Port,
		}
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

// &CustomDynamoEndpointResolver implements a custom endpoint resolver for local DynamoDB.
type CustomDynamoEndpointResolver struct {
	Host string
	Port int
}

func (cd *CustomDynamoEndpointResolver) ResolveEndpoint(context.Context, dynamodb.EndpointParameters) (smithyendpoints.Endpoint, error) {
	u, err := url.Parse("http://" + cd.Host + ":" + strconv.Itoa(cd.Port))
	if err != nil {
		return smithyendpoints.Endpoint{}, fmt.Errorf("parsing DynamoDB local URL: %w", err)
	}

	return smithyendpoints.Endpoint{
		URI: *u,
	}, nil
}
