package api

import (
	"fmt"
	"net/http"

	"github.com/ankitsalunkhe/url-shortner/shortner"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/labstack/echo/v4"
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

func (a *API) PostURL(e echo.Context) error {
	var request Request
	if err := e.Bind(&request); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Errorf("invalid body: %w", err),
		)
	}
	base62 := shortner.Base62{}
	shortUrl := base62.Generate(100000000000)
	url := Url{
		LongUrl:  request.Url,
		ShortUrl: shortUrl,
	}

	item, err := attributevalue.MarshalMap(url)
	if err != nil {
		e.JSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("failed: %v", err),
		})
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String(a.db.TableName),
		Item:      item,
	}

	output, err := a.db.Client.PutItem(e.Request().Context(), params)
	if err != nil {
		e.JSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("failed: %v", err),
		})
	}

	println(output)

	return e.JSON(200, Success{
		Url: shortUrl,
	})
}
